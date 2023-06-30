package raid

import (
	"bytes"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Settings struct {
	TelegramChannel string         `mapstructure:"telegram_channel"`
	TimezoneName    string         `mapstructure:"timezone_name"`
	Timezone        *time.Location ``
	APIKeys         []string       `mapstructure:"api_keys"`
	Debug           bool           `mapstructure:"debug"`
	Trace           bool           `mapstructure:"trace"`
	BacklogSize     int            `mapstructure:"backlog_size"`
	Host            string         `mapstructure:"host"`
	Port            uint16         `mapstructure:"port"`
}

func init() {
	viper.SetDefault("telegram_channel", "air_alert_ua")
	viper.SetDefault("timezone_name", "Europe/Kiev")
	viper.SetDefault("api_keys", []string{})
	viper.SetDefault("debug", false)
	viper.SetDefault("trace", false)
	viper.SetDefault("backlog_size", 500)
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", "10101")

	viper.SetConfigFile("settings.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/raid/")
	viper.AddConfigPath("$HOME/.raid")
	viper.AddConfigPath(".")

	if len(os.Args) > 1 {
		fname := os.Args[1]
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			log.Fatalf("settings: open settings file: %s", err)
		}

		r := bytes.NewReader(data)
		if err := viper.ReadConfig(r); err != nil {
			log.Fatalln(err)
		}
	} else {
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("settings: read config: %s", err)
		}
	}
}

func MustLoadSettings() (settings Settings) {
	if err := viper.Unmarshal(&settings); err != nil {
		log.Fatalf("settings: unmarshal: %s", err)
	}

	location, err := time.LoadLocation(settings.TimezoneName)
	if err != nil {
		log.Fatalf("settings: load timezone: %s", err)
	}

	settings.Timezone = location

	if len(settings.APIKeys) == 0 {
		log.Fatal("settings: no API keys were loaded")
	}

	log.Infof("settings: load %d API keys", len(settings.APIKeys))

	return
}
