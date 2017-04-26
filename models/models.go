package models

import (
	"time"

	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//数据库连对象需要的信息
var (
	dburl  string = "tcp(192.168.113.3:3306)"
	dbuser string = "lpd_crowd"
	dbpwd  string = "lpdE@Mlcrowd"
	dbname string = "test"
)

//初始化orm
func RegisterDB() {
	conn := dbuser + ":" + dbpwd + "@" + dburl + "/" + dbname + "?charset=utf8&loc=Local" //组合成连接串
	orm.RegisterModel(new(User), new(Category), new(Topic), new(Reply))                   // 注册 Category 如果没有会自动创建
	orm.RegisterDriver("mysql", orm.DRMySQL)                                              //注册mysql驱动
	orm.RegisterDataBase("default", "mysql", conn)                                        //设置conn中的数据库为默认使用数据库
	orm.RunSyncdb("default", false, false)                                                //后一个使用true会带上很多打印信息，数据库操作和建表操作的；第二个为true代表强制创建表
}

func AddCategory(name string) error {
	o := orm.NewOrm()
	cate := &Category{Title: name, Created: time.Now(), Views: 0, TopicTime: time.Now()}

	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}
	_, err1 := o.Insert(cate)
	if err1 != nil {
		return err1
	}
	return nil
}

func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

func GetCategory(cate string) (*Category, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("category")

	category := new(Category)
	err := qs.Filter("title", cate).One(category)

	return category, err
}

func DelCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	return err
}

func AddTopic(category, title, label, content, attachment string) error {

	// 处理标签
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	o := orm.NewOrm()
	topic := &Topic{
		Title:      title,
		Content:    content,
		Category:   category,
		Labels:     label,
		Created:    time.Now(),
		Updated:    time.Now(),
		ReplyTime:  time.Now(),
		Attachment: attachment,
	}

	_, err1 := o.Insert(topic)
	if err1 != nil {
		return err1
	}

	// 分类文章数 + 1
	cate, err1 := GetCategory(category)
	if cate != nil {
		cate.TopicCount++
		o.Update(cate)
	}

	return nil
}

func ModifyTopic(id, category, title, label, content, attachment string) error {

	// 处理标签
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	topic := &Topic{Id: cid}

	topic, err = GetTopic(id, false)

	// 如果文章分类改变了
	if !strings.EqualFold(topic.Category, category) {
		// 分类文章数 + 1
		cate, err1 := GetCategory(category)
		if err1 != nil {
			beego.Error(err1)
		}
		if cate != nil {
			cate.TopicCount++
			o.Update(cate)
		}

		// 修改前的分类数 - 1
		cate_before, err2 := GetCategory(topic.Category)
		if err2 != nil {
			beego.Error(err2)
		}
		fmt.Println("cate_before : ", cate_before)
		if cate_before != nil {
			cate_before.TopicCount--
			o.Update(cate_before)
		}
	}

	if o.Read(topic) == nil {
		topic.Title = title
		topic.Content = content
		topic.Category = category
		topic.Labels = label
		topic.Updated = time.Now()
		topic.Attachment = attachment
		_, err2 := o.Update(topic)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

func GetAllTopics(cate, label string, isDesc bool) ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	var err error
	if isDesc {
		if len(cate) > 0 {
			_, err = qs.Filter("category", cate).OrderBy("-created").All(&topics)
		} else if len(label) > 0 {
			_, err = qs.Filter("labels__icontains", label).OrderBy("-created").All(&topics)
		} else {
			_, err = qs.OrderBy("-created").All(&topics)

		}
	} else {
		_, err = qs.All(&topics)
	}
	return topics, err
}

// 根据 ID 获得文章
func GetTopic(id string, count bool) (*Topic, error) {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	qs := o.QueryTable("topic")

	topic := new(Topic)
	err = qs.Filter("id", cid).One(topic)

	// 是否增加浏览数
	if count {
		topic.Views++
		_, err = o.Update(topic)
	}

	topic.Labels = strings.Replace(strings.Replace(topic.Labels, "#", " ", -1), "$", "", -1)

	return topic, err
}

// 删除文章
func DelTopic(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	topic := &Topic{Id: cid}
	_, err = o.Delete(topic)

	topic, err = GetTopic(id, false)

	// 分类文章数 - 1
	cate, err1 := GetCategory(topic.Category)
	if err1 != nil {
		beego.Error(err1)
	}
	if cate != nil {
		cate.TopicCount--
		o.Update(cate)
	}

	return err
}

// 添加评论
func AddReply(tid, nickname, content string) error {
	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	reply := &Reply{
		Tid:     id,
		Content: content,
		Name:    nickname,
		Created: time.Now(),
	}

	_, err1 := o.Insert(reply)
	if err1 != nil {
		return err1
	}

	// 添加评论数
	topic, err := GetTopic(tid, false)
	if topic != nil {
		topic.ReplyCount++
		topic.Updated = time.Now()
		topic.ReplyTime = time.Now()
		_, err1 = o.Update(topic)
		if err1 != nil {
			return err1
		}
	}

	return nil
}

// 删除评论
func DelReply(tid, id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	reply := &Reply{Id: cid}
	_, err = o.Delete(reply)

	// 减少评论数
	topic, err := GetTopic(tid, false)
	if topic != nil {
		topic.ReplyCount--
		o.Update(topic)
	}
	return err
}

// 获得所有评论
func GetAllReplies(tid string) ([]*Reply, error) {
	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	replies := make([]*Reply, 0)
	qs := o.QueryTable("reply")
	_, err = qs.Filter("tid", id).OrderBy("-created").All(&replies)
	return replies, err
}

func AddUser(user_name, passpord string) error {

	o := orm.NewOrm()
	user := &User{
		UserName: user_name,
		PassWord: passpord,
		Created:  time.Now(),
	}

	_, err1 := o.Insert(user)
	if err1 != nil {
		return err1
	}

	return nil
}

func GetUser(user_name string) (*User, error) {

	o := orm.NewOrm()
	qs := o.QueryTable("user")
	user := new(User)
	err := qs.Filter("user_name", user_name).One(user)
	return user, err
}

type User struct {
	Id       int64
	UserName string
	PassWord string
	Created  time.Time `orm:"index"`
}

type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Content         string `orm:"size(5000)"`
	Category        string
	Labels          string
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time
	ReplyCount      int64
	ReplyLastUserId int64
}

type Reply struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}
