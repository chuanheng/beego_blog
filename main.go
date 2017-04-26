package main

import (
	"beego_blog/controllers"
	"beego_blog/models"
	_ "beego_blog/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
)

func init() {
	models.RegisterDB()
}

func main() {
	i18n.SetMessage("en-US", "conf/locale_en-US.ini")
	i18n.SetMessage("zh-CN", "conf/locale_zh-CN.ini")

	// 注册国际化模板函数
	beego.AddFuncMap("i18n", i18n.Tr)

	orm.Debug = true
	orm.RunSyncdb("default", false, true)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/register", &controllers.RegisterController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/topic", &controllers.TopicController{})
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/reply", &controllers.ReplyController{})
	beego.AutoRouter(&controllers.ReplyController{})
	beego.Router("/attachment/:all", &controllers.AttachController{})
	// 启动 beego
	beego.Run()
}
