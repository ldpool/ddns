package utils

import (
	"ddns/response"
	"fmt"
	"strconv"
	"strings"
)

func GetTpl(tpl string) string {
	master := Index
	master = strings.ReplaceAll(master, "{{.content}}", tpl)
	return master
}

func MesTpl(status, msg, url string) string {

	msgtext := GetTpl(Message)
	var title string
	if status == "success" {
		title = "操作成功"
	} else {
		title = "操作失败"
	}
	if url == "" {
		url = "-1"
	}
	msgtext = strings.ReplaceAll(msgtext, "{{.title}}", title)
	msgtext = strings.ReplaceAll(msgtext, "{{.status}}", status)
	msgtext = strings.ReplaceAll(msgtext, "{{.message}}", msg)
	msgtext = strings.ReplaceAll(msgtext, "{{.url}}", url)
	return msgtext
}

func DomainSub(domainList []response.DomainInfo) string {
	var result string
	for _, v := range domainList {
		statusText := "启用"
		statusAction := "ENABLE"
		if v.Status == "ENABLE" {
			statusText = "暂停"
			statusAction = "PAUSE"
		} else if v.Status == "PAUSE" {
			statusText = "启用"
			statusAction = "ENABLE"
		}
		domain := v.MapStatus()
		result += fmt.Sprintf(`<tr>
                <td>%d</td><td><a href="recordlist?domain=%s&grade=%s">%s</a></td><td>%s</td><td>%s</td><td>%d</td><td>%s</td><td><a href="recordlist?domain=%s&grade=%s">记录</a> <a href="domainstatus?domain=%s&status=%s">%s</a> <a href="domainremove?domain=%s">删除</a></td>
            </tr>`,
			domain.DomainId, domain.Name, domain.Grade, domain.Name, domain.GradeTitle, domain.Status, domain.RecordCount, domain.UpdatedOn, domain.Name, domain.Grade, domain.Name, statusAction, statusText, domain.Name,
		)
	}
	return result
}

func RecordSub(recordList []response.RecordInfo, domain, grade string) string {
	var result string
	for _, v := range recordList {
		statusText := "启用"
		statusAction := "ENABLE"

		if v.Status == "ENABLE" {
			statusText = "暂停"
			statusAction = "DISABLE"
		} else if v.Status == "DISABLE" {
			statusText = "启用"
			statusAction = "ENABLE"
		}
		record := v.MapStatus()
		result += fmt.Sprintf(`<tr>
		<td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td> <a href="recordeditf?action=edit&domain=%s&record_id=%d&grade=%s">修改</a> <a href="recordstatus?domain=%s&record_id=%d&status=%s">%s</a> <a href="recordremove?domain=%s&record_id=%d">删除</a></td>
	</tr>`,
			record.RecordId, record.Name, record.Type, record.Line, record.Value, record.Status, record.MX, record.TTL, domain, record.RecordId, grade, domain, record.RecordId, statusAction, statusText, domain, record.RecordId,
		)
	}
	return result
}

func RecordEdit(record response.RecordEditInfo, gradeList, lineList []string, action string) string {
	grade := optionList(gradeList, record.RecordType)
	line := optionList(lineList, record.RecordLine)
	recordcreatef := GetTpl(Recordcreatef)
	if action == "edit" {
		recordcreatef = strings.ReplaceAll(recordcreatef, "{{.title}}", "修改记录")
		recordcreatef = strings.ReplaceAll(recordcreatef, "{{.ttl}}", strconv.Itoa(record.TTL))
	} else if action == "add" {
		recordcreatef = strings.ReplaceAll(recordcreatef, "{{.title}}", "增加记录")
		recordcreatef = strings.ReplaceAll(recordcreatef, "{{.ttl}}", "600")
	}

	recordcreatef = strings.ReplaceAll(recordcreatef, "{{.domain}}", record.Domain)
	recordcreatef = strings.ReplaceAll(recordcreatef, "{{.action}}", action)
	recordcreatef = strings.ReplaceAll(recordcreatef, "{{.domain_id}}", strconv.Itoa(record.DomainId))
	recordcreatef = strings.ReplaceAll(recordcreatef, "{{.record_id}}", strconv.Itoa(record.Id))
	recordcreatef = strings.ReplaceAll(recordcreatef, "{{.sub_domain}}", record.SubDomain)
	recordcreatef = strings.ReplaceAll(recordcreatef, "{{.type_list}}", grade)
	recordcreatef = strings.ReplaceAll(recordcreatef, "{{.line_list}}", line)
	recordcreatef = strings.ReplaceAll(recordcreatef, "{{.value}}", record.Value)
	recordcreatef = strings.ReplaceAll(recordcreatef, "{{.mx}}", strconv.Itoa(record.MX))

	return recordcreatef
}

func LoginView() string {

	login := GetTpl(Login)
	login = strings.ReplaceAll(login, "{{.title}}", "")
	return login
}

func optionList(types []string, fist string) string {
	if fist != "" {
		types = append([]string{fist}, types...)
	}
	var typeList string
	for _, value := range types {
		typeList += fmt.Sprintf(`<option value="%s">%s</option>`, value, value)
	}
	return typeList
}

var (
	Index = `
	<!DOCTYPE html>
	<html lang="zh-cn">
		<head>
			<meta charset="utf-8">
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<!-- <title>{{.title}} - DNSPod</title> -->
			<link href="http://libs.baidu.com/bootstrap/3.0.3/css/bootstrap.min.css" rel="stylesheet" />
		</head>
		<body style=" margin: 0;width:40%; margin-left:30%;">
			<div class="dnspod">
				
				<h3>{{.title}}</h3>
				<div>{{.content}}</div>
			 
			</div>
		</body>
	</html>
`
	Domain = `
<div class="add_domain">
<form class="form-inline"  method="POST" action="domaincreate">
    <div class="form-group">
        <input type="text" name="domain" class="form-control" style="width:60%" placeholder="如: ddv.cn" />
        <input type="submit" class="btn btn-default" value="添加" />
    </div>
</form>
</div>
<table class="table hovered table-striped">
    <thead>
        <tr>
            <th>编号</th><th>域名</th><th>等级</th><th>状态</th><th>记录</th><th>更新</th><th>操作</th>
        </tr>
    </thead>
    <tbody>
        {{.list}}
    </tbody>
</table>
`

	Login = `
	<form method="POST" action="loginauth" style=" margin: 0;width:40%; margin-left:30%;margin-top: 60px;">
	<div class="input-group input-group-lg t">
		<span class="input-group-addon">邮箱</span>
		<input type="text" name="login_email" class="form-control">
	</div>
	<div class="input-group input-group-lg t">
		<span class="input-group-addon">密码</span>
		<input type="password" name="login_password" class="form-control">
	</div>
	<input type="submit" class="btn btn-primary btn-lg btn-block" value="提交" />
	</form>
`
	Logind = `
	<form method="POST" action="?action={{.action}}">
	<div class="input-group input-group-lg t">
		<span class="input-group-addon">D令牌</span>
		<input type="text" name="login_code" class="form-control">
	</div>
	<input type="submit" class="btn btn-primary btn-lg btn-block" value="提交" />
	</form>
`
	Message = `
	<div class="result_message bg-{{.status}}">
    <div id="message">{{.message}}</div>
    <div id="back">正在返回，请稍候...</div>
</div>
<script type="text/javascript">var url = '{{.url}}';if (url) {setTimeout(function() {if (url == '-1') {history.back();} else {location.href=url;}}, 3000);} else {document.getElementById('back').style.display='none';}</script>
`
	Record = `
	<div>
    <input type="button" class="btn btn-default" value="返回域名" onclick="location.href='/domainlist'" />
    <input type="button" class="btn btn-default" value="添加记录" onclick="location.href='/recordeditf?action=add&domain={{.domain}}&grade={{.grade}}'" />
</div>
<table class="table hovered table-striped">
    <thead>
        <tr>
            <th>编号</th><th>子域名</th><th>类型</th><th>线路</th><th>记录</th><th>状态</th><th>MX</th><th>TTL</th><th>操作</th>
        </tr>
    </thead>
    <tbody>
        {{.list}}
    </tbody>
</table>
`
	Recordcreatef = `
	<div class="add_record">
	<form method="POST" action="modifyrecord?action={{.action}}&domain={{.domain}}&domain_id={{.domain_id}}&record_id={{.record_id}}">
		<div class="form-group">
			<label for="sub_domain">主机名</label>
			<input type="text" class="form-control" id="sub_domain" name="sub_domain" placeholder="如:www" value="{{.sub_domain}}">
		</div>
		<div class="form-group">
			<label for="type">类型</label>
			<select class="form-control" id="type" name="type">
				{{.type_list}}
			</select>
		</div>
		<div class="form-group">
			<label for="line">线路</label>
			<select class="form-control" id="line" name="line">
				{{.line_list}}
			</select>
		</div>
		<div class="form-group">
			<label for="value">记录值</label>
			<input type="text" class="form-control" id="value" name="value" placeholder="如:1.1.1.1" value="{{.value}}">
		</div>
		<div class="form-group">
			<label for="mx">MX</label>
			<input type="text" class="form-control" id="mx" name="mx" placeholder="如:60" value="{{.mx}}">
		</div>
		<div class="form-group">
			<label for="ttl">TTL</label>
			<input type="text" class="form-control" id="ttl" name="ttl" placeholder="如:600" value="{{.ttl}}">
		</div>
		<input type="submit" class="btn btn-primary btn-lg btn-block" value="提交" />
	</form>
	</div>
`
)
