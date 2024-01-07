package proxyserv

import "chatrabbit/pkg/infra/log"

type ProxyService interface {
	Start() error
	Stop() error
}

type proxyService struct {
}

func NewProxyService() ProxyService {
	return &proxyService{}
}

func (s *proxyService) Start() error {
	log.Info("proxy service start")
	return nil
}

func (s *proxyService) Stop() error {
	log.Info("proxy service stop")
	return nil
}
