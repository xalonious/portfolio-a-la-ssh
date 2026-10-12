package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/joho/godotenv"
	"github.com/muesli/termenv"

	"github.com/xalonious/portfolio-a-la-ssh/internal/achievements"
	"github.com/xalonious/portfolio-a-la-ssh/internal/content"
	"github.com/xalonious/portfolio-a-la-ssh/internal/projects"
	"github.com/xalonious/portfolio-a-la-ssh/internal/ui"
)

const (
	defaultAddress = ":2323"
	defaultHostKey = "./host_ed25519"
)

func main() {
	log.SetOutput(os.Stdout)
	loadDotEnv()
	lipgloss.SetColorProfile(termenv.TrueColor)

	address := envOrDefault("SSH_PORTFOLIO_ADDR", defaultAddress)
	hostKeyPath := envOrDefault("SSH_PORTFOLIO_HOST_KEY", defaultHostKey)
	databasePath := os.Getenv("PORTFOLIO_DATABASE_PATH")
	achievementTracker := achievements.NewTracker(databasePath, log.Default())
	defer achievementTracker.Close()
	projectRepository := projects.NewSQLiteRepository(
		databasePath,
		"https://"+content.Data.Domain+"/projects",
		log.Default(),
	)

	server, err := wish.NewServer(
		wish.WithAddress(address),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithMiddleware(
			bm.Middleware(func(session ssh.Session) (tea.Model, []tea.ProgramOption) {
				width, height := initialSize(session)

				ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
				defer cancel()
				loadedProjects, loadErr := projectRepository.Load(ctx)
				if loadErr != nil {
					log.Printf("could not load published projects database=%q: %v", databasePath, loadErr)
				}

				portfolio := content.Data
				portfolio.Projects = loadedProjects
				return ui.New(width, height, portfolio, loadErr), []tea.ProgramOption{tea.WithAltScreen()}
			}),
			connectionLogger(achievementTracker),
		),
	)
	if err != nil {
		log.Fatalf("could not create SSH server: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("SSH portfolio listening on %s", address)
		log.Printf("using host key %s", hostKeyPath)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-done
	log.Println("shutting down SSH portfolio...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}

func loadDotEnv() {
	executablePath, err := os.Executable()
	if err != nil {
		log.Printf("could not locate executable for .env loading: %v", err)
		return
	}

	envPath := filepath.Join(filepath.Dir(executablePath), ".env")
	if err := godotenv.Load(envPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Printf("could not load environment file path=%q: %v", envPath, err)
		}
		return
	}
	log.Printf("loaded environment file path=%q", envPath)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

type connectionRecorder interface {
	Record(string, time.Time)
}

type recordedSessionKey struct{ session ssh.Session }

func connectionLogger(recorder connectionRecorder) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(session ssh.Session) {
			term := terminalName(session)
			if session.User() == "root" && term == "" {
				_ = session.Exit(1)
				_ = session.Close()
				return
			}
			startedAt := time.Now()
			username := session.User()
			tracked := achievements.IsVisitorToken(username)
			loggedUsername := username
			if tracked {
				loggedUsername = "[achievement visitor]"
			}

			log.Printf(
				"SSH session connected remote=%s user=%q term=%q",
				session.RemoteAddr(),
				loggedUsername,
				term,
			)

			_, _, interactive := session.Pty()
			ctx := session.Context()
			if recorder != nil && tracked && interactive && len(session.Command()) == 0 && session.Subsystem() == "" && ctx.Err() == nil {
				key := recordedSessionKey{session: session}
				ctx.Lock()
				first := ctx.Value(key) == nil
				if first {
					ctx.SetValue(key, true)
				}
				ctx.Unlock()
				if first {
					recorder.Record(username, startedAt)
				}
			}

			next(session)

			log.Printf(
				"SSH session disconnected remote=%s user=%q duration=%s",
				session.RemoteAddr(),
				loggedUsername,
				time.Since(startedAt).Round(time.Second),
			)
		}
	}
}

func initialSize(session ssh.Session) (int, int) {
	width, height := 80, 24
	if pty, _, ok := session.Pty(); ok {
		if pty.Window.Width > 0 {
			width = pty.Window.Width
		}
		if pty.Window.Height > 0 {
			height = pty.Window.Height
		}
	}
	return width, height
}

func terminalName(session ssh.Session) string {
	if pty, _, ok := session.Pty(); ok {
		return pty.Term
	}
	return ""
}
