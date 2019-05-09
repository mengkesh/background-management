package controllers

import (
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego"
)
func init(){
	conn,err:=redis.Dial("tcp",":6379")
	if err != nil {
		beego.Error("连接redis失败")
		return
	}
	defer conn.Close()
	//conn.Do("set","c1","hello world")
	//resp,err:=conn.Do("get","c1")
	//result,err:=redis.String(resp,err)
	//if err!=nil{
	//	beego.Error("获取数据失败")
	//	return
	//}
	//beego.Info("获取的数据为：",result)
	resp,err:=conn.Do("mget","r1","r2","r3")
	result,err:=redis.Values(resp,err)
	if err != nil {
		beego.Error("获取数据失败")
	}
	var v1 int
	var v2 string
	var v3 float64
	redis.Scan(result,&v1,&v2,&v3)
	beego.Info("获取的数据为：",v1,v2,v3)
}

