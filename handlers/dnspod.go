package handlers

import (
	"ddns/response"
	"ddns/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

func LoginView(c *gin.Context) {
	loginV := utils.LoginView()
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(loginV))
}
func LoginAuth(c *gin.Context) {
	login_email := c.PostForm("login_email")
	login_password := c.PostForm("login_password")
	validUser := checUserPasswd(login_email, login_password)
	if validUser {
		session := sessions.Default(c)
		session.Set("user", login_email)
		session.Save()
		msg := utils.MesTpl("success", "登陆成功", "domainlist")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
		return
	} else {
		msg := utils.MesTpl("fail", "账户或密码错误", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
		return
	}

}

func checUserPasswd(email, password string) bool {
	validEmail := "ld@omyue.com"
	validPassword := "123456"
	return email == validEmail && password == validPassword
}

// 域名列表
func DomainList(c *gin.Context) {
	var domainlist response.DomainAll
	client := utils.Api()
	request := dnspod.NewDescribeDomainListRequest()
	response, _ := client.DescribeDomainList(request)

	json.Unmarshal([]byte(response.ToJsonString()), &domainlist)
	domain_sub := utils.DomainSub(domainlist.Response.DomainList)
	content := utils.GetTpl(utils.Domain)
	content = strings.ReplaceAll(content, "{{.title}}", "域名列表")
	content = strings.ReplaceAll(content, "{{.list}}", domain_sub)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

// 记录列表
func Recordlist(c *gin.Context) {
	domain := c.Query("domain")
	action := c.Query("action")
	grade := c.Query("grade")
	var recordlist response.RecordAll
	client := utils.Api()
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = common.StringPtr(domain)
	request.Limit = common.Uint64Ptr(1000)
	response, _ := client.DescribeRecordList(request)
	json.Unmarshal([]byte(response.ToJsonString()), &recordlist)
	record_sub := utils.RecordSub(recordlist.Response.RecordList, domain, grade)
	content := utils.GetTpl(utils.Record)
	content = strings.ReplaceAll(content, "{{.domain}}", domain)
	content = strings.ReplaceAll(content, "{{.action}}", action)
	content = strings.ReplaceAll(content, "{{.grade}}", grade)
	content = strings.ReplaceAll(content, "{{.title}}", "记录列表")
	content = strings.ReplaceAll(content, "{{.list}}", record_sub)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

// 修改记录页面
func RecordEditf(c *gin.Context) {
	domain := c.Query("domain")
	record_id := c.Query("record_id")
	grade := c.Query("grade")
	action := c.Query("action")
	var recordInfo response.RecordEdit
	record_iduint, _ := strconv.ParseUint(record_id, 10, 64)
	client := utils.Api()
	request := dnspod.NewDescribeRecordRequest()
	request.Domain = common.StringPtr(domain)
	request.RecordId = common.Uint64Ptr(record_iduint)
	response, _ := client.DescribeRecord(request)
	gradeList, lineList := GetGrade(grade, domain)
	json.Unmarshal([]byte(response.ToJsonString()), &recordInfo)
	recordInfo.Response.RecordEditInfo.Domain = domain
	recordE := utils.RecordEdit(recordInfo.Response.RecordEditInfo, gradeList, lineList, action)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(recordE))
}

// 添加域名
func Domaincreate(c *gin.Context) {
	domain := c.PostForm("domain")

	client := utils.Api()
	request := dnspod.NewCreateDomainRequest()

	request.Domain = common.StringPtr(domain)
	_, err := client.CreateDomain(request)
	if err != nil {
		msg := utils.MesTpl("fail", err.(*errors.TencentCloudSDKError).GetMessage(), "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
		return
	}
	msg := utils.MesTpl("success", "添加"+domain+"成功", "")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
}

// 修改或增加记录
func ModifyRecord(c *gin.Context) {
	domain := c.Query("domain")
	action := c.Query("action")
	record_id := c.Query("record_id")
	sub_domain := c.PostForm("sub_domain")
	gradetype := c.PostForm("type")
	line := c.PostForm("line")
	value := c.PostForm("value")
	mx := c.PostForm("mx")
	ttl := c.PostForm("ttl")
	mx_int, _ := strconv.ParseUint(mx, 10, 64)
	ttl_int, _ := strconv.ParseUint(ttl, 10, 64)
	record_id_int, _ := strconv.ParseUint(record_id, 10, 64)
	client := utils.Api()
	if action == "edit" {
		request := dnspod.NewModifyRecordRequest()
		request.Domain = common.StringPtr(domain)
		request.SubDomain = common.StringPtr(sub_domain)
		request.RecordType = common.StringPtr(gradetype)
		request.RecordLine = common.StringPtr(line)
		request.Value = common.StringPtr(value)
		request.MX = common.Uint64Ptr(mx_int)
		request.TTL = common.Uint64Ptr(ttl_int)
		request.RecordId = common.Uint64Ptr(record_id_int)

		_, err := client.ModifyRecord(request)
		if err != nil {
			msg := utils.MesTpl("fail", err.(*errors.TencentCloudSDKError).GetMessage(), "")
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
			return
		}
		msg := utils.MesTpl("success", "修改"+sub_domain+"."+domain+"成功", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
		return
	} else if action == "add" {
		request := dnspod.NewCreateRecordRequest()
		request.Domain = common.StringPtr(domain)
		request.SubDomain = common.StringPtr(sub_domain)
		request.RecordType = common.StringPtr(gradetype)
		request.RecordLine = common.StringPtr(line)
		request.Value = common.StringPtr(value)
		request.MX = common.Uint64Ptr(mx_int)
		request.TTL = common.Uint64Ptr(ttl_int)
		_, err := client.CreateRecord(request)
		if err != nil {
			msg := utils.MesTpl("fail", err.(*errors.TencentCloudSDKError).GetMessage(), "")
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
			return
		}
		msg := utils.MesTpl("success", "添加"+sub_domain+"."+domain+"成功", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
		return
	}

}

// 删除域名
func Domainremove(c *gin.Context) {
	domain := c.Query("domain")
	client := utils.Api()
	request := dnspod.NewDeleteDomainRequest()

	request.Domain = common.StringPtr(domain)
	_, err := client.DeleteDomain(request)
	if err != nil {
		msg := utils.MesTpl("fail", err.(*errors.TencentCloudSDKError).GetMessage(), "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
		return
	}
	msg := utils.MesTpl("success", "删除"+domain+"成功", "")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
}

// 域名状态
func DomainStatus(c *gin.Context) {
	domain := c.Query("domain")
	status := c.Query("status")
	client := utils.Api()
	request := dnspod.NewModifyDomainUnlockRequest()

	request.Domain = common.StringPtr(domain)
	request.LockCode = common.StringPtr(status)
	_, err := client.ModifyDomainUnlock(request)
	if err != nil {
		msg := utils.MesTpl("fail", err.(*errors.TencentCloudSDKError).GetMessage(), "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
		return
	}
	msg := utils.MesTpl("success", "删除"+domain+"成功", "")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
}

// 记录状态
func RecordStatus(c *gin.Context) {
	domain := c.Query("domain")
	record_id := c.Query("record_id")
	status := c.Query("status")
	record_id_int, _ := strconv.ParseUint(record_id, 10, 64)
	client := utils.Api()
	request := dnspod.NewModifyRecordStatusRequest()
	request.Domain = common.StringPtr(domain)
	request.RecordId = common.Uint64Ptr(record_id_int)
	request.Status = common.StringPtr(status)
	_, err := client.ModifyRecordStatus(request)
	if err != nil {
		msg := utils.MesTpl("fail", err.(*errors.TencentCloudSDKError).GetMessage(), "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
		return
	}
	msg := utils.MesTpl("success", "修改"+domain+"状态成功", "")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
}

// 记录删除
func RecordRemove(c *gin.Context) {
	domain := c.Query("domain")
	record_id := c.Query("record_id")
	record_id_int, _ := strconv.ParseUint(record_id, 10, 64)
	client := utils.Api()
	request := dnspod.NewDeleteRecordRequest()
	request.Domain = common.StringPtr(domain)
	request.RecordId = common.Uint64Ptr(record_id_int)

	_, err := client.DeleteRecord(request)
	if err != nil {
		msg := utils.MesTpl("fail", err.(*errors.TencentCloudSDKError).GetMessage(), "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
		return
	}
	msg := utils.MesTpl("success", "删除"+domain+"记录成功", "")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg))
}

// 获取DomainGrade类型
func GetGrade(grade, domain string) (a, b []string) {
	var gradeResponse response.GradeResponse
	var lineListResponse response.LineListResponse
	client := utils.Api()
	request := dnspod.NewDescribeRecordTypeRequest()
	request.DomainGrade = common.StringPtr(grade)
	response, _ := client.DescribeRecordType(request)
	json.Unmarshal([]byte(response.ToJsonString()), &gradeResponse)

	request1 := dnspod.NewDescribeRecordLineListRequest()
	request1.Domain = common.StringPtr(domain)
	request1.DomainGrade = common.StringPtr(grade)
	response1, _ := client.DescribeRecordLineList(request1)
	json.Unmarshal([]byte(response1.ToJsonString()), &lineListResponse)
	var lineList []string
	for _, line := range lineListResponse.Response.LineList {
		lineList = append(lineList, line.Name)
	}
	return gradeResponse.Response.TypeList, lineList
}
