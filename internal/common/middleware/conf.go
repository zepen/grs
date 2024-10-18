package middleware

import (
	"database/sql"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/opensearch-project/opensearch-go"
	"time"
)

type RedisConf struct {
	IsCluster    bool          `yaml:"is_cluster"`
	Host         string        `yaml:"host"`
	Password     string        `yaml:"password"`
	MinIdleConn  int           `yaml:"min_idle_conn"` //最小活跃数
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolTimeout  time.Duration `yaml:"pool_timeout"`
	PoolSize     int           `yaml:"pool_size"`
	PartNum      int           `yaml:"part_num"`
	NodeNum      int           `yaml:"node_num"`
	Tls          bool          `yaml:"tls"`
}

type ElasticConf struct {
	Host string `yaml:"host"`
}

type OpenSearchConf struct {
	Host               string `yaml:"host"`
	User               string `yaml:"user"`
	Passwd             string `yaml:"passwd"`
	IsAws              bool   `yaml:"is_aws"`
	RecIndex           string `yaml:"rec_index"`
	LogIndex           string `yaml:"log_index"`
	AwsRegion          string `yaml:"aws_region"`
	AwsAccessKey       string `yaml:"aws_access_key"`
	AwsSecretAccessKey string `yaml:"aws_secret_access_key"`
}

type KafkaConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Tls  bool   `yaml:"tls"`
}

type MySqlConf struct {
	Host        string `yaml:"host"`
	User        string `yaml:"user"`
	Passwd      string `yaml:"passwd"`
	Dbname      string `yaml:"dbname"`
	MaxOpenConn int    `yaml:"max_open_conn"`
	MaxIdConn   int    `yaml:"max_id_conn"`
}

type MiddleWareControl struct {
	IsKafka      bool `yaml:"is_kafka"`
	IsRedis      bool `yaml:"is_redis"`
	IsOpenSearch bool `yaml:"is_open_search"`
	IsMySql      bool `yaml:"is_my_sql"`
}

type MiddleWareConf struct {
	KafkaConf         *KafkaConfig       `yaml:"kafka"`
	RedisConf         *RedisConf         `json:"redis" yaml:"redis"`
	OpenSearchConf    *OpenSearchConf    `yaml:"open_search"`
	ElasticConf       *ElasticConf       `yaml:"elastic"`
	MySqlConf         *MySqlConf         `yaml:"mysql"`
	MiddleWareControl *MiddleWareControl `yaml:"middle_ware_control"`
}

type MiddleWare struct {
	Conf            *MiddleWareConf
	RedisCli        *redis.Client
	RedisClusterCli *redis.ClusterClient
	OpenSearch      *opensearch.Client
	Mysql           *sql.DB
	KafkaProducer   sarama.SyncProducer
}
