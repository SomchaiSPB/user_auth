package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/SomchaiSPB/user-auth/internal/config"
	"github.com/SomchaiSPB/user-auth/internal/custom_middleware"
	"github.com/SomchaiSPB/user-auth/internal/entity"
	"github.com/SomchaiSPB/user-auth/internal/hash"
	"github.com/SomchaiSPB/user-auth/internal/logger"
	"github.com/SomchaiSPB/user-auth/internal/repository"
	"github.com/SomchaiSPB/user-auth/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jaswdr/faker"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	sqliteStorage   = "sqlite"
	postgresStorage = "postgres"
)

var (
	ErrStorageTypeNotFound = errors.New("storage type not found")
	ErrSqliteConnect       = errors.New("connecting to sqlite error")
	ErrCreateSqliteDbFile  = errors.New("creating db file error")
	ErrPostgresConnect     = errors.New("connecting to postgres error")
	ErrDBMigration         = errors.New("migrating entities error")
	ErrAddFixtures         = errors.New("adding fixtures error")
)

type App struct {
	config     *config.Config
	logger     logger.Logger
	db         *gorm.DB
	userSvc    *service.UserService
	productSvc *service.ProductService
	hasher     hash.Hasher
}

func New(config *config.Config) *App {
	return &App{config: config}
}

func (a *App) Init() error {
	l, err := logger.NewSugaredLogger()

	if err != nil {
		return err
	}

	a.logger = l

	a.hasher = hash.NewHasher()

	if err := a.initDB(); err != nil {
		return err
	}

	if a.config.WithTableTruncate() {
		if err := a.db.Migrator().DropTable(&entity.User{}); err != nil {
			log.Println("error dropping users table")
		}
		if err := a.db.Migrator().DropTable(&entity.Product{}); err != nil {
			log.Println("error dropping products table")
		}
	}

	if err := a.db.AutoMigrate(&entity.Product{}, &entity.User{}); err != nil {
		return fmt.Errorf("%w: %w", ErrDBMigration, err)
	}

	if a.config.WithFakeData() {
		if err := a.addFixtures(); err != nil {
			return fmt.Errorf("%w: %w", ErrAddFixtures, err)
		}
	}

	a.userSvc = service.NewUserSvc(repository.NewUserDBRepository(a.db))
	a.productSvc = service.NewProductSvc(repository.NewProductDBRepository(a.db))

	return nil
}

func (a *App) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go a.startServer(ctx, wg)
}

func (a *App) ShutDown() error {
	return nil
}

func (a *App) startServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	a.logger.Infof("server started at :" + a.config.HttpPort())

	go http.ListenAndServe(":"+a.config.HttpPort(), a.router())

	<-ctx.Done()
}

func (a *App) router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(custom_middleware.ContentTypeJson)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", a.HandleCreateUser)
		r.Post("/sign-in", a.HandleAuthUser)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(a.ApiTokenMiddleware)
		r.Get("/product", a.HandleGetProduct)
		r.Get("/products", a.HandleGetProducts)
	})

	return r
}

func (a *App) initDB() error {
	var db *gorm.DB
	var err error

	switch a.config.Storage() {
	case sqliteStorage:
		dbPath := fmt.Sprintf("%s/%s", "db_data", a.config.SqliteDBConfig.DbFile())

		err = os.MkdirAll("db_data", os.ModePerm)

		if err != nil {
			return fmt.Errorf("%w: %w", ErrCreateSqliteDbFile, err)
		}

		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

		if err != nil {
			return fmt.Errorf("%w:%w", ErrSqliteConnect, err)
		}
	case postgresStorage:
		db, err = gorm.Open(postgres.Open(a.config.PostgresDBConfig.Dsn()), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("%w:%w", ErrPostgresConnect, err)
		}
	default:
		return fmt.Errorf("%s: %w", a.config.Storage(), ErrStorageTypeNotFound)
	}

	a.db = db

	return nil
}

func (a *App) addFixtures() error {
	fake := faker.New()

	userRepo := repository.NewUserDBRepository(a.db)
	productRepo := repository.NewProductDBRepository(a.db)

	for i := range 10 {
		var u *entity.User
		if i == 0 {
			hashedPass, err := a.hasher.HashPassword("12345678")

			if err != nil {
				return err
			}

			u = &entity.User{
				Name:     "admin@admin.com",
				Password: hashedPass,
			}
		} else {
			hashedPass, err := a.hasher.HashPassword(fake.Internet().Password())

			if err != nil {
				return err
			}

			u = &entity.User{
				Name:     fake.Internet().Email(),
				Password: hashedPass,
			}
		}

		if _, err := userRepo.Create(u); err != nil {
			return err
		}
	}

	for i := 0; i < 20; i++ {
		var p *entity.Product
		if i == 0 {
			p = &entity.Product{
				Name:        "test",
				Description: fake.Address().Address(),
				Price:       333,
			}
		} else {
			p = &entity.Product{
				Name:        fmt.Sprintf("%d:%s", i, fake.Beer().Name()),
				Description: fake.Address().Address(),
				Price:       fake.Float64(1, 100, 10000),
			}
		}

		if _, err := productRepo.Create(p); err != nil {
			return err
		}
	}

	return nil
}
