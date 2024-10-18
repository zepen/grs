package conf

import (
	"fmt"
	"gitlab.com/cher8/lion/common/ilog"
	"gitlab.com/cher8/lion/common/middleware"
	"gitlab.com/cher8/lion/common/robot"
	"gitlab.com/cher8/lion/common/sd"
	//"runtime/trace"
)

type GrpcConf struct {
	MaxWorkers           int `yaml:"max_workers"`
	MaxConcurrentStreams int `yaml:"max_concurrent_streams"`
	connectionTimeout    int `yaml:"connection_timeout"`
}

type Conf struct {
	ServiceName   string                     `yaml:"service_name"`
	SD            *sd.SDServerConf           `yaml:"sd"`
	Logger        *ilog.LogConf              `yaml:"logger"`
	Middleware    *middleware.MiddleWareConf `yaml:"middleware"`
	RecommendConf map[string]interface{}     `yaml:"recommend_conf"`
	Robot         *robot.Robot               `yaml:"robot"`
	Grpc          *GrpcConf                  `yaml:"grpc"`
}

func GetStringRecommendConf(recommendConf map[string]interface{}, name string) string {
	col := ""
	if v, ok := recommendConf[name]; ok {
		col = v.(string)
	} else {
		panic(fmt.Sprintf("not found %s from nacos conf!", name))
	}
	return col
}

func GetFloatRecommendConf(recommendConf map[string]interface{}, name string) float32 {
	col := float32(0)
	if v, ok := recommendConf[name]; ok {
		col = float32(v.(float64))
	} else {
		panic(fmt.Sprintf("not found %s from nacos conf!", name))
	}
	return col
}

func GetBoolRecommendConf(recommendConf map[string]interface{}, name string) bool {
	col := false
	if v, ok := recommendConf[name]; ok {
		col = v.(bool)
	} else {
		panic(fmt.Sprintf("not found %s from nacos conf!", name))
	}
	return col
}

func GetIntRecommendConf(recommendConf map[string]interface{}, name string) int {
	col := 0
	if v, ok := recommendConf[name]; ok {
		col = v.(int)
	} else {
		panic(fmt.Sprintf("not found %s from nacos conf!", name))
	}
	return col
}
