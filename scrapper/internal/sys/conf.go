package sys

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DSN  string
	ADDR string
}

func LoadConfig(filename string) (*Config, error) {

	if err := godotenv.Load(filename); err != nil {
		return nil, err
	}

	conf := &Config{
		DSN:  getEnv("POSTGRES_DSN"),
		ADDR: getEnv("ADDR"),
	}

	return conf, nil
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("[%s] not found in env", key))
	}

	return val
}
