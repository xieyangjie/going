package redis

import (
	"fmt"
	"os"
	"time"
	redigo "github.com/garyburd/redigo/redis"
	"encoding/json"
)

func NewRedis(url, password string, dbIndex, maxActive, maxIdle int) (p *Pool) {
	var dialFunc = func() (c redigo.Conn, err error) {
		if len(password) > 0 {
			c, err = redigo.Dial("tcp", url, redigo.DialPassword(password))
		} else {
			c, err = redigo.Dial("tcp", url)
		}

		if err != nil {
			fmt.Println("连接 Redis 服务器失败:", url, err)
			os.Exit(-1)
		}

		_, err = c.Do("SELECT", dbIndex)
		if err != nil {
			fmt.Println("Redis 执行 SELECT 指令失败:", dbIndex, err)
			c.Close()
			os.Exit(-1)
		}

		return c, err
	}

	p = &Pool{}
	var pool = redigo.NewPool(dialFunc, maxIdle)
	pool.MaxActive = maxActive
	pool.IdleTimeout = 180 * time.Second
	pool.Wait = true
	p.p = pool

	return p
}

////////////////////////////////////////////////////////////////////////////////
type Pool struct {
	p *redigo.Pool
}

func (this *Pool) GetSession() *Session {
	var c = this.p.Get()
	return &Session{c}
}

func (this *Pool) Release(s *Session) {
	s.c.Close()
}

////////////////////////////////////////////////////////////////////////////////
type Session struct {
	c Conn
}

func (this *Session) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return this.c.Do(commandName, args...)
}

////////////////////////////////////////////////////////////////////////////////
func (this *Session) EXISTS(key string) (bool) {
	return MustBool(this.Do("EXISTS", key))
}

func (this *Session) EXPIRE(key string, seconds int) (reply interface{}, err error) {
	return this.Do("EXPIRE", key, seconds)
}

func (this *Session) DEL(key ...interface{}) (reply interface{}, err error) {
	return this.Do("DEL", key...)
}

////////////////////////////////////////////////////////////////////////////////
func (this *Session) GET(key string) (reply interface{}, err error) {
	return this.Do("GET", key)
}

func (this *Session) SET(key string, value interface{}) (reply interface{}, err error) {
	return this.Do("SET", key, value)
}

func (this *Session) SETEX(key string, value interface{}, seconds int) (reply interface{}, err error) {
	return this.Do("SETEX", key, seconds, value)
}

func (this *Session) INCR(key string) (reply interface{}, err error) {
	return this.Do("INCR", key)
}

////////////////////////////////////////////////////////////////////////////////
func (this *Session) HMSET(key string, obj interface{}) (reply interface{}, err error) {
	return this.Do("HMSET", redigo.Args{}.Add(key).AddFlat(obj)...)
}

func (this *Session) HGETALL(key string, obj interface{}) (err error) {
	var reply interface{}
	reply, err = this.Do("HGETALL", key)
	if err != nil {
		return err
	}

	err = redigo.ScanStruct(reply.([]interface{}), obj)
	return err
}

func (this *Session) HLEN(key string) (int) {
	var r, _ = Int(this.Do("HLEN", key))
	return r
}

////////////////////////////////////////////////////////////////////////////////
// 把一个对象编码成 JSON 字符串数据进行存储
func (this *Session) EncodeToJSONEX(key string, obj interface{}, seconds int) (reply interface{}, err error) {
	value, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	reply, err = this.SETEX(key, value, seconds)
	return reply, err
}

////////////////////////////////////////////////////////////////////////////////
func (this *Session) Transaction(callback func(conn Conn)) (reply interface{}, err error) {
	var c = this.c
	c.Send("MULTI")
	callback(c)
	return c.Do("EXEC")
}

func (this *Session) Pipeline(callback func(conn Conn)) (err error) {
	var c = this.c
	callback(c)
	return c.Flush()
}

////////////////////////////////////////////////////////////////////////////////
func Int(reply interface{}, err error) (int, error) {
	return redigo.Int(reply, err)
}

func Bool(reply interface{}, err error) (bool, error) {
	return redigo.Bool(reply, err)
}

func String(reply interface{}, err error) (string, error) {
	return redigo.String(reply, err)
}

func MustInt(reply interface{}, err error) (int) {
	var r, _ = Int(reply, err)
	return r
}

func MustBool(reply interface{}, err error) (bool) {
	var r, _ = Bool(reply, err)
	return r
}

func MustString(reply interface{}, err error) (string) {
	if err != nil {
		fmt.Println(err)
	}
	var r, _ = String(reply, err)
	return r
}

////////////////////////////////////////////////////////////////////////////////
type Conn interface {
	redigo.Conn
}

////////////////////////////////////////////////////////////////////////////////
const k_REDIS_KEY = "redis_conn"

type Setter interface {
	Set(key string, value interface{})
}

type Getter interface {
	MustGet(key string) interface{}
}

func FromContext(g Getter) *Session {
	return g.MustGet(k_REDIS_KEY).(*Session)
}

func ToContext(s Setter, c *Session) {
	s.Set(k_REDIS_KEY, c)
}