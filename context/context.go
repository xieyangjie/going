package context

import (
	"os"
	"path"
	"strings"
	"encoding/json"
	"errors"
	"github.com/smartwalle/going/convert"
	"sync"
)

var sharedContext *Context

type Context struct {
	sync.RWMutex
	contextDir		string
	contextList 	[]string
	currentContext	string
	data 			map[string]map[string]interface{}
}

func init() {
	SharedContext()
}

// NewContext 创建新的 Context 实例
func NewContext() (*Context) {
	var context = &Context{}
	return context
}

// SharedContext 获取共享的 Context 实例
func SharedContext() (*Context) {
	if sharedContext == nil {
		sharedContext = &Context{}
	}
	return sharedContext
}

// LoadContext 从目录加载数据
// @contextDir 配置文件所在目录
func (this *Context) LoadContexts(contextDir string) (error) {
	err := this.loadContextsWithDir(contextDir)
	return err
}

func (this *Context) loadContextsWithDir(contextDir string) error {
	this.contextDir = contextDir

	dir, err := os.Open(contextDir)
	if err != nil {
		return err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	this.data = make(map[string]map[string]interface{})
	this.contextList = make([]string, 0, len(names))

	for index, name := range names {
		path := path.Join(contextDir, name)

		file, err := os.Stat(path)
		if err != nil {
			continue
		}

		if !file.IsDir() {
			contextName := strings.Split(file.Name(), ".")[0]

			if data, err := this.loadContextWithPath(path); err == nil {
				if index == 0 {
					this.currentContext = contextName
				}
				this.data[contextName] = data
				this.contextList = append(this.contextList, contextName)
			} else {
				return err
			}
		}
	}
	return nil
}

func (this *Context) loadContextWithPath(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	var data = make(map[string]interface{})
	jsonDecoder := json.NewDecoder(file)
	err = jsonDecoder.Decode(&data)

	return data, err
}

// SetDefaultContext 设定默认的配置文件
// 设定默认的配置文件，如果指定名称的配置文件不存在，则设定为第一个配置。
// @contextName 配置名称
func (this *Context) SetDefaultContext(contextName string) bool {
	_, ok := this.data[contextName]
	if ok {
		this.currentContext = contextName
	} else if len(this.contextList) > 0{
		this.currentContext = this.contextList[0]
		ok = true
	}
	return ok
}

func (this *Context) GetDefaultContext() string {
	return this.currentContext
}

func (this *Context) ContextExist(contextName string) bool {
	_, ok := this.data[contextName]
	return ok
}

func (this *Context) KeyExist(contextName string, key string) bool {
	_, ok := this.data[contextName][key]
	return ok
}

// GetValueWithContext 从指定配置文件读取数据
// 从指定配置文件读取数据，如果指定配置(@contextName)或者key(@key)不存在，则返回默认值(@defaultValue)，并且附上对应的错误信息。
// 如果指定配置信息和key都存在，则什么对应的值，错误信息为空。
// @contextName 配置名称
// @key key
// @defaultValue 默认值
func (this *Context) GetValueWithContext(contextName string, key string, defaultValue interface{}) (interface{}, error) {
	this.RLock()
	defer this.RUnlock()

	//data 为空
	if this.data == nil || len(this.data) == 0 {
		return defaultValue, errors.New("读取配置信息失败")
	}

	//如果 contextName 为空字符串的时候，则返回默认值
	if len(contextName) == 0 {
		return defaultValue, errors.New("配置名称不能为空字符串")
	}

	//如果没有相应的配置
	if !this.ContextExist(contextName) {
		return defaultValue, errors.New("配置[" + contextName + "]不存在")
	}

	//如果没有相应的key
	if !this.KeyExist(contextName, key) {
		return defaultValue, errors.New("字段[" + key + "]不存在")
	}

	return this.data[contextName][key], nil
}

// GetValue 获取指定字段的值
// 如果字段不存在，则返回nil和错误信息；如果字段存在，则返回其值，error为空。
func (this *Context) Get(key string) (interface{}, error) {
	return this.GetValueWithContext(this.currentContext, key, nil)
}

// Get 获取指定字段的值
// 如果字段不存在，则返回提供的默认值；如果字段存在，则返回其值。
func (this *Context) GetValue(key string, defaultValue interface{}) interface{} {
	var value, _ = this.GetValueWithContext(this.currentContext, key, defaultValue)
	return value
}

func (this *Context) GetList(key string, defaultValue []interface{}) []interface{} {
	var value = this.GetValue(key, defaultValue)
	if v, ok := value.([]interface{}); ok {
		return v
	}
	return defaultValue
}

func (this *Context) GetMap(key string, defaultValue map[string]interface{}) map[string]interface{} {
	var value = this.GetValue(key, defaultValue)
	if v, ok := value.(map[string]interface{}); ok {
		return v
	}
	return defaultValue
}

func (this *Context) GetString(key string, defaultValue string) string {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToString(value)
}

func (this *Context) GetInt(key string, defaultValue int) int {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToInt(value)
}

func (this *Context) GetInt32(key string, defaultValue int32) int32 {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToInt32(value)
}

func (this *Context) GetInt64(key string, defaultValue int64) int64 {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToInt64(value)
}

func (this *Context) GetFloat(key string, defaultValue float32) float32 {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToFloat32(value)
}

func (this *Context) GetFloat64(key string, defaultValue float64) float64 {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToFloat64(value)
}

func (this *Context) GetBool(key string, defaultValue bool) bool {
	var value = this.GetValue(key, defaultValue)
	return convert.ConvertToBool(value)
}


func SetDefaultContext(contextName string) bool {
	return sharedContext.SetDefaultContext(contextName)
}

func GetValue(key string) (interface{}, error) {
	return sharedContext.Get(key)
}

func Get(key string, defaultValue interface{}) interface{} {
	return sharedContext.GetValue(key, defaultValue)
}

func GetList(key string, defaultValue []interface{}) []interface{} {
	return sharedContext.GetList(key, defaultValue)
}

func GetMap(key string, defaultValue map[string]interface{}) map[string]interface{} {
	return sharedContext.GetMap(key, defaultValue)
}

func GetString(key string, defaultValue string) string {
	return sharedContext.GetString(key, defaultValue)
}

func GetInt(key string, defaultValue int) int {
	return sharedContext.GetInt(key, defaultValue)
}

func GetInt32(key string, defaultValue int32) int32 {
	return sharedContext.GetInt32(key, defaultValue)
}

func GetInt64(key string, defaultValue int64) int64 {
	return sharedContext.GetInt64(key, defaultValue)
}

func GetFloat(key string, defaultValue float32) float32 {
	return sharedContext.GetFloat(key, defaultValue)
}

func GetFloat64(key string, defaultValue float64) float64 {
	return sharedContext.GetFloat64(key, defaultValue)
}

func GetBool(key string, defaultValue bool) bool {
	return sharedContext.GetBool(key, defaultValue)
}