package controllers

import (
	"beego_blog/models"

	"github.com/astaxie/beego"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Get() {
	c.TplName = "register.html"
}

func (ctrl *RegisterController) Post() {

	uname := ctrl.Input().Get("uname")
	psw := ctrl.Input().Get("psw")

	err := models.AddUser(uname, psw)
	if err != nil {
		beego.Error(err)
	}

	if beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("psw") == psw {
		maxAge := 1<<31 - 1
		// 注册成功，放入 cookie
		ctrl.Ctx.SetCookie("uname", uname, maxAge, "/")
		ctrl.SetSession("uname", uname)

	}
	ctrl.Redirect("/", 301)
	return
}
