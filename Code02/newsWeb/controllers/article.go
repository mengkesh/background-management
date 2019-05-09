package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"Code02/newsWeb/models"
	"math"
	"strconv"
	"github.com/gomodule/redigo/redis"
	"bytes"
	"encoding/gob"
)

type ArticleController struct {
	beego.Controller
}
//=================================================展示文章列表===============================================================
func (this *ArticleController) ShowIndex() {
	username:=this.GetSession("userName")
	if username==nil{
		this.Redirect("/login",302)
		return
	}
	this.Data["userName"]=username.(string)
	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var articles []models.Article
	//qs.All(&articles)
	/*==================================文章总数
	amountNum, _ := qs.Count()
	this.Data["amountnum"] = amountNum
	==================================总页数*/
	singlePageNum := 2
	/*PageCount := int(math.Ceil(float64(amountNum) / float64(singlePageNum)))
	this.Data["pagecount"] = PageCount*/
	//==================================获取当前页数和当前页信息
	pageNum, err := this.GetInt("pagenum")
	if err != nil {
		pageNum = 1
	}
	var amountNum int64
	typename:=this.GetString("select")
	if typename=="" {
		amountNum, _ = qs.RelatedSel("ArticleType").Count()
		_, err = qs.Limit(singlePageNum, (pageNum-1)*singlePageNum).RelatedSel("ArticleType").All(&articles)
	}else{
		amountNum, _ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typename).Count()
		_, err = qs.Limit(singlePageNum, (pageNum-1)*singlePageNum).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typename).All(&articles)
	}

	var a []models.ArticleType
	o.QueryTable("ArticleType").Filter("Articles__id",3).All(&a)
	beego.Info(a)




	if err != nil&&amountNum==0 {
		beego.Info("没有数据")
	}else if err != nil {
		beego.Error("获取数据失败", err)
		this.Layout="layout.html"
		this.TplName = "index.html"
		return
	}
	this.Data["amountnum"] = amountNum
	PageCount := int(math.Ceil(float64(amountNum) / float64(singlePageNum)))
	this.Data["pagecount"] = PageCount



	//if err != nil {
	//	beego.Error("获取数据失败", err)
	//	this.Data["currentPageNum"] = pageNum
	//	this.TplName = "index.html"
	//	return
	//}
	var articletypes []models.ArticleType
	conn,err:=redis.Dial("tcp","127.0.0.1:6379")
	if err != nil {
		beego.Error("数据库连接失败")
		return
	}
	defer conn.Close()
	resp,err:=conn.Do("get","newsweb")
	result,_:=redis.Bytes(resp,err)
	if len(result) == 0 {
		o.QueryTable("ArticleType").All(&articletypes)
		var buffer bytes.Buffer
		enc:=gob.NewEncoder(&buffer)
		enc.Encode(articletypes)
		conn.Do("set","newsweb",buffer.Bytes())
		beego.Info("从mysql中获取")
	}else {
		dec:=gob.NewDecoder(bytes.NewBuffer(result))
		dec.Decode(&articletypes)
		beego.Info(articletypes)
		beego.Info("从redis中获取")
	}








	this.Data["typename"]=typename
	this.Data["articletypes"]=articletypes
	this.Data["articles"] = articles
	this.Data["currentPageNum"] = pageNum
	this.Layout="layout.html"
	this.TplName = "index.html"
}
//=================================================添加新闻===============================================================
func (this *ArticleController) ShowAddArticle() {
	o:=orm.NewOrm()
	var articletypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articletypes)
	this.Data["articletypes"]=articletypes
	this.TplName = "add.html"
}
func (this *ArticleController) AddArticle() {
	title := this.GetString("articleName")
	content := this.GetString("content")
	thisarticletype:=this.GetString("select")
	if title == "" || content == ""||thisarticletype=="" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "上传数据错误"
		this.TplName = "add.html"
		return
	}
	o := orm.NewOrm()
	var articletype models.ArticleType
	articletype.TypeName=thisarticletype
	err:=o.Read(&articletype,"TypeName")
	if err!=nil{
		beego.Error("读取数据错误")
		this.Redirect("/article/addarticle",302)
		return
	}
	var article models.Article
	article.Title = title
	article.Content = content
	article.ArticleType=&articletype

	file, head, err := this.GetFile("uploadname")
	if head.Size != 0 {
		if err != nil {
			beego.Error("获取图片错误")
			this.Data["errmsg"] = "上传图片错误"
			this.TplName = "add.html"
			return
		}
		defer file.Close()
		if head.Size >= 5000000 {
			beego.Error("获取图片错误")
			this.Data["errmsg"] = "上传图片过大，请重新上传"
			this.TplName = "add.html"
			return
		}
		ext := path.Ext(head.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			beego.Error("获取图片错误")
			this.Data["errmsg"] = "上传图片格式错误，请重新上传"
			this.TplName = "add.html"
			return
		}
		filename := time.Now().Format("200601021504051111")
		this.SaveToFile("uploadname", "./static/img/"+filename+ext)
		article.Img = "/static/img/" + filename + ext
	}

	_, err = o.Insert(&article)
	if err != nil {
		beego.Error("写入数据库失败", err)
		this.Data["errmsg"] = "写入数据库失败"
		this.TplName = "add.html"
		return
	}
	this.Redirect("/article/index", 302)
}

//=================================================文章详细信息===============================================================
func (this *ArticleController) ShowArticleDetail() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error(err)
		this.Redirect("/article/index",302)
		return
	}
	o:=orm.NewOrm()
	//qs:=o.QueryTable("Article")
	var article models.Article
	article.Id=id
	o.Read(&article)

	//多对多查询 方法二
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)
	//											   字段名     表名
this.Data["users"]=users


	article.Readcount+=1
	_,err=o.Update(&article,"Readcount")
	if err != nil {
		beego.Error(err)
		this.Redirect("/article/index",302)
		return
	}
	username:=this.GetSession("userName")
	var user models.User
	user.Name=username.(string)
	o.Read(&user,"Name")
	m2m:=o.QueryM2M(&article,"Users")
	m2m.Add(user)
	this.Data["article"]=article
	this.TplName = "content.html"
}
//=================================================编辑文章===============================================================
func (this *ArticleController)ShowUpdateArticle(){
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error(err)
		this.Redirect("/article/index",302)
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.Read(&article)
	this.Data["article"]=article
	this.TplName="update.html"
}

func SaveFile(this *ArticleController,uploadname string,errpath string)string {
	file, head, err := this.GetFile(uploadname)
	var name string
	if head.Size != 0 {
		if err != nil {
			beego.Error("获取图片错误")
			this.Data["errmsg"] = "上传图片错误"
			this.TplName = errpath
			return name
		}
		defer file.Close()
		if head.Size >= 5000000 {
			beego.Error("获取图片错误")
			this.Data["errmsg"] = "上传图片过大，请重新上传"
			this.TplName = errpath
			return name
		}
		ext := path.Ext(head.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			beego.Error("获取图片错误")
			this.Data["errmsg"] = "上传图片格式错误，请重新上传"
			this.TplName = errpath
			return name
		}
		filename := time.Now().Format("200601021504051111")
		this.SaveToFile("uploadname", "./static/img/"+filename+ext)
		name="/static/img/"+filename+ext
	} else{
		name="none"
	}
	return name
}
func (this *ArticleController)HandleUpdate(){
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error(err)
		this.Redirect("/article/index",302)
		return
	}
	title:=this.GetString("articleName")
	content:=this.GetString("content")
	filename:=SaveFile(this,"uploadname","update.html")
	if title==""||content==""||filename=="" {
		beego.Error("获取数据失败")
		this.Redirect("/article/updatearticle?id="+strconv.Itoa(id),302)
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.Read(&article)
	article.Title=title
	article.Content=content
	if filename=="none"{
		article.Img=""
	}else{
		article.Img=filename
	}
	o.Update(&article)
	this.Redirect("/article/index",302)
}
//=================================================删除文章===============================================================
func(this *ArticleController)HandleDeleteArticle(){
	id,err:=this.GetInt("id")
	if err != nil {
		beego.Error(err)
		this.Redirect("/article/index",302)
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.Delete(&article)
	this.Redirect("/article/index",302)
}
//=================================================添加文章分类===============================================================
func (this *ArticleController)ShowAddArticleType(){
	o:=orm.NewOrm()
	qs:=o.QueryTable("ArticleType")
	var articletypes []models.ArticleType
	qs.All(&articletypes)
	this.Data["articletypes"]=articletypes
	this.TplName="addType.html"
}
func (this *ArticleController)HandleAddArticleType(){
	typename:=this.GetString("typeName")
	if typename=="" {
		beego.Error("输入类型不正确")
		this.Redirect("/article/addarticletype",302)
		return
	}
	o:=orm.NewOrm()
	var articletype models.ArticleType
	articletype.TypeName=typename
	_,err:=o.Insert(&articletype)
	if err!=nil{
		beego.Error(err)
		this.Redirect("/article/addarticletype",302)
	}
	this.Redirect("/article/addarticletype",302)
}
//=================================================删除文章分类===============================================================
func(this *ArticleController)HandleDeleteArticleType(){
	id,err:=this.GetInt("id")
	if err != nil {
		beego.Error(err)
		this.Redirect("/article/addarticletype",302)
		return
	}
	o:=orm.NewOrm()
	var articletype models.ArticleType
	count,_:=o.QueryTable("Article").Filter("ArticleType__Id",id).Count()
	if count!=0{
		beego.Error("您要删除的类型有文章使用，无法删除")
		this.Redirect("/article/addarticletype",302)
		return
	}
	articletype.Id=id
	o.Delete(&articletype)
	this.Redirect("/article/addarticletype",302)
}