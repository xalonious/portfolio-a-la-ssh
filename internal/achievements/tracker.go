package achievements

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type Logger interface {
	Printf(format string, v ...any)
}

type connection struct {
	tokenHash   string
	connectedAt int64
}

type Tracker struct {
	db     *sql.DB
	logger Logger
	queue  chan connection
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewTracker(databasePath string, logger Logger) *Tracker {
	if strings.TrimSpace(databasePath) == "" {
		if logger != nil {
			logger.Printf("SSH achievement tracking disabled: configure PORTFOLIO_DATABASE_PATH")
		}
		return nil
	}
	dsn, err := writableDSN(databasePath)
	if err != nil {
		if logger != nil {
			logger.Printf("SSH achievement tracking unavailable: %v", err)
		}
		return nil
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		if logger != nil {
			logger.Printf("SSH achievement tracking unavailable: %v", err)
		}
		return nil
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	ctx, cancel := context.WithCancel(context.Background())
	tracker := &Tracker{db: db, logger: logger, queue: make(chan connection, 64), ctx: ctx, cancel: cancel}
	tracker.wg.Add(1)
	go tracker.run()
	return tracker
}

var visitorTokenPattern = regexp.MustCompile(`^x_[0-9a-f]{6}$`)

func IsVisitorToken(username string) bool {
	return len(username) == 8 && visitorTokenPattern.MatchString(username)
}

func HashVisitorToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (t *Tracker) Record(token string, connectedAt time.Time) {
	if t == nil || !IsVisitorToken(token) {
		return
	}
	record := connection{tokenHash: HashVisitorToken(token), connectedAt: connectedAt.UnixMilli()}
	select {
	case <-t.ctx.Done():
	case t.queue <- record:
	default:
		t.warn(fmt.Errorf("write queue full; connection not recorded"))
	}
}

func (t *Tracker) run() {
	defer t.wg.Done()
	for {
		select {
		case <-t.ctx.Done():
			return
		case record := <-t.queue:
			ctx, cancel := context.WithTimeout(t.ctx, 6*time.Second)
			err := t.insert(ctx, record)
			cancel()
			if err != nil {
				t.warn(err)
			}
		}
	}
}

func (t *Tracker) insert(ctx context.Context, record connection) error {
	_, err := t.db.ExecContext(ctx,
		"INSERT INTO ssh_achievement_connections (token_hash, connected_at) VALUES (?, ?)",
		record.tokenHash, record.connectedAt,
	)
	return err
}

func (t *Tracker) warn(err error) {
	if t.logger != nil {
		t.logger.Printf("SSH achievement tracking unavailable: %v", err)
	}
}

func (t *Tracker) Close() {
	if t == nil {
		return
	}
	t.cancel()
	t.wg.Wait()
	_ = t.db.Close()
}

func writableDSN(databasePath string) (string, error) {
	absolutePath, err := filepath.Abs(databasePath)
	if err != nil {
		return "", err
	}
	databaseURL := url.URL{Scheme: "file", Opaque: filepath.ToSlash(absolutePath)}
	query := databaseURL.Query()
	query.Set("mode", "rw")
	query.Set("_pragma", "busy_timeout(5000)")
	databaseURL.RawQuery = query.Encode()
	return databaseURL.String(), nil
}
