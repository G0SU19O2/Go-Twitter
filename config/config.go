package config

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type database struct {
	URL string
}

type jwt struct {
	Secret string
	Issuer string

}

type Config struct {
	Database database
	JWT      jwt
}

func LoadEnv(filename string) {
	re := regexp.MustCompile(`^(.*` + "Go-Twitter" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	err := godotenv.Load(string(rootPath) + "/" + filename)
	if err != nil {
		godotenv.Load()
	}
}

func New() *Config {
	godotenv.Load()
	return &Config{
		Database: database{
			URL: os.Getenv("DATABASE_URL"),
		},
		JWT: jwt{
			Secret: os.Getenv("JWT_SECRET"),
			Issuer: os.Getenv("DOMAIN"),
		},
	}
}