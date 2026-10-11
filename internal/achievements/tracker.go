package achievements

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type Logger interface {
	Printf(format string, v ...any)
}

type connection struct {
	ipHash      string
	connectedAt int64
}


type Tracker struct {
	db     *sql.DB
	key    string
	logger Logger
	queue  chan connection
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewTracker(databasePath, key string, logger Logger) *Tracker {
	if strings.TrimSpace(databasePath) == "" || len([]byte(key)) < 32 {
		if logger != nil {
			logger.Printf("SSH achievement tracking disabled: configure PORTFOLIO_DATABASE_PATH and SSH_ACHIEVEMENT_IP_HASH_KEY (at least 32 bytes)")
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
	tracker := &Tracker{db: db, key: key, logger: logger, queue: make(chan connection, 64), ctx: ctx, cancel: cancel}
	tracker.wg.Add(1)
	go tracker.run()
	return tracker
}

func NormalizeRemoteIP(remote net.Addr) (string, error) {
	if remote == nil {
		return "", fmt.Errorf("missing remote address")
	}
	host, _, err := net.SplitHostPort(remote.String())
	if err != nil {
		return "", fmt.Errorf("invalid remote address")
	}
	ip, err := netip.ParseAddr(host)
	if err != nil {
		return "", fmt.Errorf("invalid remote IP")
	}
	return ip.Unmap().WithZone("").String(), nil
}

func HashIP(ip, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	_, _ = hash.Write([]byte(ip))
	return hex.EncodeToString(hash.Sum(nil))
}

func (t *Tracker) Record(remote net.Addr, connectedAt time.Time) {
	if t == nil {
		return
	}
	ip, err := NormalizeRemoteIP(remote)
	if err != nil {
		t.warn(err)
		return
	}
	record := connection{ipHash: HashIP(ip, t.key), connectedAt: connectedAt.UnixMilli()}
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
		"INSERT INTO ssh_achievement_connections (ip_hash, connected_at) VALUES (?, ?)",
		record.ipHash, record.connectedAt,
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
	// Do not create a new database if the shared path is missing/misconfigured.
	query.Set("mode", "rw")
	query.Set("_pragma", "busy_timeout(5000)")
	databaseURL.RawQuery = query.Encode()
	return databaseURL.String(), nil
}
