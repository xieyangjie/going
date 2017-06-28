package redis

import (
	//redigo "github.com/garyburd/redigo/redis"
)

func (this *Session) HEXISTS(key, field string) bool {
	return MustBool(this.Do("HEXISTS", key, field))
}

func (this *Session) HGET(key string, field string) (reply interface{}, err error) {
	return this.Do("HGET", key, field)
}

func (this *Session) HGETALL(key string) (interface{}, error) {
	return this.Do("HGETALL", key)
}

//func (this *Session) HGETALL(key string, obj interface{}) (err error) {
//	var reply interface{}
//	reply, err = this.Do("HGETALL", key)
//	if err != nil {
//		return err
//	}
//
//	err = redigo.ScanStruct(reply.([]interface{}), obj)
//	return err
//}

func (this *Session) HINCRBY(key, field string, increment int) (interface{}, error) {
	return this.Do("HINCRBY", key, field, increment)
}

func (this *Session) HLEN(key string) int64 {
	var r, _ = Int64(this.Do("HLEN", key))
	return r
}

func (this *Session) HMSET(key string, params ...interface{}) (interface{}, error) {
	var ps []interface{}
	ps = append(ps, key)
	ps = append(ps, params...)
	return this.Do("HMSET", ps...)
}

//func (this *Session) HMSET(key string, obj interface{}) (interface{}, error) {
//	return this.Do("HMSET", redigo.Args{}.Add(key).AddFlat(obj)...)
//}

func (this *Session) HSET(key, field string, value interface{}) (interface{}, error) {
	return this.Do("HSET", key, field, value)
}