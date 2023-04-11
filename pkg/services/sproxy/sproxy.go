package sproxy

type ProxyService interface {
}

type proxyService struct {
}

func NewProxyService() ProxyService {
	return &proxyService{}
}
