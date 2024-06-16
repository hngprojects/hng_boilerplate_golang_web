package config

import (
	"log"

	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"github.com/spf13/viper"
)

// Setup initialize configuration
var (
	// Params ParamsConfiguration
	Config *Configuration
)

// Params = getConfig.Params
func Setup(logger *utility.Logger, name string) *Configuration {
	var baseConfiguration *BaseConfig

	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.AutomaticEnv()

	err := viper.Unmarshal(&baseConfiguration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	configuration := baseConfiguration.SetupConfigurationn()

	// Params = configuration.Params
	Config = configuration
	logger.Info("configurations loading successfully")
	return configuration
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return Config
}
