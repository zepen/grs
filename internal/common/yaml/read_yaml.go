package yaml

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"io/ioutil"
	"path/filepath"
	"recommend-server/internal/common/sd"
	"strconv"
	"strings"
)

func ReadYamlLocal(confPath string) ([]byte, error) {
	if confPath == "" {
		panic("conf error")
	}
	yamlFile, err := filepath.Abs(confPath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	yamlRead, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	return yamlRead, nil
}

func ReadYamlRemote(conf *sd.SDServerConf) ([]byte, error) {
	nacosAddrList := strings.Split(conf.NacosAddr, ",")
	postList := strings.Split(conf.ServicePort, ",")
	var sc []constant.ServerConfig
	for i, _ := range nacosAddrList {
		serConf := constant.ServerConfig{}
		serConf.IpAddr = nacosAddrList[i]
		intPos, _ := strconv.Atoi(postList[i])
		serConf.Port = uint64(intPos)
		serConf.GrpcPort = uint64(intPos + 1000)
		sc = append(sc, serConf)
	}
	cc := constant.ClientConfig{
		NamespaceId:         "",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "log",
		CacheDir:            "cache",
		LogLevel:            "debug",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: conf.ServiceName + ".yaml",
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		panic(err)
		return nil, err
	}
	return []byte(content), nil
}
