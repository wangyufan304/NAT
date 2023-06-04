package utils

// TODO 定义一个全局信息的结构体

// 此文件是用来解析配置文件信息的
import (
	"fmt"
	"github.com/spf13/viper"
)

type ParseConfigFromYML struct {
	ViperInstance *viper.Viper
}

func ParseFile(configFileName string) *ParseConfigFromYML {
	pcfy := new(ParseConfigFromYML)
	pcfy.ViperInstance = viper.New()
	pcfy.ViperInstance.SetConfigName(configFileName)
	pcfy.ViperInstance.AddConfigPath(".")
	pcfy.ViperInstance.AddConfigPath("./config")
	pcfy.ViperInstance.AddConfigPath("../config")
	pcfy.ViperInstance.AddConfigPath("../../config")
	pcfy.ViperInstance.SetConfigType("yml")
	err := pcfy.ViperInstance.ReadInConfig()
	if err != nil {
		fmt.Println("[ReadInConfig]")
		return nil
	}
	return pcfy
}
