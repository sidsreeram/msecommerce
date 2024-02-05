package config
import "github.com/spf13/viper"

type Config struct{
	Port string `mapstructure:"PORT"`
	DBUrl string `mapstructure:"DB_URL"`
	ProductSvcUrl string `mapstructure:"PRODUCT_SVC_URL"`
}

func LoadConfig()(c Config,err error){
	viper.AddConfigPath(".")
	viper.SetConfigName("cart")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err=viper.ReadInConfig()

	if err!=nil{
		return
	}
	err=viper.Unmarshal(&c)

	return
}