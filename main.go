package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"

	"github.com/xalonious/portfolio-a-la-ssh/internal/ui"
)

const (
	defaultAddress = ":2323"
	defaultHostKey = "./host_ed25519"
)

func main() {
	log.SetOutput(os.Stdout)
	lipgloss.SetColorProfile(termenv.TrueColor)

	address := envOrDefault("SSH_PORTFOLIO_ADDR", defaultAddress)
	hostKeyPath := envOrDefault("SSH_PORTFOLIO_HOST_KEY", defaultHostKey)

	server, err := wish.NewServer(
		wish.WithAddress(address),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithMiddleware(
			bm.Middleware(func(session ssh.Session) (tea.Model, []tea.ProgramOption) {
				width, height := initialSize(session)
				return ui.New(width, height), []tea.ProgramOption{tea.WithAltScreen()}
			}),
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

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
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
