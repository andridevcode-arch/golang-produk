package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-produk/internal/database"
	"go-produk/internal/models"
)

var ErrUserProjectNotFound = errors.New("proyek tidak ditemukan")

type UserProjectRepository struct {
	pool *pgxpool.Pool
}

func NewUserProjectRepository(pool *pgxpool.Pool) *UserProjectRepository {
	return &UserProjectRepository{pool: pool}
}

const userProjectColumns = `
	up.id, up.user_id, up.template_id, up.template_version_id, up.title, up.description,
	up.status, up.is_active, up.created_at, up.updated_at,
	t.id, t.product_id, t.name, t.slug, t.preview_url, t.thumbnail_url,
	t.is_premium, t.is_active, t.sort_order, t.config_schema, t.created_at, t.updated_at,
	v.id, v.template_id, v.version, v.config_schema, v.is_active, v.created_at,
	p.id, p.category_id, p.name, p.slug, p.description, p.price,
	p.access_type, p.is_active, p.sort_order, p.created_at, p.updated_at
`

func scanUserProject(row pgx.Row) (*models.UserProject, error) {
	var project models.UserProject
	var template models.Template
	var version models.TemplateVersion
	var product models.Product

	err := row.Scan(
		&project.ID, &project.UserID, &project.TemplateID, &project.TemplateVersionID, &project.Title, &project.Description,
		&project.Status, &project.IsActive, &project.CreatedAt, &project.UpdatedAt,
		&template.ID, &template.ProductID, &template.Name, &template.Slug, &template.PreviewURL, &template.ThumbnailURL,
		&template.IsPremium, &template.IsActive, &template.SortOrder, &template.ConfigSchema, &template.CreatedAt, &template.UpdatedAt,
		&version.ID, &version.TemplateID, &version.Version, &version.ConfigSchema, &version.IsActive, &version.CreatedAt,
		&product.ID, &product.CategoryID, &product.Name, &product.Slug, &product.Description, &product.Price,
		&product.AccessType, &product.IsActive, &product.SortOrder, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	template.Products = []models.Product{product}
	version.Template = &template
	project.Template = &template
	project.TemplateVersion = &version
	return &project, nil
}

func (r *UserProjectRepository) Create(ctx context.Context, project *models.UserProject) error {
	query := `
		INSERT INTO user_projects (user_id, template_id, template_version_id, title, description, status, is_active)
		VALUES ($1, $2, $3, $4, $5, 'draft', true)
		RETURNING id, status, is_active, created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		project.UserID, project.TemplateID, project.TemplateVersionID, project.Title, project.Description,
	).Scan(&project.ID, &project.Status, &project.IsActive, &project.CreatedAt, &project.UpdatedAt)
}

func (r *UserProjectRepository) ListByUserID(ctx context.Context, userID, status string) ([]models.UserProject, error) {
	query := "SELECT " + userProjectColumns + `
		FROM user_projects up
		JOIN templates t ON t.id = up.template_id
		JOIN template_versions v ON v.id = up.template_version_id
		JOIN products p ON p.id = t.product_id
		WHERE up.user_id = $1`
	args := []any{userID}
	if status != "" {
		args = append(args, status)
		query += fmt.Sprintf(" AND up.status = $%d", len(args))
	}
	query += " ORDER BY up.created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []models.UserProject{}
	for rows.Next() {
		project, err := scanUserProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *project)
	}
	return projects, rows.Err()
}

func (r *UserProjectRepository) GetByIDAndUserID(ctx context.Context, id, userID string) (*models.UserProject, error) {
	query := "SELECT " + userProjectColumns + `
		FROM user_projects up
		JOIN templates t ON t.id = up.template_id
		JOIN template_versions v ON v.id = up.template_version_id
		JOIN products p ON p.id = t.product_id
		WHERE up.id = $1 AND up.user_id = $2`
	project, err := scanUserProject(r.pool.QueryRow(ctx, query, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserProjectNotFound
	}
	return project, err
}

func (r *UserProjectRepository) UpdateByIDAndUserID(ctx context.Context, project *models.UserProject) error {
	query := `
		UPDATE user_projects
		SET title = $1, description = $2, updated_at = $3
		WHERE id = $4 AND user_id = $5
		RETURNING updated_at`
	project.UpdatedAt = time.Now()
	err := r.pool.QueryRow(ctx, query, project.Title, project.Description, project.UpdatedAt, project.ID, project.UserID).Scan(&project.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserProjectNotFound
	}
	return err
}

func (r *UserProjectRepository) SetStatusByIDAndUserID(ctx context.Context, id, userID, status string) error {
	result, err := r.pool.Exec(ctx,
		"UPDATE user_projects SET status = $1, updated_at = now() WHERE id = $2 AND user_id = $3",
		status, id, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrUserProjectNotFound
	}
	return nil
}

func (r *UserProjectRepository) ArchiveByIDAndUserID(ctx context.Context, id, userID string) error {
	return database.WithinTx(ctx, r.pool, func(ctx context.Context, tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE published_pages SET is_active=false,updated_at=now() WHERE project_id=$1 AND user_id=$2`, id, userID); err != nil {
			return err
		}
		result, err := tx.Exec(ctx, `UPDATE user_projects SET status='archived',updated_at=now() WHERE id=$1 AND user_id=$2`, id, userID)
		if err != nil {
			return err
		}
		if result.RowsAffected() == 0 {
			return ErrUserProjectNotFound
		}
		return nil
	})
}

func (r *UserProjectRepository) GetActiveTemplateVersionID(ctx context.Context, templateID string) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		"SELECT id FROM template_versions WHERE template_id = $1 AND is_active = true",
		templateID,
	).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrTemplateVersionNotFound
	}
	return id, err
}

func (r *UserProjectRepository) IsTemplateVersionOwnedByTemplate(ctx context.Context, templateID, versionID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM template_versions WHERE id = $1 AND template_id = $2)",
		versionID, templateID,
	).Scan(&exists)
	return exists, err
}
