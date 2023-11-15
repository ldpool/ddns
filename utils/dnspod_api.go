package utils

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

func Api() (client *dnspod.Client) {
	credential := common.NewCredential(
		"SecretId",
		"SecretKey",
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	c, _ := dnspod.NewClient(credential, "", cpf)

	return c
}
