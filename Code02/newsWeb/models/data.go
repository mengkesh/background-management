package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)
type User struct {
	Id int
	Name string
	Pwd string
	Articles []*Article `orm:"reverse(many)"`
}
type Article struct{
	Id int `orm:"pk;auto"`
	Title string `orm:"unique;size(40)"`
	Content string `orm:"size(500)"`
	Img string `orm:"null"`
	Time time.Time `orm:"type(datetime);auto_now"`
	Readcount int `orm:"default(0)"`
	ArticleType *ArticleType `orm:"rel(fk);null"`
	Users []*User `orm:"rel(m2m)"`
}
type ArticleType struct {
	Id int
	TypeName string `orm:"unique"`
	Articles []*Article `orm:"reverse(many)"`
}
func init(){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/news?charset=utf8&loc=Asia%2FShanghai")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}
