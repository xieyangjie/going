package redis

func (this *Session) SADD(key string, member interface{}) (interface{}, error) {
	return this.Do("SADD", key, member)
}

func (this *Session) SCARD(key string) int64 {
	return MustInt64(this.Do("SCARD", key))
}

func (this *Session) SMEMBERS(key string) (interface{}, error) {
	return this.Do("SMEMBERS", key)
}

func (this *Session) SREM(key string, member interface{}) (interface{}, error) {
	return this.Do("SREM", key, member)
}
