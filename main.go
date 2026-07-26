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

	"github.com/xalonious/portfolio-a-la-ssh/internal/content"
	"github.com/xalonious/portfolio-a-la-ssh/internal/projects"
	"github.com/xalonious/portfolio-a-la-ssh/internal/ui"
)

const (
	defaultAddress      = ":2323"
	defaultHostKey      = "./host_ed25519"
	defaultDatabasePath = "D:/coding/portfolio/.data/portfolio.db"
)

func main() {
	log.SetOutput(os.Stdout)
	loadDotEnv()
	lipgloss.SetColorProfile(termenv.TrueColor)

	address := envOrDefault("SSH_PORTFOLIO_ADDR", defaultAddress)
	hostKeyPath := envOrDefault("SSH_PORTFOLIO_HOST_KEY", defaultHostKey)
	databasePath := envOrDefault("PORTFOLIO_DATABASE_PATH", defaultDatabasePath)
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
			connectionLogger,
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
	paths := []string{".env"}
	if executablePath, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(executablePath), ".env"))
	}

	seen := make(map[string]bool, len(paths))
	for _, path := range paths {
		absolutePath, err := filepath.Abs(path)
		if err != nil || seen[absolutePath] {
			continue
		}
		seen[absolutePath] = true

		if err := godotenv.Load(absolutePath); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Printf("could not load environment file path=%q: %v", absolutePath, err)
		}
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func connectionLogger(next ssh.Handler) ssh.Handler {
	return func(session ssh.Session) {
		startedAt := time.Now()
		term := terminalName(session)

		log.Printf(
			"SSH session connected remote=%s user=%q term=%q",
			session.RemoteAddr(),
			session.User(),
			term,
		)

		next(session)

		log.Printf(
			"SSH session disconnected remote=%s user=%q duration=%s",
			session.RemoteAddr(),
			session.User(),
			time.Since(startedAt).Round(time.Second),
		)
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
