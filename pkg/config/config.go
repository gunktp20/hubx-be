package config

import (
	"log"
	"sync"

	"github.com/gunktp20/digital-hubx-be/pkg/constant"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server        *Server              `mapstructure:"SERVER"`
		Db            *Db                  `mapstructure:"DB"`
		GcsSignedUrl  *GcsSignedUrlConfig  `mapstructure:"GCS_SIGNED_URL"`
		BusinessLogic *BusinessLogicConfig `mapstructure:"BUSINESS_LOGIC"`
	}

	Server struct {
		Port int `mapstructure:"SERVER_PORT"`
	}

	Db struct {
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
)

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig(configPath string) *Config {
	once.Do(func() {
		v := viper.New()
		v.SetConfigName("config")
		v.SetConfigType("json")
		v.AddConfigPath(configPath)

		if err := v.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		var config Config
		if err := v.Unmarshal(&config); err != nil {
			log.Fatalf("Error unmarshaling config, %s", err)
		}

		if config.Server == nil || config.Db == nil {
			log.Fatalf("Missing essential configuration")
		}

		configInstance = &config
		log.Println(constant.Green + "Configuration loaded successfully" + constant.Reset)
	})

	return configInstance
}
