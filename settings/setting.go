package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	//"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("config") // 指定配置文件名称
	viper.SetConfigType("yaml")   // 指定配置文件类型
	viper.AddConfigPath(".")      // 指定查找配置文件路径（相对路径）
	err = viper.ReadInConfig()    // 读取配置
	if err != nil {               // 读取配置失败
		fmt.Printf("Fatal error config file: %s \n", err)
		return
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
	})
	return

	//r := gin.Default()
	//if err = r.Run(fmt.Sprintf(":%d", viper.Get("port"))); err != nil {
	//	panic(err)
	//}
}
