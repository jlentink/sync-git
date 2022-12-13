package cnf

import (
	"github.com/spf13/viper"
	"log"
)

var (
	_config        *viper.Viper
	_source        *GitDestination
	_destinations  []*GitDestination
	ConfigLocation string
)

// Get the configuration.
func Get() *viper.Viper {
	if _config == nil {
		_config = viper.New()
		if ConfigLocation == "" {
			_config.SetConfigType("toml")
			_config.AddConfigPath("/etc")
			_config.AddConfigPath(".")
			_config.SetConfigName("sync-git")
		} else {
			_config.SetConfigFile(ConfigLocation)
		}

		if err := _config.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Fatalf("%s\n", err.Error())
			} else {
				log.Fatalf("Error parsing cnf file. (%s)\n", err.Error())
			}
		}
	}
	return _config
}

func GetString(key string) string {
	return Get().GetString(key)
}

func GetInt(key string) int {
	return Get().GetInt(key)
}

func GetSource() *GitDestination {
	if _source != nil {
		return _source
	}
	source := GitDestination{}
	err := Get().UnmarshalKey("source", &source)
	if err != nil {
		log.Fatalf("Error unmarshalling source %s", err)
	}
	_source = &source
	return &source
}
func GetDestinations() []*GitDestination {
	if _destinations != nil {
		return _destinations
	}
	err := Get().UnmarshalKey("destination", &_destinations)
	if err != nil {
		log.Fatalf("Error unmarshalling destinations %s", err)
	}
	return _destinations
}
