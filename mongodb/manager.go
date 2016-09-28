package mongodb

import (
	"gopkg.in/mgo.v2/bson"
	"strings"
)

const (
	//获取较新的数据
	DB_GET_NEWEST = 1

	//获取较旧的数据
	DB_GET_PAST = 2

	//默认分页数量
	DB_PAGE_SIZE = 20
)

type IPagination interface {
	GetPageId() string
	GetPageFlag() int

	GetPageSize() int
	GetPageNumber() int
	GetSortFields() []string
}

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

func (this *Session) FindAllWithPagination(cName string, query bson.M, pagination IPagination, results interface{}) (err error) {
	var pageId = pagination.GetPageId()
	var pageFlag = pagination.GetPageFlag()
	var pageSize = pagination.GetPageSize()
	var pageNumber = pagination.GetPageNumber()
	var sortFields = pagination.GetSortFields()
	if bson.IsObjectIdHex(pageId) && pageNumber <= 0 {
		if pageFlag == DB_GET_NEWEST {
			query["_id"] = bson.M{"$gt": bson.ObjectIdHex(pageId)}
		} else if pageFlag == DB_GET_PAST {
			query["_id"] = bson.M{"$lt": bson.ObjectIdHex(pageId)}
		}
	}

	if pageSize == 0 {
		pageSize = DB_PAGE_SIZE
	}

	var fieldsList = make([]string, 0, len(sortFields))
	for _, value := range sortFields {
		v := strings.TrimSpace(value)
		if len(v) > 0 {
			fieldsList = append(fieldsList, v)
		}
	}
	var collection = this.C(cName)
	var q = collection.Find(query).Sort(fieldsList...)
	if pageNumber > 0 {
		q = q.Skip((pageNumber - 1) * pageSize)
	}

	if pageSize > 0 {
		q = q.Limit(pageSize)
	}

	err = q.All(results)
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
