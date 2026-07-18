package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configs struct{
    Api ApiConfig
	Db DatabaseConfig
	Auth AuthConfig
}

type ApiConfig struct {
   Addr string	
}

type DatabaseConfig struct {
	Url string
}

type AuthConfig struct {
	RefreshTokenSecret string
	AccessTokenSecret string
}

const ENV_PREFIX = "MOVIE_SERVICE_"

func InitConfigs() Configs {
	// Fail silently for production
	_ = godotenv.Load()
	
	return Configs{
		Api: ApiConfig{
			Addr: getEnv(ENV_PREFIX+"ADDR", ":8080"),
		},
	}
}

func getEnv(key string, defaultValue string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func getEnvFromInt(key string, defaultValue int) int {

	if value, ok := os.LookupEnv(key); ok {
		num, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}

		return num
	}

	return defaultValue
}