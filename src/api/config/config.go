package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
	"api/internal/security"
)

// Config represents application configuration
type Config struct {
	DBHost          string
	DBPort          int
	DBUser          string
	DBPassword      string
	DBName          string
	SecretKey       string
	TokenAge        int
	GraphQLPort     int
	LogServerPort   int
	LogMergeMin     int
	LogMoveMin      float64
	RateLimitReqSec int
	RateLimitBurst  int
	CacheProvider   int
	CacheConnString string
	CachePassword   string
	CacheIndex      int
	CacheAge        int
}

const (
	DBHost          = "DB_HOST"
	DBPort          = "DB_PORT"
	DBUser          = "DB_USER"
	DBPassword      = "DB_PASSWORD"
	DBName          = "DB_NAME"
	SecretKey       = "SECRETE_KEY"
	TokenAge        = "TOKEN_AGE"
	GraphQLPort     = "GRAPHQL_PORT"
	LogServerPort   = "LOG_SERVER_PORT"
	LogMergeMin     = "LOG_MERGE_MIN"
	LogMoveMin      = "LOG_MOVE_MIN"
	RateLimitReqSec = "RATE_LIMIT_REQ_SEC"
	RateLimitBurst  = "RATE_LIMIT_BURST"
	CacheAge        = "CACHE_AGE"
	CacheConnString = "CACHE_CONNECTION_STRING"
	CacheIndex      = "CACHE_INDEX"
	CachePassword   = "CACHE_PASSWORD"
)

var instance *Config
var once sync.Once
const KEY = "2b2aac4013cff37435cb22ba6b0e338e2afbef15ac02abcf0a89de5d06d6ae10"

// LoadConfig loads the configuration from environment variables
func NewConfig() *Config {
	once.Do(func() {
		relativePath := "../../config/encrypted.env"

		// Get the absolute path
		absolutePath, err := filepath.Abs(relativePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		content, err := os.ReadFile(absolutePath)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		decryptedEnv, err := security.Decrypt(string(content), KEY)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Load the decrypted environment variables into Viper
		viper.SetConfigType("env")
		err = viper.ReadConfig(bytes.NewBufferString(decryptedEnv))
		if err != nil {
			fmt.Println("Failed to load env variables:", err)
			return
		}

		viper.AutomaticEnv()

		// Create a Config instance and set values from Viper
		instance = &Config{
			DBHost:          viper.GetString(DBHost),
			DBPort:          viper.GetInt(DBPort),
			DBUser:          viper.GetString(DBUser),
			DBPassword:      viper.GetString(DBPassword),
			DBName:          viper.GetString(DBName),
			SecretKey:       viper.GetString(SecretKey),
			TokenAge:        viper.GetInt(TokenAge),
			GraphQLPort:     viper.GetInt(GraphQLPort),
			LogServerPort:   viper.GetInt(LogServerPort),
			LogMergeMin:     viper.GetInt(LogMergeMin),
			LogMoveMin:      viper.GetFloat64(LogMoveMin),
			RateLimitReqSec: viper.GetInt(RateLimitReqSec),
			RateLimitBurst:  viper.GetInt(RateLimitBurst),
		}
	})
	return instance
}

// GetConfig returns the singleton configuration instance
func GetConfig() *Config {
	return instance
}
