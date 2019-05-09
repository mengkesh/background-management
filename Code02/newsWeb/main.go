package main

import (
	_"Code02/newsWeb/controllers"
	_ "Code02/newsWeb/routers"
	"github.com/astaxie/beego"
	_"Code02/newsWeb/models"
)
func NextPage (pageindex,pagecount int) int{
	if pageindex>=pagecount{
		return pagecount
	}
	return pageindex+1
}
func PrePage (pageindex int) int{
	if pageindex<=1{
		return 1
	}
	return pageindex-1
}
func AddOne (key int)int{
	return key+1
}
func main() {
	beego.AddFuncMap("nextpage",NextPage)
	beego.AddFuncMap("prepage",PrePage)
	beego.AddFuncMap("addone",AddOne)
	beego.Run()
}

