package redis

//DEL 删除给定的一个或多个 key 。
func (this *Session) DEL(keys ...string) (interface{}, error) {
	var ks []interface{}
	for _, k := range keys {
		ks = append(ks, k)
	}
	return this.Do("DEL", ks...)
}

//EXISTS 检查给定 key 是否存在。
func (this *Session) EXISTS(key string) bool {
	return MustBool(this.Do("EXISTS", key))
}

//EXPIRE 为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除。
func (this *Session) EXPIRE(key string, seconds int) (interface{}, error) {
	return this.Do("EXPIRE", key, seconds)
}

//EXPIREAT 作用和 EXPIRE 类似，都用于为 key 设置生存时间。不同在于 EXPIREAT 命令接受的时间参数是 UNIX 时间戳(unix timestamp)。
func (this *Session) EXPIREAT(key string, timestamp int64) (interface{}, error) {
	return this.Do("EXPIREAT", key, timestamp)
}

//KEYS 查找所有符合给定模式 pattern 的 key 。
func (this *Session) KEYS(pattern string) (interface{}, error) {
	return this.Do("KEYS", pattern)
}

//MIGRATE 将 key 原子性地从当前实例传送到目标实例的指定数据库上，一旦传送成功， key 保证会出现在目标实例上，而当前实例上的 key 会被删除。
func (this *Session) MIGRATE(host, port, key string, destinationDB int, timeout int, options string) (interface{}, error) {
	return this.Do("MIGRATE", host, port, key, destinationDB, timeout, options)
}

//MOVE 将当前数据库的 key 移动到给定的数据库 db 当中。
func (this *Session) MOVE(key string, destinationDB int) (interface{}, error) {
	return this.Do("MOVE", key, destinationDB)
}

//PERSIST 移除给定 key 的生存时间，将这个 key 从『易失的』(带生存时间 key )转换成『持久的』(一个不带生存时间、永不过期的 key)。
func (this *Session) PERSIST(key string) (interface{}, error) {
	return this.Do("PERSIST", key)
}

//PEXPIRE 这个命令和 EXPIRE 命令的作用类似，但是它以毫秒为单位设置 key 的生存时间，而不像 EXPIRE 命令那样，以秒为单位。
func (this *Session) PEXPIRE(key string, milliseconds int) (interface{}, error) {
	return this.Do("PEXPIRE", key, milliseconds)
}

//RANDOMKEY 从当前数据库中随机返回(不删除)一个 key
func (this *Session) RANDOMKEY() (string, error) {
	return this.String(this.Do("RANDOMKEY"))
}

//RENAME 将 key 改名为 newKey。
func (this *Session) RENAME(key, newKey string) (interface{}, error) {
	return this.Do("RENAME", key, newKey)
}

//RENAMENX 当且仅当 newkey 不存在时，将 key 改名为 newkey。
func (this *Session) RENAMENX(key, newKey string) (interface{}, error) {
	return this.Do("RENAMENX", key, newKey)
}

//SORT 返回或保存给定列表、集合、有序集合 key 中经过排序的元素。
func (this *Session) SORT(key string, params ...interface{}) (interface{}, error) {
	var ps []interface{}
	ps = append(ps, key)
	ps = append(ps, params...)
	return this.Do("SORT", ps...)
}

//TTL 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。
func (this *Session) TTL(key string) (int64) {
	return MustInt64(this.Do("TTL", key))
}

//TYPE 返回 key 所储存的值的类型。
func (this *Session) TYPE(key string) (interface{}, error) {
	return this.Do("TYPE", key)
}