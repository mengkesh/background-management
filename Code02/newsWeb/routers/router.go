package routers

import (
	"Code02/newsWeb/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeRouter,filterFunc)
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register", &controllers.UserController{},"get:ShowRegister;post:Register")
    beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:Login")
    beego.Router("/article/index", &controllers.ArticleController{},"get,post:ShowIndex")
	beego.Router("/article/addarticle", &controllers.ArticleController{},"get:ShowAddArticle;post:AddArticle")
	beego.Router("/article/showarticledetail", &controllers.ArticleController{},"get:ShowArticleDetail")
	beego.Router("/article/updatearticle", &controllers.ArticleController{},"get:ShowUpdateArticle;post:HandleUpdate")
	beego.Router("/article/deletearticle", &controllers.ArticleController{},"get:HandleDeleteArticle")
	beego.Router("/article/addarticletype", &controllers.ArticleController{},"get:ShowAddArticleType;post:HandleAddArticleType")
	beego.Router("/article/deletearticletype", &controllers.ArticleController{},"get:HandleDeleteArticleType")
	beego.Router("/article/logout",&controllers.UserController{},"get:LogOut")
}
func filterFunc(ctx *context.Context){
	userName:=ctx.Input.Session("userName")
	if userName==nil{
		ctx.Redirect(302,"/login")
		return
	}
}
