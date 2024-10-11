package util

import "github.com/spf13/viper"

// Config holds all configurations of the application.
// Viper will read the values from a config file or environment variables.
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"` // `mapstructure` tag is used by Viper to map env variables or config file values to struct fields.
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfiguration reads configuration from a file at the given path or from environment variables.
func LoadConfiguration(path string) (config Config, err error) {
	viper.AddConfigPath(path)  // Set the path to look for the configuration file.
	viper.SetConfigName("app") // Name of the config file (without extension).
	viper.SetConfigType("env") // The type of the config file (in this case, .env format).

	// Viper will check if an environment variable with the same name as a config key exists.
	// If it does, the environment variable value will take precedence over the value in the config file.
	viper.AutomaticEnv()

	// Read in the configuration file, if it exists.
	err = viper.ReadInConfig()
	if err != nil {
		return // Return early if there's an error reading the config file.
	}

	// Unmarshal the config file or environment variables into the Config struct.
	err = viper.Unmarshal(&config)
	return
}
