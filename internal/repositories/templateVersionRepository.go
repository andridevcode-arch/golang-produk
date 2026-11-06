package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-produk/internal/models"
)

var ErrTemplateVersionNotFound = errors.New("versi template tidak ditemukan")

type TemplateVersionRepository struct{ pool *pgxpool.Pool }

func NewTemplateVersionRepository(pool *pgxpool.Pool) *TemplateVersionRepository {
	return &TemplateVersionRepository{pool: pool}
}

func scanTemplateVersion(row pgx.Row) (*models.TemplateVersion, error) {
	var v models.TemplateVersion
	var template models.Template
	var templateID, templateProductID, templateName, templateSlug, templatePreview string
	var templateThumbnail *string
	var templatePremium, templateActive bool
	var templateSort int
	var templateConfig []byte
	var templateCreated, templateUpdated time.Time
	if err := row.Scan(&v.ID, &v.TemplateID, &v.Version, &v.ConfigSchema, &v.IsActive, &v.CreatedAt,
		&templateID, &templateProductID, &templateName, &templateSlug, &templatePreview, &templateThumbnail,
		&templatePremium, &templateActive, &templateSort, &templateConfig, &templateCreated, &templateUpdated); err != nil {
		return nil, err
	}
	template.ID = templateID
	template.ProductID = templateProductID
	template.Name = templateName
	template.Slug = templateSlug
	template.PreviewURL = templatePreview
	template.ThumbnailURL = templateThumbnail
	template.IsPremium = templatePremium
	template.IsActive = templateActive
	template.SortOrder = templateSort
	template.ConfigSchema = templateConfig
	template.CreatedAt = templateCreated
	template.UpdatedAt = templateUpdated
	v.Template = &template
	return &v, nil
}

const templateVersionSelect = `
	v.id, v.template_id, v.version, v.config_schema, v.is_active, v.created_at,
	t.id, t.product_id, t.name, t.slug, t.preview_url, t.thumbnail_url,
	t.is_premium, t.is_active, t.sort_order, t.config_schema, t.created_at, t.updated_at
`

func (r *TemplateVersionRepository) Create(ctx context.Context, v *models.TemplateVersion) error {
	query := `INSERT INTO template_versions (template_id, version, config_schema, is_active) VALUES ($1, $2, $3, false) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, v.TemplateID, v.Version, v.ConfigSchema).Scan(&v.ID, &v.CreatedAt)
}

func (r *TemplateVersionRepository) GetAll(ctx context.Context, templateID string) ([]models.TemplateVersion, error) {
	query := "SELECT " + templateVersionSelect + " FROM template_versions v JOIN templates t ON t.id = v.template_id WHERE v.template_id = $1 ORDER BY v.created_at DESC"
	rows, err := r.pool.Query(ctx, query, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []models.TemplateVersion{}
	for rows.Next() {
		v, err := scanTemplateVersion(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *v)
	}
	return result, rows.Err()
}

func (r *TemplateVersionRepository) GetByID(ctx context.Context, templateID, id string) (*models.TemplateVersion, error) {
	query := "SELECT " + templateVersionSelect + " FROM template_versions v JOIN templates t ON t.id = v.template_id WHERE v.template_id = $1 AND v.id = $2"
	v, err := scanTemplateVersion(r.pool.QueryRow(ctx, query, templateID, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTemplateVersionNotFound
	}
	return v, err
}

func (r *TemplateVersionRepository) GetByVersion(ctx context.Context, templateID, version string) (*models.TemplateVersion, error) {
	query := "SELECT " + templateVersionSelect + " FROM template_versions v JOIN templates t ON t.id = v.template_id WHERE v.template_id = $1 AND v.version = $2"
	v, err := scanTemplateVersion(r.pool.QueryRow(ctx, query, templateID, version))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTemplateVersionNotFound
	}
	return v, err
}

func (r *TemplateVersionRepository) SetActive(ctx context.Context, templateID, id string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, "UPDATE template_versions SET is_active = false WHERE template_id = $1", templateID); err != nil {
		return err
	}
	result, err := tx.Exec(ctx, "UPDATE template_versions SET is_active = true WHERE template_id = $1 AND id = $2", templateID, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrTemplateVersionNotFound
	}
	return tx.Commit(ctx)
}

func (r *TemplateVersionRepository) Delete(ctx context.Context, templateID, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM template_versions WHERE template_id = $1 AND id = $2", templateID, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrTemplateVersionNotFound
	}
	return nil
}
