package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

const (
	baseConfig = "config"
)

var (
	App = &Application{}
	cfg *config
)

type config struct {
	data   map[string]toml.Primitive
	source string
}

type Application struct {
	Name string `toml:"name"`
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func Init() {
	cfg = newConfig()
	path, _ := os.Getwd()
	file := fmt.Sprintf("%s/config/%s.toml", path, baseConfig)

	data := make(map[string]toml.Primitive)

	str, _ := ioutil.ReadFile(file)
	_, _ = toml.Decode(string(str), &data)
	cfg.data = data

	DecodeKey("app", App)

	fmt.Println(App)

}

func DecodeKey(key string, v interface{}) error {
	data, exist := cfg.data[key]
	if (!exist) {
		fmt.Println("error")
	}

	return toml.PrimitiveDecode(data, v)
}

func newConfig() *config {
	return &config{
		data: make(map[string]toml.Primitive),
	}
}
