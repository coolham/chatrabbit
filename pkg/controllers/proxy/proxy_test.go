package proxy

import (
	"chatrabbit/config"
	"chatrabbit/config/common"
	"chatrabbit/pkg/services/proxyserv"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kataras/iris/v12"
	iristest "github.com/kataras/iris/v12/httptest"
	"github.com/kataras/iris/v12/mvc"
	. "github.com/smartystreets/goconvey/convey"
)

func newApp() *iris.Application {
	app := iris.New()
	mvc.Configure(app.Party("/{any:path}"), proxyApp)
	return app
}

func proxyApp(app *mvc.Application) {
	service := proxyserv.NewProxyService()
	ctx := context.Background()
	app.Register(ctx, service)
	app.Handle(new(ProxyController))
}

func TestMain(m *testing.M) {
	// 设置 CONFIGFILE 环境变量
	os.Setenv("RABBIT_CONFIG_FILE", "../../../conf/config.yaml.example")

	// 初始化配置
	err := config.InitConfig("")
	if err != nil {
		fmt.Printf("failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	m.Run()
}

func TestProxyGetModels(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/models" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"data": [{"id": "model-1"}, {"id": "model-2"}]}`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	configServe, _ := config.GetConfig()
	originalProxyURL := configServe.GetString(common.PROXY_URL)
	configServe.SetString(common.PROXY_URL, mockServer.URL)

	app := newApp()
	h := iristest.New(t, app)

	Convey("Test Proxy GET /v1/models", t, func() {
		h.GET("/v1/models").
			Expect().
			Status(iris.StatusOK).
			JSON().Object().Value("data").Array().Length().IsEqual(2)
	})

	configServe.SetString(common.PROXY_URL, originalProxyURL)
}

func TestProxyGetModelDetails(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/models/model-1" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"id": "model-1", "object": "model", "created": 1610078138}`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	configServe, _ := config.GetConfig()
	originalProxyURL := configServe.GetString(common.PROXY_URL)
	configServe.SetString(common.PROXY_URL, mockServer.URL)

	app := newApp()
	h := iristest.New(t, app)

	Convey("Test Proxy GET /v1/models/model-1", t, func() {
		h.GET("/v1/models/model-1").
			Expect().
			Status(iris.StatusOK).
			JSON().Object().Value("id").IsEqual("model-1")
	})

	configServe.SetString(common.PROXY_URL, originalProxyURL)
}

func TestProxyGetFiles(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/files" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"data": [{"id": "file-1"}, {"id": "file-2"}]}`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	configServe, _ := config.GetConfig()
	originalProxyURL := configServe.GetString(common.PROXY_URL)
	configServe.SetString(common.PROXY_URL, mockServer.URL)

	app := newApp()
	h := iristest.New(t, app)

	Convey("Test Proxy GET /v1/files", t, func() {
		h.GET("/v1/files").
			Expect().
			Status(iris.StatusOK).
			JSON().Object().Value("data").Array().Length().IsEqual(2)
	})

	configServe.SetString(common.PROXY_URL, originalProxyURL)
}

func TestProxyPostCreateFile(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/v1/files" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"id": "file-1", "object": "file", "created": 1610078138}`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	configServe, _ := config.GetConfig()
	originalProxyURL := configServe.GetString(common.PROXY_URL)
	configServe.SetString(common.PROXY_URL, mockServer.URL)

	app := newApp()
	h := iristest.New(t, app)

	Convey("Test Proxy POST /v1/files", t, func() {
		h.POST("/v1/files").
			WithJSON(map[string]interface{}{"name": "test.txt", "content": "Hello, World!"}).
			Expect().
			Status(iris.StatusOK).
			JSON().Object().Value("id").IsEqual("file-1")
	})

	configServe.SetString(common.PROXY_URL, originalProxyURL)
}

func TestProxyPostGenerateText(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/v1/generate" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"id": "gen-1", "object": "text", "created": 1610078138, "text": "Generated text"}`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	configServe, _ := config.GetConfig()
	originalProxyURL := configServe.GetString(common.PROXY_URL)
	configServe.SetString(common.PROXY_URL, mockServer.URL)

	app := newApp()
	h := iristest.New(t, app)

	Convey("Test Proxy POST /v1/generate", t, func() {
		h.POST("/v1/generate").
			WithJSON(map[string]interface{}{"prompt": "Hello, AI!"}).
			Expect().
			Status(iris.StatusOK).
			JSON().Object().Value("text").IsEqual("Generated text")
	})

	configServe.SetString(common.PROXY_URL, originalProxyURL)
}

func TestProxyPutUpdateFile(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && r.URL.Path == "/v1/files/file-1" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"id": "file-1", "object": "file", "updated": 1610078138}`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	configServe, _ := config.GetConfig()
	originalProxyURL := configServe.GetString(common.PROXY_URL)
	configServe.SetString(common.PROXY_URL, mockServer.URL)

	app := newApp()
	h := iristest.New(t, app)

	Convey("Test Proxy PUT /v1/files/file-1", t, func() {
		h.PUT("/v1/files/file-1").
			WithJSON(map[string]interface{}{"name": "updated.txt"}).
			Expect().
			Status(iris.StatusOK).
			JSON().Object().Value("id").IsEqual("file-1")
	})

	configServe.SetString(common.PROXY_URL, originalProxyURL)
}

func TestProxyDeleteFile(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete && r.URL.Path == "/v1/files/file-1" {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	configServe, _ := config.GetConfig()
	originalProxyURL := configServe.GetString(common.PROXY_URL)
	configServe.SetString(common.PROXY_URL, mockServer.URL)

	app := newApp()
	h := iristest.New(t, app)

	Convey("Test Proxy DELETE /v1/files/file-1", t, func() {
		h.DELETE("/v1/files/file-1").
			Expect().
			Status(iris.StatusNoContent)
	})

	configServe.SetString(common.PROXY_URL, originalProxyURL)
}
