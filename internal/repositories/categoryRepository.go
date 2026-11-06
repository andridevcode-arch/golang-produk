package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-produk/internal/models"
)

var (
	ErrCategoryNotFound = errors.New("kategori tidak ditemukan")
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) Create(ctx context.Context, cat *models.Category) error {
	query := `
		INSERT INTO categories (parent_id, name, slug, icon, description, sort_order, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	now := time.Now()
	return r.pool.QueryRow(ctx, query,
		cat.ParentID, cat.Name, cat.Slug, cat.Icon, cat.Description, cat.SortOrder, cat.IsActive, now, now,
	).Scan(&cat.ID, &cat.CreatedAt, &cat.UpdatedAt)
}

func (r *CategoryRepository) GetAll(ctx context.Context, onlyActive bool) ([]models.Category, error) {
	query := `
		SELECT id, parent_id, name, slug, icon, description, sort_order, is_active, created_at, updated_at
		FROM categories
	`
	if onlyActive {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY sort_order ASC, name ASC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(
			&cat.ID, &cat.ParentID, &cat.Name, &cat.Slug, &cat.Icon,
			&cat.Description, &cat.SortOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*models.Category, error) {
	query := `
		SELECT id, parent_id, name, slug, icon, description, sort_order, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`
	var cat models.Category
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&cat.ID, &cat.ParentID, &cat.Name, &cat.Slug, &cat.Icon,
		&cat.Description, &cat.SortOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	query := `
		SELECT id, parent_id, name, slug, icon, description, sort_order, is_active, created_at, updated_at
		FROM categories
		WHERE slug = $1
	`
	var cat models.Category
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&cat.ID, &cat.ParentID, &cat.Name, &cat.Slug, &cat.Icon,
		&cat.Description, &cat.SortOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepository) Update(ctx context.Context, cat *models.Category) error {
	query := `
		UPDATE categories
		SET parent_id = $1, name = $2, slug = $3, icon = $4, description = $5,
		    sort_order = $6, is_active = $7, updated_at = $8
		WHERE id = $9
		RETURNING updated_at
	`
	cat.UpdatedAt = time.Now()
	err := r.pool.QueryRow(ctx, query,
		cat.ParentID, cat.Name, cat.Slug, cat.Icon, cat.Description,
		cat.SortOrder, cat.IsActive, cat.UpdatedAt, cat.ID,
	).Scan(&cat.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return err
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}
	return nil
}
