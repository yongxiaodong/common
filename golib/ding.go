package golib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"time"
)

func signDing(t int64, secret string) string {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

// 钉钉报警，Text格式
func Ding(token, secret, content string, mobiles []string) error {
	timestamp := time.Now().Unix() * 1000
	sign := signDing(timestamp, secret)
	params := map[string]string{
		"access_token": token,
		"timestamp":    fmt.Sprintf("%d", timestamp),
		"sign":         sign,
	}
	body := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": content,
		},
		"at": map[string]interface{}{
			"atMobiles": mobiles,
			"isAtAll":   false,
		},
	}
	resp, err := resty.New().R().SetQueryParams(params).SetBody(body).Post("https://oapi.dingtalk.com/robot/send")
	if err != nil {
		return errors.Wrap(err, "发送Ding消息失败")
	}
	eStr := gjson.Get(resp.String(), "errmsg").String()
	if eStr != "ok" {
		return errors.Errorf("发送Ding消息失败: %s", eStr)
	}
	return nil
}

// 钉钉报警，Markdown格式
func DingMarkdown(token, secret, content, title string, mobiles []string) error {
	timestamp := time.Now().Unix() * 1000
	sign := signDing(timestamp, secret)
	params := map[string]string{
		"access_token": token,
		"timestamp":    fmt.Sprintf("%d", timestamp),
		"sign":         sign,
	}
	body := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": title,
			"text":  content,
		},
		"at": map[string]interface{}{
			"atMobiles": mobiles,
			"isAtAll":   false,
		},
	}
	resp, err := resty.New().R().SetQueryParams(params).SetBody(body).Post("https://oapi.dingtalk.com/robot/send")
	if err != nil {
		return errors.Wrap(err, "发送Ding消息失败")
	}
	eStr := gjson.Get(resp.String(), "errmsg").String()
	if eStr != "ok" {
		return errors.Errorf("发送Ding消息失败: %s", eStr)
	}
	return nil
}
