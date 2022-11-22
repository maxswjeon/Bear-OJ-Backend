package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/jackpal/gateway"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/maxswjeon/contest-backend/routes"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/server"
)

func set_trusted_proxies(engine *gin.Engine) {
	trusted_proxies := os.Getenv("TRUSTED_PROXIES")

	if trusted_proxies == "" {
		log.Print("No trusted proxy given, disabling X-Forwarded-* headers")
		engine.SetTrustedProxies(nil)
		return
	}

	if strings.TrimSpace(strings.ToLower(trusted_proxies)) == "gateway" {
		log.Print("Trusting only gateway")

		ip, err := gateway.DiscoverGateway()
		if err != nil {
			log.Panic("Failed to discover gateway")
		}
		log.Printf("Found %s as gateway", ip.String())

		engine.SetTrustedProxies([]string{ip.String()})
		return
	}

	trusted_proxies_list := strings.Split(strings.TrimSpace(trusted_proxies), ",")
	log.Printf("Trusted proxies: %v", trusted_proxies_list)
	engine.SetTrustedProxies(trusted_proxies_list)
}

func connect_to_database() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASS"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_SSLMODE"), os.Getenv("TZ"))
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func set_cors_headers(engine *gin.Engine) {
	cors_origin := os.Getenv("CORS_ORIGIN")
	if cors_origin == "" {
		log.Print("No CORS_ORIGIN given, allowing requests from any origin")
		engine.Use(cors.Default())
		return
	}
	cors_origin_list := strings.Split(strings.TrimSpace(cors_origin), ",")
	log.Printf("Allowing Requests from these origins: %v", cors_origin_list)

	cors_methods := os.Getenv("CORS_METHODS")
	if cors_methods == "" {
		cors_methods = "GET,POST,PUT,DELETE,PATCH"
	}
	cors_methods_list := strings.Split(strings.TrimSpace(cors_methods), ",")
	log.Printf("Allowing Requests with these methods: %v", cors_methods_list)

	cors_headers := os.Getenv("CORS_HEADERS")
	if cors_headers == "" {
		cors_headers = "Origin,Accept,Content-Type,Authorization"
	}
	cors_headers_list := strings.Split(strings.TrimSpace(cors_headers), ",")
	log.Printf("Allowing Requests with these headers: %v", cors_headers_list)

	cors_creditentials := os.Getenv("CORS_CREDENTIALS")
	if cors_creditentials == "" {
		cors_creditentials = "true"
	}
	cors_creditentials_bool := cors_creditentials == "true"
	if cors_creditentials_bool {
		log.Print("Allowing Credentials")
	} else {
		log.Print("Not allowing Credentials")
	}

	cors_cache_time := os.Getenv("CORS_CACHE_TIME")
	if cors_cache_time == "" {
		cors_cache_time = "12h"
	}
	cors_cache_time_duration, err := time.ParseDuration(cors_cache_time)
	if err != nil {
		log.Panicf("Failed to parse CORS_CACHE_TIME (given \"%s\") with error: %v", cors_cache_time, err)
	}
	log.Printf("Caching CORS results for %d seconds", cors_cache_time_duration/time.Second)

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     cors_origin_list,
		AllowMethods:     cors_methods_list,
		AllowHeaders:     cors_headers_list,
		AllowCredentials: cors_creditentials_bool,
		MaxAge:           cors_cache_time_duration,
	}))
}

func set_session(engine *gin.Engine) {
	redis_connections_raw := os.Getenv("SESSION_REDIS_CONNECTIONS")
	if redis_connections_raw == "" {
		redis_connections_raw = "16"
	}
	redis_connections, err := strconv.Atoi(redis_connections_raw)
	if err != nil {
		log.Panicf("Failed to parse SESSION_REDIS_CONNECTIONS (given \"%s\") with error: %v", redis_connections_raw, err)
	}

	authentication_key := os.Getenv("SESSION_REDIS_AUTHENTICATION_KEY")
	if authentication_key == "" {
		log.Panicf("SESSION_REDIS_AUTHENTICATION_KEY not set")
	}

	encryption_key := os.Getenv("SESSION_REDIS_ENCRYPTION_KEY")
	if encryption_key == "" {
		log.Panicf("SESSION_REDIS_ENCRYPTION_KEY not set")
	}

	if len(encryption_key) != 32 {
		log.Panicf("Length of SESSION_REDIS_ENCRYPTION_KEY is not 32 characters")
	}

	store, err := redis.NewStoreWithDB(
		redis_connections,
		"tcp",
		os.Getenv("SESSION_REDIS_URL"),
		os.Getenv("SESSION_REDIS_PASSWORD"),
		"1",
		[]byte(authentication_key), []byte(encryption_key))

	store.Options(sessions.Options{
		Path:     "/",
		Domain:   fmt.Sprintf(".%s", os.Getenv("DOMAIN")),
		MaxAge:   6 * 60 * 60, // 6 Hours
		Secure:   true,
		HttpOnly: true,
	})

	if err != nil {
		log.Panicf("Failed to create session store with error: %v", err)
	}

	engine.Use(sessions.Sessions(os.Getenv("SESSION_NAME"), store))
}

func main() {
	Load_env()
	Ensure_env()

	db, err := connect_to_database()
	if err != nil {
		log.Panicf("Failed to connect to database with error: %v", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	err = db.AutoMigrate(&schemas.Contest{}, &schemas.Submit{}, &schemas.User{}, &schemas.Problem{}, &schemas.Admin{}, &schemas.InternalProblem{})
	if err != nil {
		log.Panicf("Failed to migrate database with error: %v", err)
	}

	engine := gin.Default()
	set_trusted_proxies(engine)
	set_cors_headers(engine)
	set_session(engine)

	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db", db)
	})

	var count int64
	db.Model(&schemas.Admin{}).Find(&schemas.Admin{}).Count(&count)
	if count == 0 {
		log.Print("No admins found, creating admin with credentials from environment")
		admin := schemas.Admin{
			Username: os.Getenv("ADMIN_USERNAME"),
			Password: os.Getenv("ADMIN_PASSWORD"),
		}
		db.Create(&admin)
	}

	var internalProblems = []schemas.InternalProblem{}
	db.Find(&internalProblems)
	for _, internalProblem := range internalProblems {
		internalProblem.Valid = false
		db.Save(&internalProblem)
	}

	routes.Apply(engine)

	s, err := server.NewServer()
	if err != nil {
		log.Panicf("Failed to create server with error: %v", err)
	}

	defer s.Close()
	go s.Listen(db)

	engine.Use(func(ctx *gin.Context) {
		ctx.Set("server", s)
	})

	engine.Run()
}
