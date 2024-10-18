package index

import (
	"sync"
)

type IndexConf struct { //配置
	Id      string   `yaml:"id"` //index的唯一id
	IvtList []string `yaml:"ivt_list"`
	OSSPath string   `yaml:"oss_path"`
	//OssClient         *ossutil.OSSConf     `yaml:"oss_client"`
	RefreshFullSecond int64 `json:"refresh_full_second" yaml:"refresh_full_second"` //刷新全量 秒
	Dynamic           bool  `yaml:"dynamic"`
	//DynamicKafka      *KafkaConsumerConfig `yaml:"dynamic_kafka"`
}

type IdDocIdMap struct {
	mux  sync.RWMutex
	Data map[string]uint32
}

type Index struct {
	Id           string               //index id alpha beta
	ivt          map[string]*sync.Map //倒排索引  key：域 value： (k:term;v:bitmap)
	fwd          []*map[string]string //正排索引，下标是docId
	idDocIdMap   *IdDocIdMap          //支持并发读写
	version      string               //全量索引版本  for eg.  ：goods-index/20210603/1441
	ts           int64                //全量索引的时间戳
	maxFullDocId uint32               //全量索引最大的docId
	maxDocId     uint32               //索引最大的docId
	//deleteRR               *roaring.Bitmap //删除的压缩位图
	//missRelateRR           *roaring.Bitmap //相关性miss的压缩位图
	//KafkaConsumer          *KafkaConsumer
	Dynamic                bool //是否增量
	DynamicCount           int
	DynamicDeleteCount     int
	DynamicInsertCount     int
	DynamicMissRelateScore int //增量的相关性打分miss数量
	DynamicHitRelateScore  int //赠浪的存在intentId 可以计算相关性分数的监控
}

type IndexWrapper struct {
	Id    string
	Conf  *IndexConf
	mux   sync.RWMutex //读写锁 读多写少
	Index *Index       //内存索引
	//OSSClient *ossutil.OSSClient //ossclient
}
