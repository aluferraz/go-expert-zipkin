package configs

import (
	"github.com/spf13/viper"
	"os"
	"path"
)

type conf struct {
	WebserverPort string `mapstructure:"WEBSERVER_PORT"`
	WeatherApiKey string `mapstructure:"WEATHER_API_KEY"`
	WeatherApiURL string `mapstructure:"WEATHER_API_URL"`
	CepApiURL     string `mapstructure:"CEP_API_URL"`
}

func defaultAndBindings() error {
	defaultConfigs := map[string]string{
		"WEBSERVER_PORT":  "8080",
		"WEATHER_API_KEY": "",
		"WEATHER_API_URL": "https://api.weatherapi.com/v1/current.json?key=%s&q=%s",
		"CEP_API_URL":     "https://viacep.com.br/ws/%s/json/",
	}
	for envKey, envValue := range defaultConfigs {
		err := viper.BindEnv(envKey)
		if err != nil {
			return err
		}
		viper.SetDefault(envKey, envValue)
	}
	return nil

}
func LoadConfig(workdir string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	_, err := os.Stat(path.Join(workdir, ".env"))
	if err == nil {
		viper.SetConfigType("env")
		viper.AddConfigPath(workdir)
		viper.SetConfigFile(".env")
		err = viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}
	viper.AutomaticEnv()
	err = defaultAndBindings()
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
