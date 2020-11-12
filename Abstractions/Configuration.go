package Abstractions

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoyofx/yoyogo/Utils"
	"log"
	"path"
)

type Configuration struct {
	context *ConfigurationContext
	config  *viper.Viper
}

func NewConfiguration(configContext *ConfigurationContext) *Configuration {
	defaultConfig := viper.New()
	if configContext.enableEnv {
		defaultConfig.AutomaticEnv()
		defaultConfig.SetEnvPrefix("YYG")
	}
	if configContext.enableFlag {
		flag.String("app", "", "application name")
		flag.String("port", "", "application port")
		flag.String("profile", configContext.profile, "application profile")
		flag.String("conf", ".", "config dir")
		pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
		pflag.Parse()
		_ = defaultConfig.BindPFlags(pflag.CommandLine)
	}

	if pf := defaultConfig.GetString("profile"); pf != "" {
		configContext.profile = pf
	}

	if cf := defaultConfig.GetString("conf"); cf != "" {
		configContext.configDir = cf
	}

	configName := configContext.configName + "_" + configContext.profile
	configFilePath := path.Join(configContext.configDir, configName+"."+configContext.configType)
	exists, _ := Utils.PathExists(configFilePath)
	if !exists {
		configName = configContext.configName
	}

	defaultConfig.AddConfigPath(configContext.configDir)
	defaultConfig.SetConfigName(configName)
	defaultConfig.SetConfigType(configContext.configType)
	if err := defaultConfig.ReadInConfig(); err != nil {
		panic(err)
		return nil
	}
	log.Println(configFilePath)

	return &Configuration{
		context: configContext,
		config:  defaultConfig,
	}
}

func (c *Configuration) Get(name string) interface{} {
	return c.config.Get(name)
}

func (c *Configuration) GetString(name string) string {
	return c.config.GetString(name)
}

func (c *Configuration) GetInt(name string) int {
	return c.config.GetInt(name)
}

func (c *Configuration) GetSection(name string) IConfiguration {
	section := c.config.Sub(name)

	if section != nil {
		return &Configuration{config: section}
	}
	return nil
}

func (c *Configuration) Unmarshal(obj interface{}) {
	err := c.config.Unmarshal(obj)
	if err != nil {
		fmt.Println("unmarshal config is failed, err:", err)
	}
}

func (c *Configuration) GetProfile() string {
	return c.context.profile
}

func (c *Configuration) GetConfDir() string {
	return c.context.configDir
}
