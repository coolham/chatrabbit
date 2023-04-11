package cproxy

import (
	"bytes"
	"chatrabbit/config"
	"chatrabbit/config/common"
	"chatrabbit/pkg/infra/log"
	"chatrabbit/pkg/services/sproxy"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type ProxyController struct {
	Ctx      iris.Context
	Service  sproxy.ProxyService
	Sessions *sessions.Sessions
}

func IsFilterHeader(key string) bool {

	switch key {
	case "X-Real-Ip":
		return true
	case "X-Forwarded-For":
		return true
	}
	return false
}

func (c *ProxyController) Get() mvc.Result {
	// 获取请求URL
	reqUrl := c.Ctx.FullRequestURI()
	log.Infof("new proxy %s request, %s", c.Ctx.Method(), reqUrl)

	queryString := c.Ctx.Request().URL.Query()

	path := c.Ctx.Path()
	log.Infof("req path: %s", path)

	newPath := strings.Replace(path, "/proxy", "", 1)
	log.Infof("new path:%s", newPath)

	// 将GET参数转换为url参数字符串
	var arr []string
	for k, v := range queryString {
		for _, value := range v {
			arr = append(arr, fmt.Sprintf("%v=%v", k, value))
		}
	}
	newQuery := strings.Join(arr, "&")

	// 拼接url
	proxyUrl := config.GetString(common.PROXY_URL) + newPath
	baseUrl, _ := url.Parse(proxyUrl)
	baseUrl.RawQuery = newQuery
	newUrl := baseUrl.String()
	log.Infof("new url, %s", newUrl)

	// 创建新的请求
	req, err := http.NewRequest(c.Ctx.Method(), newUrl, c.Ctx.Request().Body)
	if err != nil {
		log.Errorf("new request error, %v", err)
		return mvc.Response{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	// 传递Header
	for key, values := range c.Ctx.Request().Header {
		if IsFilterHeader(key) {
			continue
		}
		for _, value := range values {
			// log.Infof("request header, %s=%s", key, value)
			req.Header.Set(key, value)
		}
	}
	// 发送请求并获取响应
	startTime := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("send request error, %v", err)
		return mvc.Response{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	// 获取响应Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read response body error, %v", err)
		return mvc.Response{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}
	defer resp.Body.Close()
	finishedTime := time.Now()
	elapseSeconds := finishedTime.Sub(startTime).Seconds()

	// log.Infof("proxy response, code=%d, body=%s", resp.StatusCode, string(body))
	log.Infof("proxy response, code=%d, elapse=%.1f, url=%s", resp.StatusCode, elapseSeconds, newUrl)

	// 设置响应Header
	for key, values := range resp.Header {
		for _, value := range values {
			// log.Infof("set response header, %s=%s", key, value)
			c.Ctx.Header(key, value)
		}
	}

	// 设置响应状态
	c.Ctx.StatusCode(resp.StatusCode)

	// 输出响应Body
	return mvc.Response{
		Code:        resp.StatusCode,
		ContentType: "application/json",
		Content:     body,
	}
}

func (c *ProxyController) Post() mvc.Result {
	// 获取请求URL
	reqUrl := c.Ctx.FullRequestURI()
	log.Infof("new proxy %s request, %s", c.Ctx.Method(), reqUrl)

	queryString := c.Ctx.Request().URL.Query()

	path := c.Ctx.Path()
	log.Infof("req path: %s", path)

	newPath := strings.Replace(path, "/proxy", "", 1)
	log.Infof("new path:%s", newPath)

	// 将GET参数转换为url参数字符串
	var arr []string
	for k, v := range queryString {
		for _, value := range v {
			arr = append(arr, fmt.Sprintf("%v=%v", k, value))
		}
	}
	newQuery := strings.Join(arr, "&")
	fmt.Println("new query:", newQuery)

	// 拼接url
	proxyUrl := config.GetString(common.PROXY_URL) + newPath
	baseUrl, _ := url.Parse(proxyUrl)
	baseUrl.RawQuery = newQuery
	newUrl := baseUrl.String()
	log.Infof("new url, %s", newUrl)

	// 创建新的请求
	req, err := http.NewRequest(c.Ctx.Method(), newUrl, c.Ctx.Request().Body)
	if err != nil {
		log.Errorf("new request error, %v", err)
		return mvc.Response{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	// 传递Header
	for key, values := range c.Ctx.Request().Header {
		if IsFilterHeader(key) {
			continue
		}
		for _, value := range values {
			// log.Infof("request header, %s=%s", key, value)
			req.Header.Set(key, value)
		}
	}
	// 处理POST请求
	reqBody, err := ioutil.ReadAll(c.Ctx.Request().Body)
	if err != nil {
		log.Errorf("read request body error, %v", err)
		return mvc.Response{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	// 发送请求并获取响应
	startTime := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("send request error, %v", err)
		return mvc.Response{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	// 获取响应Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read response body error, %v", err)
		return mvc.Response{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}
	defer resp.Body.Close()
	finishedTime := time.Now()
	elapseSeconds := finishedTime.Sub(startTime).Seconds()

	// log.Infof("proxy response, code=%d, body=%s", resp.StatusCode, string(body))
	log.Infof("proxy response, code=%d, elapse=%.1f, url=%s", resp.StatusCode, elapseSeconds, newUrl)

	// 设置响应Header
	for key, values := range resp.Header {
		for _, value := range values {
			c.Ctx.Header(key, value)
		}
	}

	// 设置响应状态
	c.Ctx.StatusCode(resp.StatusCode)

	// 输出响应Body
	return mvc.Response{
		Code:        resp.StatusCode,
		ContentType: "application/json",
		Content:     body,
	}
}
