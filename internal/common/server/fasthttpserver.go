package server

import (
	"fmt"
	"time"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

const (
	DefaultReadTimeoutMs  = 1000
	DefaultWriteTimeoutMs = 1000
	DefaultMaxConn        = 10000
)

type FastHttpServerConf struct {
	Name                 string `yaml:"name"` //根据名字生成metrics等系统级路由
	Addr                 string `yaml:"addr"`
	MaxConn              int    `yaml:"max_conn"` //最大连接数
	ReadTimeMillSeconds  int64  `yaml:"read_time_ms"`
	WriteTimeMillSeconds int64  `yaml:"write_time_ms"`
}

type FastHttpServer struct {
	addr   string
	server fasthttp.Server
}

func validation(conf *FastHttpServerConf) error {
	if conf == nil || conf.Name == "" || conf.Addr == "" {
		return fmt.Errorf("fasthttp conf error %v", conf)
	}
	return nil
}

func NewFastHttpServer(conf *FastHttpServerConf, router *router.Router) (*FastHttpServer, error) {
	if err := validation(conf); err != nil {
		return nil, err
	}
	if conf.MaxConn <= 0 {
		conf.MaxConn = DefaultMaxConn
	}
	if conf.ReadTimeMillSeconds <= 0 {
		conf.ReadTimeMillSeconds = DefaultReadTimeoutMs
	}
	if conf.WriteTimeMillSeconds <= 0 {
		conf.WriteTimeMillSeconds = DefaultWriteTimeoutMs
	}
	router.GET("/metrics", prometheusHandler())
	router.GET(fmt.Sprintf("/%v/health", conf.Name), health)
	s := &FastHttpServer{
		addr: conf.Addr,
		server: fasthttp.Server{
			Concurrency:        conf.MaxConn,
			Handler:            router.Handler,
			ReadTimeout:        time.Duration(conf.ReadTimeMillSeconds) * time.Millisecond,
			WriteTimeout:       time.Duration(conf.WriteTimeMillSeconds) * time.Millisecond,
			TCPKeepalive:       true,
			IdleTimeout:        time.Second * 90,
			TCPKeepalivePeriod: time.Second * 90,
		},
	}
	err := s.Run()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *FastHttpServer) Run() error {
	go func() {
		err := s.server.ListenAndServe(s.addr)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

func (s *FastHttpServer) Stop() error {
	return s.server.Shutdown()
}

func prometheusHandler() fasthttp.RequestHandler {
	return fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
}

func health(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(200)
	_, err := ctx.WriteString("true")
	if err != nil {
		return
	}
	return
}
