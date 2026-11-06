package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-produk/internal/models"
)

var ErrTemplateNotFound = errors.New("template tidak ditemukan")

type TemplateRepository struct {
	pool *pgxpool.Pool
}

func NewTemplateRepository(pool *pgxpool.Pool) *TemplateRepository {
	return &TemplateRepository{pool: pool}
}

const templateColumns = `
	t.id, t.product_id, t.name, t.slug, t.preview_url, t.thumbnail_url,
	t.is_premium, t.is_active, t.sort_order, t.config_schema, t.created_at, t.updated_at,
	p.id, p.category_id, p.name, p.slug, p.description, p.price,
	p.access_type, p.is_active, p.sort_order, p.created_at, p.updated_at`

func scanTemplate(row pgx.Row) (*models.Template, error) {
	var t models.Template
	var product models.Product
	var productID, productName, productSlug, accessType *string
	var description *string
	var price *float64
	var productActive *bool
	var productSort *int
	var productCreated, productUpdated *time.Time
	var productCategoryID *string

	err := row.Scan(
		&t.ID, &t.ProductID, &t.Name, &t.Slug, &t.PreviewURL, &t.ThumbnailURL,
		&t.IsPremium, &t.IsActive, &t.SortOrder, &t.ConfigSchema, &t.CreatedAt, &t.UpdatedAt,
		&productID, &productCategoryID, &productName, &productSlug, &description, &price,
		&accessType, &productActive, &productSort, &productCreated, &productUpdated,
	)
	if err != nil {
		return nil, err
	}
	if productID != nil {
		product.ID = *productID
		product.CategoryID = valueOrEmpty(productCategoryID)
		product.Name = valueOrEmpty(productName)
		product.Slug = valueOrEmpty(productSlug)
		product.Description = description
		product.Price = price
		product.AccessType = valueOrEmpty(accessType)
		product.IsActive = productActive != nil && *productActive
		if productSort != nil {
			product.SortOrder = *productSort
		}
		if productCreated != nil {
			product.CreatedAt = *productCreated
		}
		if productUpdated != nil {
			product.UpdatedAt = *productUpdated
		}
		t.Products = []models.Product{product}
	}
	return &t, nil
}

func valueOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func (r *TemplateRepository) Create(ctx context.Context, t *models.Template) error {
	query := `
		INSERT INTO templates (product_id, name, slug, preview_url, thumbnail_url, is_premium, is_active, sort_order, config_schema, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, created_at, updated_at`
	now := time.Now()
	return r.pool.QueryRow(ctx, query, t.ProductID, t.Name, t.Slug, t.PreviewURL, t.ThumbnailURL,
		t.IsPremium, t.IsActive, t.SortOrder, t.ConfigSchema, now, now).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TemplateRepository) GetAll(ctx context.Context, productID string, premium *bool, onlyActive bool) ([]models.Template, error) {
	query := "SELECT " + templateColumns + " FROM templates t JOIN products p ON p.id = t.product_id"
	args := []any{}
	conditions := []string{}
	if productID != "" {
		args = append(args, productID)
		conditions = append(conditions, fmt.Sprintf("t.product_id = $%d", len(args)))
	}
	if premium != nil {
		args = append(args, *premium)
		conditions = append(conditions, fmt.Sprintf("t.is_premium = $%d", len(args)))
	}
	if onlyActive {
		conditions = append(conditions, "t.is_active = true")
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY t.sort_order ASC, t.name ASC"
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Template
	for rows.Next() {
		t, err := scanTemplate(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *t)
	}
	return result, rows.Err()
}

func (r *TemplateRepository) GetByID(ctx context.Context, id string) (*models.Template, error) {
	query := "SELECT " + templateColumns + " FROM templates t JOIN products p ON p.id = t.product_id WHERE t.id = $1"
	t, err := scanTemplate(r.pool.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTemplateNotFound
	}
	return t, err
}

func (r *TemplateRepository) GetBySlug(ctx context.Context, productID, slug string) (*models.Template, error) {
	query := "SELECT " + templateColumns + " FROM templates t JOIN products p ON p.id = t.product_id WHERE t.product_id = $1 AND t.slug = $2"
	t, err := scanTemplate(r.pool.QueryRow(ctx, query, productID, slug))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTemplateNotFound
	}
	return t, err
}

func (r *TemplateRepository) Update(ctx context.Context, t *models.Template) error {
	t.UpdatedAt = time.Now()
	query := `UPDATE templates SET name=$1, slug=$2, preview_url=$3, thumbnail_url=$4, is_premium=$5, is_active=$6, sort_order=$7, config_schema=$8, updated_at=$9 WHERE id=$10 RETURNING updated_at`
	err := r.pool.QueryRow(ctx, query, t.Name, t.Slug, t.PreviewURL, t.ThumbnailURL, t.IsPremium, t.IsActive, t.SortOrder, t.ConfigSchema, t.UpdatedAt, t.ID).Scan(&t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTemplateNotFound
	}
	return err
}

func (r *TemplateRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM templates WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrTemplateNotFound
	}
	return nil
}
