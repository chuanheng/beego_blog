package controllers

import (
	"beego_blog/models"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	// 如果是退出登录，清空 Cookie 和 Session
	if c.GetString("exist") != "" {
		c.Ctx.SetCookie("uname", "", -1, "/")
		c.DelSession("uname")
	}
	c.TplName = "login.html"
}

func (ctrl *LoginController) Post() {

	uname := ctrl.Input().Get("uname")
	psw := ctrl.Input().Get("psw")
	autoLogin := ctrl.Input().Get("autoLogin") == "on"

	user, err := models.GetUser(uname)
	if err != nil {
		beego.Error(err)
		ctrl.Redirect("/login", 302)
		return
	}

	if user.PassWord == psw {
		maxAge := 0
		if autoLogin {
			maxAge = 1<<31 - 1
		}
		// 登录成功，放入 cookie
		ctrl.Ctx.SetCookie("uname", uname, maxAge, "/")
		ctrl.SetSession("uname", uname)
		ctrl.Redirect("/", 301)
	} else {
		ctrl.Redirect("/login", 301)
	}

	return
}

// 检查账户
func checkAccount(ctx *context.Context) bool {
	fmt.Println("current user :", ctx.Input.CruSession.Get("uname"))
	if ctx.Input.CruSession.Get("uname") != "" &&
		ctx.Input.CruSession.Get("uname") != nil {
		return true
	}
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		return false
	}
	uname := ck.Value
	if beego.AppConfig.String("uname") == uname {
		ctx.Input.CruSession.Set("uname", uname)
		return true
	}
	return false
}
