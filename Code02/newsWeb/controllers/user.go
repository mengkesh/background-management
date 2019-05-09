package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"Code02/newsWeb/models"
	"encoding/base64"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister() {

	this.TplName = "register.html"
}
func (this *UserController) Register() {
	name := this.GetString("userName")
	pwd := this.GetString("password")
	if name == "" || pwd == "" {
		beego.Error("输入数据不完整")
		this.TplName = "register.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = name
	user.Pwd = pwd
	id, err := o.Insert(&user)
	if err != nil {
		beego.Error("写入数据库错误")
		this.TplName = "register.html"
		return
	}
	beego.Info(id)
	this.Redirect("/login",302)
}
func (this *UserController) ShowLogin() {
	userName:=this.Ctx.GetCookie("userName")
	dec,_:=base64.StdEncoding.DecodeString(userName)
	if userName!="" {
		this.Data["username"]=string(dec)
		this.Data["checked"]="checked"
	}else{
		this.Data["username"]=""
		this.Data["checked"]=""
	}
	this.TplName = "login.html"
}
func (this *UserController) Login() {
	name := this.GetString("userName")
	pwd := this.GetString("password")
	if name == "" || pwd == "" {
		beego.Error("输入数据不完整")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = name
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户名错误")
		this.TplName = "login.html"
		return
	}
	if user.Pwd != pwd {
		beego.Error("密码错误")
		this.TplName = "login.html"
		return
	}
	remember:=this.GetString("remember")
	enc:=base64.StdEncoding.EncodeToString([]byte(name))
	if remember=="on"{
		this.Ctx.SetCookie("userName",enc,600)
	}else{
		this.Ctx.SetCookie("userName",name,-1)
	}
	this.SetSession("userName",name)
	this.Redirect("/article/index",302)
}
func (this *UserController)LogOut(){
	this.DelSession("userName")
	this.Redirect("/login",302)
}
