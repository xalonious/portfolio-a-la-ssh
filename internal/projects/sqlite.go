package projects

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/xalonious/portfolio-a-la-ssh/internal/content"
	_ "modernc.org/sqlite"
)

const busyTimeoutMilliseconds = 5000

const publishedProjectsQuery = `
SELECT
	slug,
	title,
	description,
	image_path,
	image_alt,
	image_layout,
	technologies_json,
	repository_url,
	featured,
	sort_order,
	case_study_json,
	published_at,
	updated_at
FROM published_projects
ORDER BY sort_order ASC, title ASC`

type Logger interface {
	Printf(format string, v ...any)
}

type SQLiteRepository struct {
	databasePath     string
	caseStudyBaseURL string
	logger           Logger
}

func NewSQLiteRepository(databasePath, caseStudyBaseURL string, logger Logger) SQLiteRepository {
	return SQLiteRepository{
		databasePath:     databasePath,
		caseStudyBaseURL: strings.TrimRight(caseStudyBaseURL, "/"),
		logger:           logger,
	}
}

func (r SQLiteRepository) Load(ctx context.Context) ([]content.Project, error) {
	if strings.TrimSpace(r.databasePath) == "" {
		return nil, errors.New("PORTFOLIO_DATABASE_PATH is not configured")
	}

	dsn, err := readOnlyDSN(r.databasePath)
	if err != nil {
		return nil, fmt.Errorf("build read-only SQLite DSN: %w", err)
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open SQLite database: %w", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	rows, err := db.QueryContext(ctx, publishedProjectsQuery)
	if err != nil {
		return nil, fmt.Errorf("query published projects: %w", err)
	}
	defer rows.Close()

	var (
		loaded  []content.Project
		rowNum  int
		skipped int
	)
	for rows.Next() {
		rowNum++

		var row projectRow
		if err := rows.Scan(
			&row.slug,
			&row.title,
			&row.description,
			&row.imagePath,
			&row.imageAlt,
			&row.imageLayout,
			&row.technologiesJSON,
			&row.repositoryURL,
			&row.featured,
			&row.sortOrder,
			&row.caseStudyJSON,
			&row.publishedAt,
			&row.updatedAt,
		); err != nil {
			skipped++
			r.logMalformed(rowNum, "", fmt.Errorf("scan row: %w", err))
			continue
		}

		project, err := r.projectFromRow(row)
		if err != nil {
			skipped++
			r.logMalformed(rowNum, row.slug.String, err)
			continue
		}
		loaded = append(loaded, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read published project rows: %w", err)
	}
	if skipped > 0 && len(loaded) == 0 {
		return nil, fmt.Errorf("all %d published project rows were malformed", skipped)
	}

	return loaded, nil
}

type projectRow struct {
	slug             sql.NullString
	title            sql.NullString
	description      sql.NullString
	imagePath        sql.NullString
	imageAlt         sql.NullString
	imageLayout      sql.NullString
	technologiesJSON sql.NullString
	repositoryURL    sql.NullString
	featured         sql.NullInt64
	sortOrder        sql.NullInt64
	caseStudyJSON    sql.NullString
	publishedAt      sql.NullString
	updatedAt        sql.NullString
}

func (r SQLiteRepository) projectFromRow(row projectRow) (content.Project, error) {
	slug, err := requiredText("slug", row.slug)
	if err != nil {
		return content.Project{}, err
	}
	title, err := requiredText("title", row.title)
	if err != nil {
		return content.Project{}, err
	}
	description, err := requiredText("description", row.description)
	if err != nil {
		return content.Project{}, err
	}
	technologiesJSON, err := requiredText("technologies_json", row.technologiesJSON)
	if err != nil {
		return content.Project{}, err
	}

	var technologies []string
	if err := json.Unmarshal([]byte(technologiesJSON), &technologies); err != nil {
		return content.Project{}, fmt.Errorf("decode technologies_json: %w", err)
	}
	if technologies == nil {
		return content.Project{}, errors.New("technologies_json must be a JSON array")
	}
	for i, technology := range technologies {
		technologies[i] = strings.TrimSpace(technology)
		if technologies[i] == "" {
			return content.Project{}, fmt.Errorf("technologies_json contains an empty item at index %d", i)
		}
	}

	caseStudyURL := ""
	if row.caseStudyJSON.Valid {
		caseStudyJSON := bytes.TrimSpace([]byte(row.caseStudyJSON.String))
		if len(caseStudyJSON) > 0 && !bytes.Equal(caseStudyJSON, []byte("null")) {
			if !json.Valid(caseStudyJSON) {
				return content.Project{}, errors.New("case_study_json contains malformed JSON")
			}
			caseStudyURL = r.caseStudyBaseURL + "/" + url.PathEscape(slug)
		}
	}

	repositoryURL := ""
	if row.repositoryURL.Valid {
		repositoryURL = strings.TrimSpace(row.repositoryURL.String)
	}

	return content.Project{
		Title:       title,
		Description: description,
		Repo:        repositoryURL,
		CaseStudy:   caseStudyURL,
		Tech:        technologies,
	}, nil
}

func requiredText(name string, value sql.NullString) (string, error) {
	if !value.Valid || strings.TrimSpace(value.String) == "" {
		return "", fmt.Errorf("%s is null or empty", name)
	}
	return strings.TrimSpace(value.String), nil
}

func readOnlyDSN(databasePath string) (string, error) {
	absolutePath, err := filepath.Abs(databasePath)
	if err != nil {
		return "", err
	}

	databaseURL := url.URL{
		Scheme: "file",
		Opaque: filepath.ToSlash(absolutePath),
	}
	query := databaseURL.Query()
	query.Set("mode", "ro")
	query.Set("_pragma", fmt.Sprintf("busy_timeout(%d)", busyTimeoutMilliseconds))
	databaseURL.RawQuery = query.Encode()
	return databaseURL.String(), nil
}

func (r SQLiteRepository) logMalformed(rowNum int, slug string, err error) {
	if r.logger == nil {
		return
	}
	r.logger.Printf("skipping malformed published project row=%d slug=%q: %v", rowNum, slug, err)
}
