package mongodb

import (
	"fmt"
	"github.com/smartwalle/pool"
	"gopkg.in/mgo.v2"
	"time"
)

func NewMongoDB(dbURL, dbName string, maxActive, maxIdle int) (p *Pool) {
	var s, err = mgo.Dial(dbURL)
	if err != nil {
		fmt.Println("连接数据库失败:", dbURL, dbName, err)
		return nil
	}

	var np = &pool.Pool{}
	np.Dial = dialFunc(s, dbName)
	np.MaxIdle = maxIdle
	np.MaxActive = maxActive
	np.IdleTimeout = time.Minute * 10
	np.Wait = true

	p = &Pool{np}
	return p
}

func dialFunc(s *mgo.Session, dbName string) func() (pool.IConnection, error) {
	var dialFunc = func() (pool.IConnection, error) {
		var ns = &Session{}
		ns.dbName = dbName
		ns.session = s.Clone()
		return ns, nil
	}
	return dialFunc
}

////////////////////////////////////////////////////////////////////////////////
type Pool struct {
	p *pool.Pool
}

func (this *Pool) GetSession() *Session {
	var c, err = this.p.Get()
	if err == nil {
		var ns = c.(*Session)
		return ns
	}
	return nil
}

func (this *Pool) Release(s *Session) {
	this.p.Release(s, false)
}

////////////////////////////////////////////////////////////////////////////////
type Session struct {
	session *mgo.Session
	dbName  string
}

// 不能主动调用此方法
func (this *Session) Close() {
	this.session.Close()
	this.dbName = ""
	this.session = nil
}

func (this *Session) C(name string) *mgo.Collection {
	return this.session.DB(this.dbName).C(name)
}

////////////////////////////////////////////////////////////////////////////////
const k_MONGODB_KEY = "mongo_session"

type Setter interface {
	Set(key string, value interface{})
}

type Getter interface {
	MustGet(key string) interface{}
}

func FromContext(g Getter) *Session {
	return g.MustGet(k_MONGODB_KEY).(*Session)
}

func ToContext(s Setter, c *Session) {
	s.Set(k_MONGODB_KEY, c)
}

////////////////////////////////////////////////////////////////////////////////
//type session struct {
//	wrapper *sessionWrapper
//}
//
//func (this *session) Release() {
//	this.wrapper.release()
//}
//
//func (this *session) C(name string) *mgo.Collection {
//	return this.wrapper.session.DB(this.wrapper.dbName).C(name)
//}

//package mongodb
//
//import (
//"fmt"
//"github.com/smartwalle/pool"
//"gopkg.in/mgo.v2"
//"os"
//"time"
//)
//
//var mainSession *session
//
//func InitMongoDB(dbURL, dbName string, maxActive, maxIdle int) {
//	var s, err = mgo.Dial(dbURL)
//	if err != nil {
//		fmt.Println("连接数据库失败:", dbURL, dbName, err)
//		os.Exit(-1)
//	}
//
//	mainSession = &session{session: s}
//	mainSession.dbName = dbName
//
//	var p = &pool.Pool{}
//	p.Dial = dialFunc(mainSession)
//	p.MaxIdle = maxIdle
//	p.MaxActive = maxActive
//	p.IdleTimeout = time.Minute * 10
//	p.Wait = true
//	mainSession.pool = p
//}
//
//func dialFunc(s *session) func() (pool.IConnection, error) {
//	var dialFunc = func() (pool.IConnection, error) {
//		return s.clone(), nil
//	}
//	return dialFunc
//}
//
//func NewSession() *session {
//	var c, err = mainSession.pool.Get()
//	if err == nil {
//		var s = c.(*session)
//		return s
//	}
//	return nil
//}
//
//////////////////////////////////////////////////////////////////////////////////
//type session struct {
//	session *mgo.Session
//	dbName  string
//	pool    *pool.Pool
//}
//
//func (this *session) clone() *session {
//	var s = &session{session: this.session.Clone()}
//	s.dbName = this.dbName
//	s.pool = this.pool
//	return s
//}
//
//// 该方法提供给 pool 调用,  不能显示调用此方法. session 使用完成之后, 应该调用 Release 进行归还.
//func (this *session) Close() {
//	this.session.Close()
//	this.dbName = ""
//	this.pool = nil
//	this.session = nil
//}
//
//func (this *session) Release() {
//	this.pool.Release(this, false)
//}
//
//func (this *session) C(name string) *mgo.Collection {
//	return this.session.DB(this.dbName).C(name)
//}
