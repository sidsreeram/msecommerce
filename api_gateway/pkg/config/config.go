package config

import "github.com/spf13/viper"


type Config struct{
	Port string `mapstructure:"PORT"`
	UserSvcUrl string `mapstructure:"USER_SVC_URL"`
	ProductSvcUrl string `mapstructure:"PRODUCT_SVC_URL"`
	CartSvcUrl string `mapstructure:"CART_SVC_URL"`
	Secret string `mapstructure:"SECRET"`
}

func LoadConfig () (c Config,err error){
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err=viper.ReadInConfig()
	if err!=nil{
		return
	}

	err=viper.Unmarshal(&c)
	return
}