package initialize

import (
	"be-better/core/global"
	"be-better/utils"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

var configEnv = "local"
var configFile = "application"
var envFile = "env.dat"
var env = ""

func Viper(path... string) *viper.Viper {
	var viperConfig string
	var viperEnv string
	if len(path) == 0{
		flag.StringVar(&viperConfig,"c","","choose config file.")
		flag.StringVar(&viperEnv,"env","","choose config file.")
		flag.Parse()
		if viperConfig == ""{
			if viperEnv == ""{
				if viperEnv = os.Getenv("DEPLOY_ENV"); viperEnv == ""{
					exists, _ := utils.PathExists(envFile)
					if exists {
						envByte, err := ioutil.ReadFile(envFile)
						if err != nil{
							fmt.Printf("读取env文件错误： %v\n", err.Error())
						}
						env = string(envByte)
					}
				}
			}
			viperConfig = "application"
			if env != "" {
				viperConfig = viperConfig + "-" + env
			}
			viperConfig += ".yaml"
			viperConfig = "./assets/config/" + viperConfig
			exists, _ := utils.PathExists(viperConfig)
			if !exists {
				viperConfig = "./assets/config/application.yaml"
			}
			fmt.Printf("读取配置文件路径为: %v\n", viperConfig)
		} else {
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为: %v\n", viperConfig)
		}
	}else {
		viperConfig = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为: %v\n", viperConfig)
	}

	v := viper.New()
	v.SetConfigFile(viperConfig)
	err := v.ReadInConfig()
	if err != nil{
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:",e.Name)
		if err := v.Unmarshal(&global.GlobalConfig); err != nil{
			panic(fmt.Sprint(err))
		}
	})

	if err := v.Unmarshal(&global.GlobalConfig); err != nil{
		panic(fmt.Sprint(err))
	}

	return v
}
