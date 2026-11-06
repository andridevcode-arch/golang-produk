package app

import (
	"context"
	"fmt"
	"log"

	"go-produk/internal/config"
	"go-produk/internal/database"
	"go-produk/internal/repositories"
	"go-produk/internal/services"
	"go-produk/pkg/token"
)

// App menampung semua dependencies
type App struct {
	Config                 *config.Config
	Supabase               *database.Supabase
	Redis                  *database.RedisClient
	Mongo                  *database.Mongo
	AuthService            *services.AuthService
	UserService            *services.UserService
	CategoryService        *services.CategoryService
	ProductService         *services.ProductService
	TemplateService        *services.TemplateService
	TemplateVersionService *services.TemplateVersionService
	UserProjectService     *services.UserProjectService
	PublishedPageService   *services.PublishedPageService
	TokenMgr               *token.Manager
}

// New membuat instance App dan menginisialisasi semua koneksi database
func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	app := &App{Config: cfg}

	// Inisialisasi Supabase (optional)
	app.Supabase, err = database.NewSupabase(ctx, cfg.Postgres)
	if err != nil {
		log.Printf("warning: postgres tidak tersedia (%v), beberapa fitur mungkin terbatas", err)
		app.Supabase = nil
	}

	// Inisialisasi Redis (optional)
	app.Redis, err = database.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Printf("warning: redis tidak tersedia (%v), beberapa fitur mungkin terbatas", err)
		app.Redis = nil
	}

	// Inisialisasi MongoDB (optional)
	app.Mongo, err = database.NewMongo(ctx, cfg.Mongo)
	if err != nil {
		log.Printf("warning: mongo tidak tersedia (%v), beberapa fitur mungkin terbatas", err)
		app.Mongo = nil
	}

	// Inisialisasi Auth, User, Category, dan Product Services (memerlukan Supabase)
	if app.Supabase != nil {
		tokenMgr := token.NewManager(cfg.JWTSecret, cfg.JWTExpiresIn)
		userRepo := repositories.NewUserRepository(app.Supabase.Pool())
		categoryRepo := repositories.NewCategoryRepository(app.Supabase.Pool())
		productRepo := repositories.NewProductRepository(app.Supabase.Pool())
		templateRepo := repositories.NewTemplateRepository(app.Supabase.Pool())
		templateVersionRepo := repositories.NewTemplateVersionRepository(app.Supabase.Pool())
		userProjectRepo := repositories.NewUserProjectRepository(app.Supabase.Pool())
		publishedPageRepo := repositories.NewPublishedPageRepository(app.Supabase.Pool())

		app.AuthService = services.NewAuthService(userRepo, tokenMgr, cfg.JWTExpiresIn, cfg.SupabaseURL, cfg.SupabaseServiceRoleKey)
		app.UserService = services.NewUserService(userRepo)
		app.CategoryService = services.NewCategoryService(categoryRepo)
		app.ProductService = services.NewProductService(productRepo)
		app.TemplateService = services.NewTemplateService(templateRepo, productRepo)
		app.TemplateVersionService = services.NewTemplateVersionService(templateVersionRepo, templateRepo)
		app.UserProjectService = services.NewUserProjectService(userProjectRepo, templateRepo)
		app.PublishedPageService = services.NewPublishedPageService(publishedPageRepo)
		app.TokenMgr = tokenMgr
	} else {
		log.Printf("warning: auth dan product service tidak tersedia karena postgres tidak tersedia")
	}

	log.Println("aplikasi berhasil diinisialisasi")
	return app, nil
}

// Shutdown menutup semua koneksi database
func (a *App) Shutdown(ctx context.Context) {
	if a.Supabase != nil {
		a.Supabase.Close()
	}
	if a.Redis != nil {
		_ = a.Redis.Close()
	}
	if a.Mongo != nil {
		_ = a.Mongo.Close(ctx)
	}
	log.Println("semua koneksi database ditutup")
}
