package controllers

import (
	"beego_blog/models"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	opt := c.Input().Get("opt")
	// 判断操作类型
	switch opt {
	// 如果是修改文章
	case "modify":
		id := c.Input().Get("id")
		if len(id) == 0 {
			break
		}
		topic, err := models.GetTopic(id, false)
		if err != nil {
			beego.Error(err)
		}
		c.Data["Topic"] = topic
		c.Data["IsModify"] = true
		c.Data["IsTopic"] = true
		c.TplName = "topic_add.html"
		return
	// 如果是删除文章
	case "del":
		id := c.Input().Get("id")
		if len(id) == 0 {
			break
		}
		err := models.DelTopic(id)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/topic", 301)
	}

	// 获取所有文章
	var err error
	c.Data["Topics"], err = models.GetAllTopics("", "", false)
	if err != nil {
		beego.Error(err)
	}

	c.Data["IsTopic"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.TplName = "topic.html"
}

// 添加或者修改文章
func (c *TopicController) Post() {
	// 判断是否是博客管理员
	if checkAccount(c.Ctx) {

		// 解析表单
		opt := c.Input().Get("opt")
		tid := c.Input().Get("id")
		category := c.Input().Get("category")
		title := c.Input().Get("title")
		label := c.Input().Get("label")
		content := c.Input().Get("content")

		// 获取附件
		_, fh, err := c.GetFile("attachment")
		if err != nil {
			beego.Error(err)
		}

		var attachment string
		if fh != nil {
			// 保存附件
			attachment = fh.Filename
			beego.Info(attachment)
			err = c.SaveToFile("attachment", path.Join("attachment", attachment))
			// filename: tmp.go
			// attachment/tem.go
			if err != nil {
				beego.Error(err)
			}
		}

		if opt != "" && opt == "modify" {
			err := models.ModifyTopic(
				tid,
				category,
				title,
				label,
				content,
				attachment)
			if err != nil {
				beego.Error(err)
			}
		} else {
			err := models.AddTopic(
				category,
				title,
				label,
				content,
				attachment)
			if err != nil {
				beego.Error(err)
			}
		}
	}
	c.Redirect("/topic", 301)
}

// 跳转到发送消息界面
func (c *TopicController) Add() {
	c.Data["IsTopic"] = true
	c.TplName = "topic_add.html"
}

// 浏览博客
func (c *TopicController) View() {
	tid := c.Ctx.Input.Params()["0"]
	topic, err := models.GetTopic(tid, true)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	replies, err1 := models.GetAllReplies(tid)
	if err1 != nil {
		beego.Error(err1)
	}

	c.Data["Replies"] = replies
	c.Data["Topic"] = topic
	c.Data["Labels"] = strings.Fields(topic.Labels)
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.TplName = "topic_view.html"
}
