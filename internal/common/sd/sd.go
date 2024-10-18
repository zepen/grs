package sd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// 服务发现

type SDServerConf struct {
	NacosAddr   string `yaml:"nacos_addr"`   // 注册中心地址
	ServicePort string `yaml:"service_port"` // 注册中心服务端口
	ClientPort  uint64 `yaml:"client_port"`  // 注册服务端口
	ServiceName string `yaml:"service_name"` // 注册服务名
	Scheme      string `yaml:"scheme"`
}

type SDServer struct {
	Conf    *SDServerConf
	nacos   naming_client.INamingClient
	offline bool
	mux     sync.RWMutex
}

func NewServiceDiscover(conf *SDServerConf) (*SDServer, error) {
	nacosAddrList := strings.Split(conf.NacosAddr, ",")
	postList := strings.Split(conf.ServicePort, ",")
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
	}
	var serverConfigs []constant.ServerConfig
	for i, _ := range nacosAddrList {
		serConf := constant.ServerConfig{}
		serConf.IpAddr = nacosAddrList[i]
		intPos, _ := strconv.Atoi(postList[i])
		serConf.Port = uint64(intPos)
		serConf.GrpcPort = uint64(intPos + 1000)
		serverConfigs = append(serverConfigs, serConf)
	}
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	//fmt.Println(serverConfigs)
	if err != nil {
		panic(err)
	}
	sdServer := &SDServer{
		Conf:  conf,
		nacos: client,
	}
	sdServer.Register()
	//go sdServer.RegisterInterval()
	return sdServer, nil
}

// 注册

func (c *SDServer) Register() {
	registerParam := vo.RegisterInstanceParam{}
	registerParam.Ip = IP()
	registerParam.ServiceName = c.Conf.ServiceName
	registerParam.Port = c.Conf.ClientPort
	registerParam.Weight = 10
	registerParam.Enable = true
	registerParam.Healthy = true
	registerParam.Ephemeral = true
	registerServiceInstance(c.nacos, registerParam)
}

// 定期注册

func (c *SDServer) RegisterInterval() {
	for {
		// 1分钟注册一次
		time.Sleep(time.Minute * 1)
		if c.IsOffline() {
			log.Println("offline stop register interval")
			break
		}
		updateParam := vo.UpdateInstanceParam{}
		updateParam.Ip = IP()
		updateParam.ServiceName = c.Conf.ServiceName
		updateParam.Port = c.Conf.ClientPort
		updateParam.Weight = 10
		updateParam.Enable = true
		updateParam.Healthy = true
		updateParam.Ephemeral = true
		updateServiceInstance(c.nacos, updateParam)
	}
}

func (c *SDServer) IsOffline() bool {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.offline
}

// 服务下线

func (c *SDServer) Deregister() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.offline = true
	deRegisterParam := vo.DeregisterInstanceParam{}
	deRegisterParam.Ip = IP()
	deRegisterParam.ServiceName = c.Conf.ServiceName
	deRegisterParam.Port = c.Conf.ClientPort
	deRegisterParam.Ephemeral = true
	deRegisterServiceInstance(c.nacos, deRegisterParam)
}

func registerServiceInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) {
	//fmt.Println(param)
	success, err := client.RegisterInstance(param)
	if !success || err != nil {
		panic("RegisterServiceInstance failed!" + err.Error())
	}
	fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func batchRegisterServiceInstance(client naming_client.INamingClient, param vo.BatchRegisterInstanceParam) {
	success, err := client.BatchRegisterInstance(param)
	if !success || err != nil {
		panic("BatchRegisterServiceInstance failed!" + err.Error())
	}
	fmt.Printf("BatchRegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func deRegisterServiceInstance(client naming_client.INamingClient, param vo.DeregisterInstanceParam) {
	success, err := client.DeregisterInstance(param)
	if !success || err != nil {
		panic("DeRegisterServiceInstance failed!" + err.Error())
	}
	fmt.Printf("DeRegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func updateServiceInstance(client naming_client.INamingClient, param vo.UpdateInstanceParam) {
	success, err := client.UpdateInstance(param)
	if !success || err != nil {
		panic("UpdateInstance failed!" + err.Error())
	}
	fmt.Printf("UpdateServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func getService(client naming_client.INamingClient, param vo.GetServiceParam) {
	service, err := client.GetService(param)
	if err != nil {
		panic("GetService failed!" + err.Error())
	}
	fmt.Printf("GetService,param:%+v, result:%+v \n\n", param, service)
}

func selectAllInstances(client naming_client.INamingClient, param vo.SelectAllInstancesParam) {
	instances, err := client.SelectAllInstances(param)
	if err != nil {
		panic("SelectAllInstances failed!" + err.Error())
	}
	fmt.Printf("SelectAllInstance,param:%+v, result:%+v \n\n", param, instances)
}

func selectInstances(client naming_client.INamingClient, param vo.SelectInstancesParam) {
	instances, err := client.SelectInstances(param)
	if err != nil {
		panic("SelectInstances failed!" + err.Error())
	}
	fmt.Printf("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
}

func selectOneHealthyInstance(client naming_client.INamingClient, param vo.SelectOneHealthInstanceParam) {
	instances, err := client.SelectOneHealthyInstance(param)
	if err != nil {
		panic("SelectOneHealthyInstance failed!")
	}
	fmt.Printf("SelectOneHealthyInstance,param:%+v, result:%+v \n\n", param, instances)
}

func subscribe(client naming_client.INamingClient, param *vo.SubscribeParam) {
	err := client.Subscribe(param)
	if err != nil {
		return
	}
}

func unSubscribe(client naming_client.INamingClient, param *vo.SubscribeParam) {
	err := client.Unsubscribe(param)
	if err != nil {
		return
	}
}

func getAllService(client naming_client.INamingClient, param vo.GetAllServiceInfoParam) {
	service, err := client.GetAllServicesInfo(param)
	if err != nil {
		panic("GetAllService failed!")
	}
	fmt.Printf("GetAllService,param:%+v, result:%+v \n\n", param, service)
}
