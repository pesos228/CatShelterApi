package main

import (
	"api/catshelter/internal/domain"
	"api/catshelter/internal/handler"
	"api/catshelter/internal/repository"
	"api/catshelter/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := initEnv()
	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("Connection to DB failed : %v", err)
	}

	tokenAuth := jwtauth.New("HS256", []byte(cfg.Secret), nil)

	migrateTables(db)

	roleRepository := repository.NewRoleRepositoryImpl(db)
	userRepository := repository.NewUserReposioryImpl(db)
	refreshTokenRepository := repository.NewRefreshTokenRepositoryImpl(db)

	authService := service.NewAuthService(userRepository, roleRepository)
	tokenService := service.NewTokenService(tokenAuth, refreshTokenRepository, userRepository)
	userService := service.NewUserService(userRepository)

	authHandler := handler.NewAuthHandler(authService, tokenService)
	userHandler := handler.NewUserHandler(userService)

	err = initRoles(context.Background(), roleRepository)
	if err != nil {
		log.Fatalf("Bad init roles in DB: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))

		r.Get("/user/info", userHandler.AboutMe)
		r.Post("/api/auth/login", authHandler.Login)
		r.Post("/api/auth/register", authHandler.Register)
		r.With(httprate.LimitByIP(5, 1*time.Minute)).Post("/api/update-session", authHandler.UpdateSession)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/api/auth/logout", authHandler.Logout)
	})

	log.Printf("The server starts on port %s\n", cfg.HTTPport)
	http.ListenAndServe(cfg.HTTPport, r)
}
func initEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file, system variables are used")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("The DATABASE_URL environment variable is not set")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		log.Println("The HTTP_PORT environment variable is not set, the default port 3000 is used")
		httpPort = ":3000"
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("The SECRET environment variable is not set")
	}

	return &Config{
		DatabaseUrl: databaseUrl,
		HTTPport:    httpPort,
		Secret:      secret,
	}
}

func migrateTables(db *gorm.DB) {
	db.AutoMigrate(&domain.Role{})
	db.AutoMigrate(&domain.Cat{})
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&repository.RefreshToken{})
}
func initRoles(ctx context.Context, r repository.RoleRepository) error {
	err := isExistsElseCreateRole("admin", r, ctx)
	if err != nil {
		return err
	}
	err = isExistsElseCreateRole("user", r, ctx)
	if err != nil {
		return err
	}
	return nil
}

func isExistsElseCreateRole(checkRole string, r repository.RoleRepository, ctx context.Context) error {
	role, err := r.FindByName(ctx, checkRole)
	if err != nil {
		if !errors.Is(err, repository.ErrRoleNotFound) {
			return err
		}
		role, _ = domain.NewRole(checkRole)
		err = r.Save(ctx, role)
		if err != nil {
			return err
		}
		log.Printf("role '%s' created\n", checkRole)
		return nil
	}
	if role != nil {
		log.Printf("role '%s' already exists\n", checkRole)
	}
	return nil
}

type Config struct {
	DatabaseUrl string
	HTTPport    string
	Secret      string
}
