{{define "navbar"}}
<a class="navbar-brand" href="/">我的 beego 博客</a>
<div>
    <ul class="nav navbar-nav">
        <li {{if .IsHome}}class="active"{{end}}><a href="/">首页</a></li>
        <li {{if .IsCategory}}class="active"{{end}}><a href="/category">分类</a></li>
        <li {{if .IsTopic}}class="active"{{end}}><a href="/topic">文章</a></li>
    </ul>
</div>

<div class="pull-right">
    <ul class="nav navbar-nav">
        {{if .IsLogin}}
            <li><a href="/login?exist=true">退出</a></li>
        {{else}}
            <li><a href="/login">登录</a></li>
			<li><a href="/register">注册</a></li>
        {{end}}
    </ul>
</div>
{{end}}