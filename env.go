package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	name     string
	required bool
}

var ENVIRONMENT_VARIABLES = []Env{
	{
		name:     "GIN_MODE",
		required: false,
	},
	{
		name:     "TRUSTED_PROXIES",
		required: false,
	},
	{
		name:     "DOMAIN",
		required: true,
	},
	{
		name:     "CORS_ORIGIN",
		required: false,
	},
	{
		name:     "CORS_METHODS",
		required: false,
	},
	{
		name:     "CORS_HEADERS",
		required: false,
	},
	{
		name:     "CORS_CREDENTIALS",
		required: false,
	},
	{
		name:     "CORS_CACHE_TIME",
		required: false,
	},

	{
		name:     "ADMIN_USERNAME",
		required: true,
	},
	{
		name:     "ADMIN_PASSWORD",
		required: true,
	},

	{
		name:     "SESSION_NAME",
		required: true,
	},
	{
		name:     "SESSION_REDIS_URL",
		required: true,
	},
	{
		name:     "SESSION_REDIS_PASSWORD",
		required: false,
	},
	{
		name:     "SESSION_REDIS_CONECTIONS",
		required: false,
	},
	{
		name:     "SESSION_REDIS_AUTHENTICATION_KEY",
		required: true,
	},
	{
		name:     "SESSION_REDIS_ENCRYPTION_KEY",
		required: true,
	},

	{
		name:     "DATABASE_HOST",
		required: true,
	},
	{
		name:     "DATABASE_PORT",
		required: true,
	},
	{
		name:     "DATABASE_NAME",
		required: true,
	},
	{
		name:     "DATABASE_USER",
		required: true,
	},
	{
		name:     "DATABASE_PASS",
		required: true,
	},
	{
		name:     "DATABASE_SSLMODE",
		required: true,
	},

	{
		name:     "SMTP_DOMAIN",
		required: true,
	},
	{
		name:     "SMTP_PORT",
		required: true,
	},
	{
		name:     "SMTP_EHLO",
		required: true,
	},
	{
		name:     "SMTP_STARTTLS",
		required: true,
	},
	{
		name:     "SMTP_USER",
		required: true,
	},
	{
		name:     "SMTP_PASS",
		required: true,
	},
	{
		name:     "SMTP_SENDER_NAME",
		required: true,
	},
	{
		name:     "SMTP_SENDER_MAIL",
		required: true,
	},

	{
		name:     "SERVER_PORT",
		required: true,
	},
}

func Load_env() {
	err := godotenv.Load()

	if err != nil && !os.IsNotExist(err) {
		log.Panicf("Error loading .env file with error %v", err)
	}
}

func Ensure_env() {
	for _, env := range ENVIRONMENT_VARIABLES {
		if _, ok := os.LookupEnv(env.name); !ok && env.required {
			log.Panicf("%s environment variable not set, but it is required", env.name)
		}
	}
}
