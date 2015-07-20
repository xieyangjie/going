package config

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"path"
	"github.com/smartwalle/going/convert"
)

////////////////////////////////////////////////////////////////////////////////
type Config struct {
	sync.RWMutex
	configPath string
	data       map[string]interface{}
}

func NewConfig() *Config {
	var config = &Config{}
	return config
}

func (this *Config) SetConfigFile(filePath string) error {
	this.configPath = filePath

	var err error
	this.data, err = this.loadConfigWithPath(filePath)

	if err != nil || this.data == nil {
		this.data = make(map[string]interface{})
	}

	return err
}

func (this *Config) loadConfigWithPath(filePath string) (map[string]interface{}, error) {
	var file, err = os.Open(filePath)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	var data = make(map[string]interface{})
	var jsonDecoder = json.NewDecoder(file)
	err = jsonDecoder.Decode(&data)
	return data, err
}

func (this *Config) SaveConfig() error {
	if this.data == nil {
		this.data = make(map[string]interface{})
	}

	if this.configPath == "" {
		return errors.New("保存路径不能为空")
	}

	//创建目录
	var dir, _ = path.Split(this.configPath)
	if dir != "" {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				os.MkdirAll(dir, os.ModeDir|os.ModePerm)
			}
		}
	}

	file, err := os.Create(this.configPath)
	defer file.Close()

	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(this.data, "", "	")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

func (this *Config) Reset() {
	this.data = make(map[string]interface{})
	this.SaveConfig()
}

func (this *Config) KeyExist(key string) bool {
	_, ok := this.data[key]
	return ok
}

func (this *Config) SetValue(key string, value interface{}) {
	this.Lock()
	defer this.Unlock()
	this.data[key] = value
}

func (this *Config) RemoveKey(key string) {
	this.Lock()
	defer this.Unlock()
	if this.KeyExist(key) {
		delete(this.data, key)
	}
}

func (this *Config) GetValue(key string, defaultValue interface{}) interface{} {
	if len(key) == 0 {
		return nil
	}

	this.RLock()
	defer this.RUnlock()

	if !this.KeyExist(key) {
		return defaultValue
	}
	return this.data[key]
}

func (this *Config) GetList(key string, defaultValue []interface{}) []interface{} {
	var value = this.GetValue(key, defaultValue)
	if v, ok := value.([]interface{}); ok {
		return v
	}
	return defaultValue
}

func (this *Config) GetMap(key string, defaultValue map[string]interface{}) map[string]interface{} {
	var value = this.GetValue(key, defaultValue)
	if v, ok := value.(map[string]interface{}); ok {
		return v
	}
	return defaultValue
}

func (this *Config) GetString(key string, defaultValue string) string {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToString(value)
}

func (this *Config) GetInt(key string, defaultValue int) int {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToInt(value)
}

func (this *Config) GetInt64(key string, defaultValue int64) int64 {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToInt64(value)
}

func (this *Config) GetFloat(key string, defaultValue float32) float32 {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToFloat32(value)
}

func (this *Config) GetFloat64(key string, defaultValue float64) float64 {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToFloat64(value)
}

func (this *Config) GetBool(key string, defaultValue bool) bool {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToBool(value)
}

////////////////////////////////////////////////////////////////////////////////
var sharedConfig *Config

func SharedConfig() *Config {
	if sharedConfig == nil {
		sharedConfig = &Config{}
	}
	return sharedConfig
}
