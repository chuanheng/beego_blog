package controllers

import (
	"beego_blog/models"

	"github.com/astaxie/beego"
)

type ReplyController struct {
	beego.Controller
}

func (c *ReplyController) Add() {
	tid := c.Input().Get("tid")
	err := models.AddReply(tid,
		c.Input().Get("nickname"),
		c.Input().Get("content"))
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view/"+tid, 301)
}

func (c *ReplyController) Del() {
	rid := c.Input().Get("rid")
	tid := c.Input().Get("tid")
	err := models.DelReply(tid, rid)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view/"+tid, 301)
}
