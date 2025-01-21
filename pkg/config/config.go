package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/gunktp20/digital-hubx-be/pkg/constant"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server        *ServerConfig        `mapstructure:"SERVER"`
		Db            *DbConfig            `mapstructure:"DB"`
		GcsSignedUrl  *GcsSignedUrlConfig  `mapstructure:"GCS_SIGNED_URL"`
		BusinessLogic *BusinessLogicConfig `mapstructure:"BUSINESS_LOGIC"`
		Swagger       *SwaggerConfig       `mapstructure:"SWAGGER"`
		Logger        *LoggerConfig        `mapstructure:"LOGGER"`
		CORS          *CORSConfig          `mapstructure:"CORS"`
	}

	ServerConfig struct {
		Port           int `mapstructure:"SERVER_PORT"`
		ReadBufferSize int `mapstructure:"READ_BUFFER_SIZE"`
	}

	DbConfig struct {
		SSLMode         string `mapstructure:"DB_SSLMODE"`
		Dns             string `mapstructure:"DB_DNS"`
		SSLRootCertPath string `mapstructure:"DB_SSLROOTCERT_PATH"`
		SSLKeyPath      string `mapstructure:"DB_SSLKEY_PATH"`
		SSLCertPath     string `mapstructure:"DB_SSLCERT_PATH"`
		DbName          string `mapstructure:"DB_NAME"`
		DbSecret        string `mapstructure:"DB_PASSWORD"`
		DbUser          string `mapstructure:"DB_USER"`
		DbPort          string `mapstructure:"DB_PORT"`
		DbHost          string `mapstructure:"DB_HOST"`
	}

	GcsSignedUrlConfig struct {
		ApiKey      string `mapstructure:"GCS_SIGNED_URL_API_KEY"`
		ServiceName string `mapstructure:"GCS_SIGNED_URL_SERVICE_NAME"`
		UploadUrl   string `mapstructure:"GCS_SIGNED_URL_UPLOAD"`
		DownloadUrl string `mapstructure:"GCS_SIGNED_URL_DOWNLOAD"`
		Expired     int    `mapstructure:"GCS_SIGNED_URL_EXPIRED"`
		BucketName  string `mapstructure:"GCS_SIGNED_URL_BUCKET_NAME"`
		Path        string `mapstructure:"GCS_SIGNED_URL_PATH"`
	}

	BusinessLogicConfig struct {
		MaxCancelPerClass                   int `mapstructure:"MAX_CANCEL_PER_CLASS"`
		DaysBeforeClassStartForCancellation int `mapstructure:"DAYS_BEFORE_CLASS_START_FOR_CANCELLATION"`
		MaxCapacityPerSession               int `mapstructure:"MAX_CAPACITY_PER_SESSION"`
	}

	SwaggerConfig struct {
		Enabled bool `mapstructure:"ENABLED"`
	}

	LoggerConfig struct {
		Enabled bool `mapstructure:"ENABLED"`
	}

	CORSConfig struct {
		AllowOrigins     string `mapstructure:"ALLOW_ORIGINS"`
		AllowMethods     string `mapstructure:"ALLOW_METHODS"`
		AllowHeaders     string `mapstructure:"ALLOW_HEADERS"`
		AllowCredentials bool   `mapstructure:"ALLOW_CREDENTIALS"`
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig(configPath string) (*Config, error) {
	once.Do(func() {
		v := viper.New()
		v.SetConfigName("config")
		v.SetConfigType("json")
		v.AddConfigPath(configPath)

		// Set default values
		v.SetDefault("SERVER.SERVER_PORT", 3000)
		v.SetDefault("SERVER.READ_BUFFER_SIZE", 60000)
		v.SetDefault("SWAGGER.ENABLED", true)

		if err := v.ReadInConfig(); err != nil {
			log.Printf("Error reading config file: %s", err)
			return
		}

		var config Config
		if err := v.Unmarshal(&config); err != nil {
			log.Printf("Error unmarshaling config: %s", err)
			return
		}

		// Validate Config
		if err := validateConfig(&config); err != nil {
			log.Printf("Invalid configuration: %s", err)
			return
		}

		configInstance = &config
		log.Println(constant.Green + "Configuration loaded successfully" + constant.Reset)
	})

	if configInstance == nil {
		return nil, fmt.Errorf("failed to load configuration")
	}
	return configInstance, nil
}

func validateConfig(config *Config) error {
	if config.Server == nil || config.Db == nil {
		return fmt.Errorf("missing essential configuration")
	}

	if config.Server.Port <= 0 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.BusinessLogic.MaxCapacityPerSession <= 0 {
		return fmt.Errorf("invalid max capacity per session: %d", config.BusinessLogic.MaxCapacityPerSession)
	}

	if config.Db.DbName == "" || config.Db.DbUser == "" || config.Db.DbHost == "" {
		return fmt.Errorf("database configuration is incomplete")
	}

	return nil
}
