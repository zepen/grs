package middleware

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
	"log"
	"net/http"
	"strings"
)

var M *MiddleWare // 中间键

func NewRedisClient(conf *RedisConf) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Host,
		Password:     conf.Password,
		DB:           0,
		MaxRetries:   0,
		MinIdleConns: conf.MinIdleConn,
		DialTimeout:  conf.DialTimeout,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		PoolTimeout:  conf.PoolTimeout,
		PoolSize:     conf.PoolSize,
	})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		log.Printf("[Error] Redis connection is fail: %s", err)
		return nil
	}
	return client
}

func NewRedisClusterClient(conf *RedisConf) *redis.ClusterClient {
	var addrs []string
	for i := 1; i <= conf.PartNum; i++ {
		for j := 1; j <= conf.NodeNum; j++ {
			addrs = append(addrs, fmt.Sprintf("redis-cluster-000%d-00%d.%s", i, j, conf.Host))
		}
	}
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: conf.Password,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: conf.Tls,
		},
	})
	_, err := clusterClient.Ping(context.TODO()).Result()
	if err != nil {
		log.Printf("[Error] Redis connection is fail: %s", err)
		return nil
	}
	return clusterClient
}

func NewElasticClient(conf *ElasticConf) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: strings.Split(conf.Host, ","),
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("[Error] ElasticSearch create is fail: %s", err)
		return nil
	}
	res, err := es.Info()
	if err != nil {
		log.Fatalf("[Error] No elasticSearch info: %s", err)
		return nil
	} else {
		fmt.Printf("%s", res)
	}
	return es
}

func getCredentialProvider(accessKey, secretAccessKey, token string) aws.CredentialsProviderFunc {
	return func(ctx context.Context) (aws.Credentials, error) {
		c := &aws.Credentials{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretAccessKey,
			SessionToken:    token,
		}
		return *c, nil
	}
}

func NewOpenSearchClient(conf *OpenSearchConf) *opensearch.Client {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Addresses: strings.Split(conf.Host, ","),
		Username:  conf.User,
		Password:  conf.Passwd,
	})
	if err != nil {
		log.Fatalf("[Error] OpenSearch connection is fail: %s", err)
		return nil
	}
	info, err := client.Info()
	if err != nil {
		log.Fatalf("[Error] OpenSearch connection is fail, miss info: %s", err)
		return nil
	} else {
		log.Fatalf("%s", info)
	}
	return client
}

func NewAwsOpenSearchClient(conf *OpenSearchConf) *opensearch.Client {
	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(conf.AwsRegion),
		config.WithCredentialsProvider(
			getCredentialProvider(conf.AwsAccessKey, conf.AwsSecretAccessKey, ""),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := awsv2.NewSignerWithService(awsCfg, "es")
	if err != nil {
		log.Fatal(err)
	}
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{conf.Host},
		Signer:    signer,
	})
	if err != nil {
		log.Fatal("client creation err", err)
	}
	return client
}

func NewMySqlClient(conf *MySqlConf) *sql.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		conf.User,
		conf.Passwd,
		conf.Host,
		conf.Dbname,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("connect mysql fail!")
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		panic("connect mysql fail!")
	}
	return db
}

func NewKafkaProducer(conf *KafkaConfig) sarama.SyncProducer {
	_config := sarama.NewConfig()
	_config.Producer.RequiredAcks = sarama.WaitForAll
	_config.Producer.Retry.Max = 5
	_config.Producer.Return.Successes = true
	if conf.Tls {
		_config.Net.TLS.Enable = true
		_config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	producer, err := sarama.NewSyncProducer([]string{fmt.Sprintf("%s:%s", conf.Host, conf.Port)}, _config)
	if err != nil {
		log.Printf("Failed to create Kafka producer: %s", err)
		return nil
	}
	return producer
}

func NewMiddleWare(conf *MiddleWareConf) error {
	middleWare := &MiddleWare{}
	middleWare.Conf = conf
	if conf.MiddleWareControl.IsMySql {
		mysql := NewMySqlClient(conf.MySqlConf)
		if mysql != nil {
			log.Printf("create mysql connection completed!")
			middleWare.Mysql = mysql
		}
	}
	if conf.MiddleWareControl.IsRedis {
		if conf.RedisConf.IsCluster {
			redisClusterCli := NewRedisClusterClient(conf.RedisConf)
			if redisClusterCli != nil {
				log.Printf("create redis cluster connection completed!")
				middleWare.RedisClusterCli = redisClusterCli
			}
		} else {
			redisCli := NewRedisClient(conf.RedisConf)
			if redisCli != nil {
				log.Printf("create redis connection completed!")
				middleWare.RedisCli = redisCli
			}
		}
	}
	if conf.MiddleWareControl.IsOpenSearch {
		if conf.OpenSearchConf.IsAws {
			openS := NewAwsOpenSearchClient(conf.OpenSearchConf)
			if openS != nil {
				log.Printf("create AWS openSearch connection completed!")
				middleWare.OpenSearch = openS
			}
		} else {
			openS := NewOpenSearchClient(conf.OpenSearchConf)
			if openS != nil {
				log.Printf("create openSearch connection completed!")
				middleWare.OpenSearch = openS
			}
		}
	}
	if conf.MiddleWareControl.IsKafka {
		kafkaProducer := NewKafkaProducer(conf.KafkaConf)
		if kafkaProducer != nil {
			log.Printf("create kafka producer completed!")
			middleWare.KafkaProducer = kafkaProducer
		}
	}
	M = middleWare
	return nil
}
