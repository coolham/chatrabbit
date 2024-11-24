package proxy

import (
	"chatrabbit/config"
	"chatrabbit/pkg/infra/log"
	"fmt"
	"net"
	"strings"

	"github.com/kataras/iris/v12"
)

// IsFilterHeader checks if a header should be filtered out
func IsFilterHeader(key string) bool {
	switch strings.ToLower(key) {
	case "x-real-ip", "x-forwarded-for", "user-agent", "referer", "cookie":
		return true
	}
	return false
}

// 验证并提取 Authorization 头部
func extractAuthorizationToken(ctx iris.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		log.Errorf("Authorization header is missing")
		return "", fmt.Errorf("authorization token is required")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		log.Errorf("Authorization token is empty")
		return "", fmt.Errorf("authorization token is empty")
	}
	log.Debugf("Authorization token: %s", token)
	return token, nil
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
func getTargetDomainAndScheme(configServ config.Config, host string, scheme string) (string, string, error) {
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
