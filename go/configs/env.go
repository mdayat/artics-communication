package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseURL    string
	AllowedOrigins string
	SecretKey      string
	OriginURL      string
	CookieDomain   string
}

func LoadEnv(filenames ...string) (Env, error) {
	if err := godotenv.Load(filenames...); err != nil {
		return Env{}, err
	}

	env := Env{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
		SecretKey:      os.Getenv("SECRET_KEY"),
		OriginURL:      os.Getenv("ORIGIN_URL"),
		CookieDomain:   os.Getenv("COOKIE_DOMAIN"),
	}

	return env, nil
}
