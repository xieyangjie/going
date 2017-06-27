package redis

import (
	"encoding/json"
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"os"
	"time"
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
	return NewSession(c)
}

func (this *Pool) Release(s *Session) {
	s.c.Close()
}

////////////////////////////////////////////////////////////////////////////////
func NewSession(c Conn) *Session {
	if c == nil {
		return nil
	}
	return &Session{c: c}
}

type Session struct {
	c Conn
}

func (this *Session) Conn() redigo.Conn {
	return this.c
}

func (this *Session) Do(commandName string, args ...interface{}) (interface{}, error) {
	return this.c.Do(commandName, args...)
}

////////////////////////////////////////////////////////////////////////////////
func (this *Session) EXISTS(key string) bool {
	return MustBool(this.Do("EXISTS", key))
}

func (this *Session) EXPIRE(key string, seconds int) (interface{}, error) {
	return this.Do("EXPIRE", key, seconds)
}

func (this *Session) DEL(key ...interface{}) (interface{}, error) {
	return this.Do("DEL", key...)
}

////////////////////////////////////////////////////////////////////////////////
func (this *Session) GET(key string) (interface{}, error) {
	return this.Do("GET", key)
}

func (this *Session) SET(key string, value interface{}) (interface{}, error) {
	return this.Do("SET", key, value)
}

func (this *Session) SETEX(key string, value interface{}, seconds int) (interface{}, error) {
	return this.Do("SETEX", key, seconds, value)
}

////////////////////////////////////////////////////////////////////////////////
func (this *Session) INCRBY(key string, increment int) (interface{}, error) {
	return this.Do("INCRBY", key, increment)
}

func (this *Session) INCR(key string) (interface{}, error) {
	return this.Do("INCR", key)
}

func (this *Session) DECR(key string) (interface{}, error) {
	return this.Do("DECR", key)
}

////////////////////////////////////////////////////////////////////////////////
func (this *Session) SADD(key string, member interface{}) (interface{}, error) {
	return this.Do("SADD", key, member)
}

func (this *Session) SREM(key string, member interface{}) (interface{}, error) {
	return this.Do("SREM", key, member)
}

func (this *Session) SMEMBERS(key string) (interface{}, error) {
	return this.Do("SMEMBERS", key)
}

func (this *Session) SCARD(key string) int {
	return MustInt(this.Do("SCARD", key))
}

////////////////////////////////////////////////////////////////////////////////
func (this *Session) HEXISTS(key, field string) bool {
	return MustBool(this.Do("HEXISTS", key, field))
}

func (this *Session) HSET(key, field string, value interface{}) (interface{}, error) {
	return this.Do("HSET", key, field, value)
}

func (this *Session) HINCRBY(key, field string, increment int) (interface{}, error) {
	return this.Do("HINCRBY", key, field, increment)
}

func (this *Session) HMSET(key string, obj interface{}) (interface{}, error) {
	return this.Do("HMSET", redigo.Args{}.Add(key).AddFlat(obj)...)
}

func (this *Session) HGET(key string, field string) (reply interface{}, err error) {
	return this.Do("HGET", key, field)
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

func (this *Session) HLEN(key string) int {
	var r, _ = Int(this.Do("HLEN", key))
	return r
}

////////////////////////////////////////////////////////////////////////////////

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
func (this *Session) Int(reply interface{}, err error) (int, error) {
	return redigo.Int(reply, err)
}

func (this *Session) Bool(reply interface{}, err error) (bool, error) {
	return redigo.Bool(reply, err)
}

func (this *Session) String(reply interface{}, err error) (string, error) {
	return redigo.String(reply, err)
}

func (this *Session) Strings(reply interface{}, err error) ([]string, error) {
	return redigo.Strings(reply, err)
}

func (this *Session) Float64(reply interface{}, err error) (float64, error) {
	return redigo.Float64(reply, err)
}

func (this *Session) MustInt(reply interface{}, err error) int {
	var r, _ = Int(reply, err)
	return r
}

func (this *Session) MustBool(reply interface{}, err error) bool {
	var r, _ = Bool(reply, err)
	return r
}

func (this *Session) MustString(reply interface{}, err error) string {
	var r, _ = String(reply, err)
	return r
}

func (this *Session) MustStrings(reply interface{}, err error) ([]string) {
	var r, _ = Strings(reply, err)
	return r
}

func (this *Session) MustFloat64(reply interface{}, err error) float64 {
	var r, _ = Float64(reply, err)
	return r
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

func Strings(reply interface{}, err error) ([]string, error) {
	return redigo.Strings(reply, err)
}

func Float64(reply interface{}, err error) (float64, error) {
	return redigo.Float64(reply, err)
}

func MustInt(reply interface{}, err error) int {
	var r, _ = Int(reply, err)
	return r
}

func MustBool(reply interface{}, err error) bool {
	var r, _ = Bool(reply, err)
	return r
}

func MustString(reply interface{}, err error) string {
	var r, _ = String(reply, err)
	return r
}

func MustStrings(reply interface{}, err error) ([]string) {
	var r, _ = Strings(reply, err)
	return r
}

func MustFloat64(reply interface{}, err error) float64 {
	var r, _ = Float64(reply, err)
	return r
}

////////////////////////////////////////////////////////////////////////////////
type Conn interface {
	redigo.Conn
}

////////////////////////////////////////////////////////////////////////////////
const k_REDIS_KEY = "redis_conn"

type Context interface {
	Set(key string, value interface{})

	MustGet(key string) interface{}
}

func FromContext(ctx Context) *Session {
	return ctx.MustGet(k_REDIS_KEY).(*Session)
}

func ToContext(ctx Context, s *Session) {
	ctx.Set(k_REDIS_KEY, s)
}
