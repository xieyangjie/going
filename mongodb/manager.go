package mongodb

import (
	"gopkg.in/mgo.v2/bson"
)

////////////////////////////////////////////////////////////////////////////////
// 添加数据
func (this *Session) Insert(cName string, obj ...interface{}) (err error) {
	var c = this.C(cName)
	err = c.Insert(obj...)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// UpdateWithId 更新指定id的文档
func (this *Session) UpdateWithId(cName string, id bson.ObjectId, update bson.M) (err error) {
	var c = this.C(cName)
	err = c.UpdateId(id, update)
	return err
}

// Update 更新满足条件的第一个文档
func (this *Session) UpdateOne(cName string, selector, update bson.M) (err error) {
	var c = this.C(cName)
	err = c.Update(selector, update)
	return err
}

// UpdateAll 更新满足所有满足条件的文档
func (this *Session) UpdateAll(cName string, selector, update bson.M) (err error) {
	var c = this.C(cName)
	_, err = c.UpdateAll(selector, update)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// FindWithId 查询指定id的文档
func (this *Session) FindWithId(cName string, id bson.ObjectId, result interface{}) (err error) {
	var c = this.C(cName)
	var query = bson.M{"_id": id}
	err = c.Find(query).One(result)
	return err
}

// FindOne 查询满足条件的第一个文档
func (this *Session) FindOne(cName string, query bson.M, result interface{}) (err error) {
	var c = this.C(cName)
	err = c.Find(query).One(result)
	return err
}

// FindAll 查询满足条件的所有文档
func (this *Session) FindAll(cName string, query bson.M, result interface{}) (err error) {
	var c = this.C(cName)
	err = c.Find(query).All(result)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// RemoveWithId 物理删除指定id的文档
func (this *Session) RemoveWithId(cName string, id bson.ObjectId) (err error) {
	var c = this.C(cName)
	err = c.Remove(bson.M{"_id": id})
	return err
}

// Remove 物理删除满足条件的第一个文档
func (this *Session) RemoveOne(cName string, selector bson.M) (err error) {
	var c = this.C(cName)
	err = c.Remove(selector)
	return err
}

// RemoveAll 物理删除满足条件的所有文档
func (this *Session) RemoveAll(cName string, selector bson.M) (err error) {
	var c = this.C(cName)
	_, err = c.RemoveAll(selector)
	return err
}