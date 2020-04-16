package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

const (
	baseConfig = "config"
	debug      = "debug"
	release    = "release"
)

var (
	App = &Application{}
	cfg *config
)

// Toml config
type config struct {
	data   map[string]toml.Primitive
	source string
}

// Application config
type Application struct {
	Name string `toml:"name"`
	Host string `toml:"host"`
	Port int    `toml:"port"`
	Mode string `toml:"mode"`
}

func Init() error {
	// New config
	cfg = newConfig()

	// Path
	path, err := os.Getwd()
	if err != nil {
		return errors.New("get wd error=" + err.Error())
	}

	// Config file
	file := fmt.Sprintf("%s/config/%s.toml", path, baseConfig)
	str, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.New("read config file error=" + err.Error())
	}

	data := make(map[string]toml.Primitive)

	// Toml decode
	_, err = toml.Decode(string(str), &data)
	if err != nil {
		return errors.New("toml decode error=" + err.Error())
	}

	cfg.data = data

	// Decode application
	err = DecodeKey("app", App)
	if err != nil {
		return errors.New("decode application error=" + err.Error())
	}

	return nil
}

// Decode one key from toml config
func DecodeKey(key string, v interface{}) error {
	data, exist := cfg.data[key]
	if !exist {
		return errors.New("decode key is not exist! key=" + key)
	}

	return toml.PrimitiveDecode(data, v)
}

func (a *Application) IsRelease() bool {
	return a.Mode == release
}

func (a *Application) IsDebug() bool {
	return a.Mode == debug
}

// New toml config
func newConfig() *config {
	return &config{
		data: make(map[string]toml.Primitive),
	}
}
