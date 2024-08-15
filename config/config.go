package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var Conf = new(Config)

func Init(filePath string) {
	//1.指定配置文件路径
	//(1)如果指定了配置文件，优先使用显式指定的值
	//(2)如果没有显式指定配置文件，通过系统环境变量获取当前所处的运行模式(通过设置测试和线上运行环境变量能使得线上线下配置隔离)
	//不用修改任何代码而且线上和线下的配置文件隔离开
	//(3)如果既没有显式指定配置文件路径，也没有设置系统环境变量，则从默认的配置文件路径找
	if filePath != "" {
		viper.SetConfigFile(filePath)
	} else if os.Getenv("RUN_MODE") != "" {
		viper.SetConfigFile("./config/config.yaml")
	} else {
		viper.SetConfigFile("./config/config.yaml")
	}
	//2.读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		//信息读取失败
		log.Fatalf("viper.ReadInConfig() failed: %v\n", err)
	}
	//将读取到的信息反序列化到全部变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		log.Fatalf("viper.Unmarshal() failed: %v\n", err)
	}
	//4.监听配置文件，当配置文件发生修改后立即更新配置文件信息到全局变量Conf
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件%s修改了...\n", e.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			log.Fatalf("viper.Unmarshal() failed: %v\n", err)
		}
	})

}

type Config struct {
	App      `mapstructure:"app"`
	Server   `mapstructure:"server"`
	JWT      `mapstructure:"jwt"`
	MySQL    `mapstructure:"mysql"`
	Redis    `mapstructure:"redis"`
	Kafka    `mapstructure:"kafka"`
	RabbitMQ `mapstructure:"rabbitmq"`
	Log      `mapstructure:"log"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Host    string `mapstructure:"host"`
	Version string `mapstructure:"version"`
}

const (
	ModeDebug = "debug"
	ModeProd  = "prod"
)

type Server struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

type JWT struct {
	Secret       string        `mapstructure:"secret"`
	TokenExpired time.Duration `mapstructure:"tokenExpired"`
}

type MySQL struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Database      string `mapstructure:"database"`
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
	TablePrefix   string `mapstructure:"tablePrefix"`
	SingularTable bool   `mapstructure:"singularTable"`
	Charset       string `mapstructure:"charset"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

type Log struct {
	Filename   string `mapstructure:"filename"`
	Ext        string `mapstructure:"ext"`
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge"`
	MaxBackups int    `mapstructure:"maxBackups"`
	SaveDir    string `mapstructure:"saveDir"`
}

type Kafka struct {
	Host  string `mapstructure:"host"`
	Port  int    `mapstructure:"port"`
	Group string `mapstructure:"group"`
}

type RabbitMQ struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
