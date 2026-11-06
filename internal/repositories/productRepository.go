package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-produk/internal/models"
)

var (
	ErrProductNotFound = errors.New("produk tidak ditemukan")
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, p *models.Product) error {
	query := `
		INSERT INTO products (category_id, name, slug, description, price, access_type, is_active, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		p.CategoryID, p.Name, p.Slug, p.Description, p.Price, p.AccessType, p.IsActive, p.SortOrder,
		p.CreatedAt, p.UpdatedAt,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProductRepository) GetAll(ctx context.Context, categoryID, accessType string, onlyActive bool) ([]models.Product, error) {
	query := `
		SELECT p.id, p.category_id, p.name, p.slug, p.description, p.price,
		       p.access_type, p.is_active, p.sort_order, p.created_at, p.updated_at,
		       c.id as c_id, c.parent_id, c.name as c_name, c.slug as c_slug,
		       c.icon, c.description as c_description, c.is_active as c_is_active,
		       c.sort_order as c_sort_order, c.created_at as c_created_at, c.updated_at as c_updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
	`
	args := []interface{}{}
	conditions := []string{}

	if categoryID != "" {
		conditions = append(conditions, "p.category_id = $"+string(rune('1'+len(conditions))))
		args = append(args, categoryID)
	}
	if accessType != "" {
		conditions = append(conditions, "p.access_type = $"+string(rune('1'+len(conditions))))
		args = append(args, accessType)
	}
	if onlyActive {
		conditions = append(conditions, "p.is_active = true")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY p.sort_order ASC, p.name ASC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var catID, catName, catSlug string
		var catParentID, catIcon, catDescription *string
		var catIsActive bool
		var catSortOrder *int
		var catCreatedAt, catUpdatedAt *time.Time

		err := rows.Scan(
			&p.ID, &p.CategoryID, &p.Name, &p.Slug, &p.Description, &p.Price,
			&p.AccessType, &p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
			&catID, &catParentID, &catName, &catSlug,
			&catIcon, &catDescription, &catIsActive,
			&catSortOrder, &catCreatedAt, &catUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if catID != "" {
			p.Categories = []models.Category{{
				ID:          catID,
				ParentID:    catParentID,
				Name:        catName,
				Slug:        catSlug,
				Icon:        catIcon,
				Description: catDescription,
				IsActive:    catIsActive,
			}}
			if catSortOrder != nil {
				p.Categories[0].SortOrder = *catSortOrder
			}
			if catCreatedAt != nil {
				p.Categories[0].CreatedAt = *catCreatedAt
			}
			if catUpdatedAt != nil {
				p.Categories[0].UpdatedAt = *catUpdatedAt
			}
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	query := `
		SELECT p.id, p.category_id, p.name, p.slug, p.description, p.price,
		       p.access_type, p.is_active, p.sort_order, p.created_at, p.updated_at,
		       c.id as c_id, c.parent_id, c.name as c_name, c.slug as c_slug,
		       c.icon, c.description as c_description, c.is_active as c_is_active,
		       c.sort_order as c_sort_order, c.created_at as c_created_at, c.updated_at as c_updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`
	var p models.Product
	var catID, catName, catSlug string
	var catParentID, catIcon, catDescription *string
	var catIsActive bool
	var catSortOrder *int
	var catCreatedAt, catUpdatedAt *time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.CategoryID, &p.Name, &p.Slug, &p.Description, &p.Price,
		&p.AccessType, &p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		&catID, &catParentID, &catName, &catSlug,
		&catIcon, &catDescription, &catIsActive,
		&catSortOrder, &catCreatedAt, &catUpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if catID != "" {
		p.Categories = []models.Category{{
			ID:          catID,
			ParentID:    catParentID,
			Name:        catName,
			Slug:        catSlug,
			Icon:        catIcon,
			Description: catDescription,
			IsActive:    catIsActive,
		}}
		if catSortOrder != nil {
			p.Categories[0].SortOrder = *catSortOrder
		}
		if catCreatedAt != nil {
			p.Categories[0].CreatedAt = *catCreatedAt
		}
		if catUpdatedAt != nil {
			p.Categories[0].UpdatedAt = *catUpdatedAt
		}
	}

	return &p, nil
}

func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	query := `
		SELECT p.id, p.category_id, p.name, p.slug, p.description, p.price,
		       p.access_type, p.is_active, p.sort_order, p.created_at, p.updated_at,
		       c.id as c_id, c.parent_id, c.name as c_name, c.slug as c_slug,
		       c.icon, c.description as c_description, c.is_active as c_is_active,
		       c.sort_order as c_sort_order, c.created_at as c_created_at, c.updated_at as c_updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.slug = $1
	`
	var p models.Product
	var catID, catName, catSlug string
	var catParentID, catIcon, catDescription *string
	var catIsActive bool
	var catSortOrder *int
	var catCreatedAt, catUpdatedAt *time.Time

	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.CategoryID, &p.Name, &p.Slug, &p.Description, &p.Price,
		&p.AccessType, &p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		&catID, &catParentID, &catName, &catSlug,
		&catIcon, &catDescription, &catIsActive,
		&catSortOrder, &catCreatedAt, &catUpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if catID != "" {
		p.Categories = []models.Category{{
			ID:          catID,
			ParentID:    catParentID,
			Name:        catName,
			Slug:        catSlug,
			Icon:        catIcon,
			Description: catDescription,
			IsActive:    catIsActive,
		}}
		if catSortOrder != nil {
			p.Categories[0].SortOrder = *catSortOrder
		}
		if catCreatedAt != nil {
			p.Categories[0].CreatedAt = *catCreatedAt
		}
		if catUpdatedAt != nil {
			p.Categories[0].UpdatedAt = *catUpdatedAt
		}
	}

	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *models.Product) error {
	query := `
		UPDATE products
		SET category_id = $1, name = $2, slug = $3, description = $4, price = $5,
		    access_type = $6, is_active = $7, sort_order = $8, updated_at = $9
		WHERE id = $10
		RETURNING updated_at
	`
	p.UpdatedAt = time.Now()
	err := r.pool.QueryRow(ctx, query,
		p.CategoryID, p.Name, p.Slug, p.Description, p.Price,
		p.AccessType, p.IsActive, p.SortOrder, p.UpdatedAt, p.ID,
	).Scan(&p.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrProductNotFound
	}
	return nil
}
