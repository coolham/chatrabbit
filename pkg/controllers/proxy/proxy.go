package proxy

import (
	"chatrabbit/api/response"
	"chatrabbit/config"
	"chatrabbit/pkg/infra/log"
	"chatrabbit/pkg/services/proxyserv"
	"fmt"
	"io"
	"net"
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
	Service  proxyserv.ProxyService
	Sessions *sessions.Sessions
}

// IsFilterHeader checks if a header should be filtered out
func IsFilterHeader(key string) bool {
	switch strings.ToLower(key) {
	case "x-real-ip", "x-forwarded-for", "user-agent", "referer", "cookie":
		return true
	}
	return false
}

// 获取服务器的IP地址
func getServerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Errorf("failed to get server IP, %v", err)
		return "unknown"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown"
}

// 获取目标域名和协议
func getTargetDomainAndScheme(configServ config.Config, host string) (string, string, error) {
    log.Infof("prepare target domain, host=%s", host)


	domainMappingsInterface := configServ.GetStringMap("proxy.domain_mappings")
	domainMappings := make(map[string]string)

	for key, value := range domainMappingsInterface {
		if strValue, ok := value.(string); ok {
			domainMappings[key] = strValue
			log.Debugf("map key=%s, v=%s", key, strValue)
		} else {
			log.Errorf("invalid domain mapping for key: %s", key)
			return "", "", fmt.Errorf("invalid domain mapping for key: %s", key)
		}
	}

	targetDomain, exists := domainMappings[host]
	if !exists {
		log.Errorf("no mapping found for domain: %s", host)
		return "", "", fmt.Errorf("no mapping found for domain: %s", host)
	}

	// 使用配置中的默认协议
	targetScheme := configServ.GetString("proxy.default_scheme")
	if targetScheme == "" {
		targetScheme = "https" // 默认将 HTTP 转换为 HTTPS
	}

	return targetDomain, targetScheme, nil
}

// handleRequest processes the proxy request
func (c *ProxyController) handleRequest(method string) mvc.Result {
	configServ, err := config.GetConfig()
	if err != nil {
		log.Errorf("get config serve error, %v", err)
		return response.ErrCodeResp(err)
	}

	// 获取请求URL
	reqUrl := c.Ctx.FullRequestURI()
	log.Infof("new proxy %s request, %s", method, reqUrl)

	// 解析请求URL
	parsedUrl, err := url.Parse(reqUrl)
	if err != nil {
		log.Errorf("failed to parse request URL, %v", err)
		return response.ErrCodeResp(err)
	}
        log.Infof("parsed url=%s", parsedUrl)

	// 获取目标域名和协议
	targetDomain, targetScheme, err := getTargetDomainAndScheme(configServ, parsedUrl.Host)
	if err != nil {
		return response.ErrCodeResp(err)
	}

	// 替换域名和协议
	parsedUrl.Host = targetDomain
	parsedUrl.Scheme = targetScheme
	newUrl := parsedUrl.String()
	log.Infof("new url, %s", newUrl)

	// 创建新的请求
	req, err := http.NewRequest(method, newUrl, c.Ctx.Request().Body)
	if err != nil {
		log.Errorf("new request error, %v", err)
		return response.ErrCodeResp(err)
	}

	// 传递Header
	for key, values := range c.Ctx.Request().Header {
		if IsFilterHeader(key) {
			continue
		}
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	// 替换敏感头信息
	serverIP := getServerIP()
	req.Header.Set("X-Forwarded-For", serverIP)
	req.Header.Set("User-Agent", c.Ctx.GetHeader("User-Agent"))

	// 发送请求并获取响应
	startTime := time.Now()
	client := &http.Client{
		Timeout: 30 * time.Second, // 设置超时
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 禁用自动重定向
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("send request error, %v", err)
		return mvc.Response{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	// 获取响应Body
	body, err := io.ReadAll(resp.Body)
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

// Get is the proxy handler
func (c *ProxyController) Get() mvc.Result {
	return c.handleRequest("GET")
}

// Post is the proxy handler
func (c *ProxyController) Post() mvc.Result {
	return c.handleRequest("POST")
}

// Put is the proxy handler
func (c *ProxyController) Put() mvc.Result {
	return c.handleRequest("PUT")
}

// Delete is the proxy handler
func (c *ProxyController) Delete() mvc.Result {
	return c.handleRequest("DELETE")
}

// Patch is the proxy handler
func (c *ProxyController) Patch() mvc.Result {
	return c.handleRequest("PATCH")
}

// Head is the proxy handler
func (c *ProxyController) Head() mvc.Result {
	return c.handleRequest("HEAD")
}

// Options is the proxy handler
func (c *ProxyController) Options() mvc.Result {
	return c.handleRequest("OPTIONS")
}
