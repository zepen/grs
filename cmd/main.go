package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"grs/apis"
	"grs/conf"
	"grs/internal"
	"grs/internal/common/sd"
	readYaml "grs/internal/common/yaml"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "conf", "", "conf values")
}

func main() {
	flag.Parse()
	readNacRead, _ := readYaml.ReadYamlLocal(confPath)
	naConfig := &sd.SDServerConf{}
	err := yaml.Unmarshal(readNacRead, naConfig)
	if err != nil {
		panic(err)
	}
	log.Printf("\nThis is %s env\n", os.Getenv("REC_ENV"))
	yamlRead, err := readYaml.ReadYamlRemote(naConfig)
	if err != nil {
		panic(err)
	}
	log.Printf("Conf from remote, conf len = %d\n\n", len(yamlRead))
	config := &conf.Conf{}
	err1 := yaml.Unmarshal(yamlRead, config)
	if err1 != nil {
		panic(err1)
	}
	s := &internal.Service{}
	config.SD = naConfig
	s.Init(config)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", naConfig.ClientPort))
	if err != nil {
		panic(err)
	}
	grpc.NumStreamWorkers(uint32(config.Grpc.MaxWorkers))
	grpc.MaxConcurrentStreams(uint32(config.Grpc.MaxConcurrentStreams))
	gs := grpc.NewServer()
	apis.RegisterRecommenderServer(gs, s)
	log.Printf("server listening at %v\n", lis.Addr())
	if err = gs.Serve(lis); err != nil {
		fmt.Printf("failed to serve, %v\n", err)
		panic(err)
	}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	go func() {
		sig := <-sigs
		log.Println(sig)
		s.Close()
		done <- true
	}()
	log.Println("awaiting signal...")
	<-done
	log.Println(fmt.Sprintf("%v:stop", config.ServiceName))
}
