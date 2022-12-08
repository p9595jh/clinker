package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	V = viper.GetViper()
)

func init() {
	viper.SetConfigFile("./.env")
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("\033[33mno .env file - guess it's a test mode (will be closed after 31s)\033[0m")
			go func() { time.Sleep(time.Second * 31); panic(err) }()
		} else {
			panic(err)
		}
	}

	// "THIS_IS_SAMPLE" -> "this.is.sample"
	for _, key := range viper.AllKeys() {
		key2 := strings.ReplaceAll(key, "_", ".")
		viper.Set(key2, viper.Get(key))
	}
}

func ToMap() map[string]any {
	keys := viper.AllKeys()
	configMap := make(map[string]any, len(keys))
	for _, key := range keys {
		if !strings.Contains(key, "_") {
			configMap[key] = viper.Get(key)
		}
	}
	return configMap
}

func ToNestedMap() map[any]any {
	keys := viper.AllKeys()
	configMap := make(map[any]any, len(keys))
	for _, key := range keys {
		if strings.Contains(key, "_") {
			continue
		}
		paths := strings.Split(key, ".")
		lastPathIdx := len(paths) - 1
		pos := &configMap
		for _, path := range paths[:lastPathIdx] {
			if inner, ok := (*pos)[path]; ok {
				m := inner.(map[any]any)
				pos = &m
			} else {
				m := map[any]any{}
				(*pos)[path] = m
				pos = &m
			}
		}
		(*pos)[paths[lastPathIdx]] = viper.Get(key)
	}
	return configMap
}
