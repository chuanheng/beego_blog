package controllers

import (
	"beego_blog/models"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

type BaseController struct {
	beego.Controller
	i18n.Locale
}

func (c *BaseController) Prepare() {
	lang := c.GetString("lang")
	if lang == "zh-CN" {
		c.Lang = lang
	} else {
		c.Lang = "en-US"
	}
	c.Data["Lang"] = lang
}

type MainController struct {
	BaseController
}

func (c *MainController) Get() {

	c.Data["hi"] = c.Tr("hi")
	c.Data["bye"] = c.Tr("bye")

	c.Data["IsHome"] = true
	var err error
	c.Data["Topics"], err = models.GetAllTopics(
		c.Input().Get("cate"),
		c.Input().Get("label"),
		true)
	if err != nil {
		beego.Error(err)
	}
	c.TplName = "home.html"

	c.Data["Categories"], err = models.GetAllCategories()
	c.Data["IsLogin"] = checkAccount(c.Ctx)
}
