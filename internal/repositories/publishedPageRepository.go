package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-produk/internal/database"
	"go-produk/internal/models"
)

var ErrPublishedPageNotFound = errors.New("halaman publik tidak ditemukan")

type PublishedPageRepository struct{ pool *pgxpool.Pool }

func NewPublishedPageRepository(pool *pgxpool.Pool) *PublishedPageRepository {
	return &PublishedPageRepository{pool: pool}
}

const publishedPageColumns = `pp.id, pp.project_id, pp.user_id, u.username, pp.slug, pp.custom_path, pp.preview_token, pp.allow_preview, pp.is_active, pp.view_count, pp.published_at, pp.created_at, pp.updated_at`

func scanPublishedPage(row pgx.Row) (*models.PublishedPage, error) {
	page := &models.PublishedPage{}
	err := row.Scan(&page.ID, &page.ProjectID, &page.UserID, &page.Username, &page.Slug, &page.CustomPath, &page.PreviewToken, &page.AllowPreview, &page.IsActive, &page.ViewCount, &page.PublishedAt, &page.CreatedAt, &page.UpdatedAt)
	return page, err
}

func (r *PublishedPageRepository) GetOwnerProject(ctx context.Context, projectID, userID string) (*models.UserProject, error) {
	return NewUserProjectRepository(r.pool).GetByIDAndUserID(ctx, projectID, userID)
}

func (r *PublishedPageRepository) GetUsername(ctx context.Context, userID string) (string, error) {
	var username string
	err := r.pool.QueryRow(ctx, `SELECT COALESCE(username, '') FROM users WHERE id=$1 AND is_active=true AND is_deleted=false`, userID).Scan(&username)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrUserNotFound
	}
	return username, err
}

func (r *PublishedPageRepository) GetByProjectAndUserID(ctx context.Context, projectID, userID string) (*models.PublishedPage, error) {
	page, err := scanPublishedPage(r.pool.QueryRow(ctx, `SELECT `+publishedPageColumns+` FROM published_pages pp JOIN users u ON u.id = pp.user_id WHERE pp.project_id = $1 AND pp.user_id = $2`, projectID, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPublishedPageNotFound
	}
	return page, err
}

func (r *PublishedPageRepository) Upsert(ctx context.Context, page *models.PublishedPage) error {
	return r.pool.QueryRow(ctx, `INSERT INTO published_pages (project_id,user_id,slug,custom_path,preview_token,allow_preview) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (project_id) DO UPDATE SET slug=EXCLUDED.slug,custom_path=EXCLUDED.custom_path,allow_preview=EXCLUDED.allow_preview,updated_at=now() RETURNING id,preview_token,is_active,view_count,published_at,created_at,updated_at`, page.ProjectID, page.UserID, page.Slug, page.CustomPath, page.PreviewToken, page.AllowPreview).Scan(&page.ID, &page.PreviewToken, &page.IsActive, &page.ViewCount, &page.PublishedAt, &page.CreatedAt, &page.UpdatedAt)
}

func (r *PublishedPageRepository) Publish(ctx context.Context, projectID, userID string, active bool) error {
	return database.WithinTx(ctx, r.pool, func(ctx context.Context, tx pgx.Tx) error {
		var one int
		var err error
		if active {
			err = tx.QueryRow(ctx, `UPDATE published_pages SET is_active=true,published_at=now(),updated_at=now() WHERE project_id=$1 AND user_id=$2 RETURNING 1`, projectID, userID).Scan(&one)
		} else {
			err = tx.QueryRow(ctx, `UPDATE published_pages SET is_active=false,updated_at=now() WHERE project_id=$1 AND user_id=$2 RETURNING 1`, projectID, userID).Scan(&one)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPublishedPageNotFound
		}
		if err != nil {
			return err
		}
		status := "draft"
		if active {
			status = "published"
		}
		result, err := tx.Exec(ctx, `UPDATE user_projects SET status=$1,updated_at=now() WHERE id=$2 AND user_id=$3`, status, projectID, userID)
		if err != nil {
			return err
		}
		if result.RowsAffected() == 0 {
			return ErrUserProjectNotFound
		}
		return nil
	})
}

func (r *PublishedPageRepository) GetPublic(ctx context.Context, username, slug string, preview bool) (*models.PublishedPage, error) {
	query := `SELECT ` + publishedPageColumns + ` FROM published_pages pp JOIN users u ON u.id=pp.user_id JOIN user_projects up ON up.id=pp.project_id WHERE LOWER(u.username)=LOWER($1) AND LOWER(pp.slug)=LOWER($2) AND u.is_active=true AND u.is_deleted=false`
	if !preview {
		query += ` AND pp.is_active=true AND up.is_active=true AND up.status='published'`
	}
	page, err := scanPublishedPage(r.pool.QueryRow(ctx, query, username, slug))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPublishedPageNotFound
	}
	return page, err
}

func (r *PublishedPageRepository) IncrementViews(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE published_pages SET view_count=view_count+1,updated_at=now() WHERE id=$1`, id)
	return err
}
