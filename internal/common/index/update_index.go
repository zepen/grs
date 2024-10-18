package index

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"recommend-server/internal/common/ilog"
	"recommend-server/internal/common/middleware"
	"strings"
	"time"
)

var RecIndex *Index

func generateDsl(dsl string) *strings.Reader {
	ilog.Log.Infof("all recall %s", dsl)
	return strings.NewReader(dsl)
}

func generateParseList(resMap map[string]interface{}) ([]interface{}, int) {
	if hits, ok := resMap["hits"]; ok {
		totalValue := int(hits.(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
		subHits := hits.(map[string]interface{})["hits"]
		return subHits.([]interface{}), totalValue
	} else {
		subHits := make([]interface{}, 0)
		return subHits, 0
	}
}

func query2OpenSearch(ctx context.Context, queryBody *strings.Reader) map[string]interface{} {
	var resMap map[string]interface{}
	// 构建查询语句，在OpenSearch中查询
	search := opensearchapi.SearchRequest{
		Index:   []string{middleware.M.Conf.OpenSearchConf.RecIndex},
		Body:    queryBody,
		Timeout: 5 * time.Minute,
	}
	searchResp, err := search.Do(ctx, middleware.M.OpenSearch)
	if err != nil {
		ilog.Log.Errorf("No recall result! s%", err)
	} else {
		var buf bytes.Buffer
		_, err = buf.ReadFrom(searchResp.Body)
		if err != nil {
			ilog.Log.Errorf("Reading the response body:", err)
		}
		if err := json.Unmarshal(buf.Bytes(), &resMap); err != nil {
			ilog.Log.Errorf("Parsing the response body:", err)
		}
		defer searchResp.Body.Close()
	}
	return resMap
}

//func NewIndex(mode string, forwardColsStr string, invertedColsStr string, args ...string) {
//	forwardCols := strings.Split(forwardColsStr, ",")
//	invertedCols := strings.Split(invertedColsStr, ",")
//	var queryCols []string
//	for _, col := range forwardCols {
//		queryCols = append(queryCols, `"`+col+`"`)
//	}
//	// 全量召回
//	recallQueryBody := generateDsl(strings.Join(queryCols, ","))
//	recallResultMap := query2OpenSearch(context.Background(), recallQueryBody)
//	recallList, listLen := generateParseList(recallResultMap)
//	if mode == "init" {
//		if listLen == 0 {
//			panic("Data is empty, please check openSearch or query!")
//		}
//		RecIndex = &Index{
//			Fwi: &ForwardIndex{
//				Id2Cols: map[string]map[string]string{},
//			},
//			Ivr: &InvertedIndex{
//				Cols2Id: map[string]map[string]string{},
//			},
//		}
//		RecIndex.BuildIndex(recallList, forwardCols, invertedCols)
//	} else {
//		updateIndex(recallList, forwardCols, invertedCols, args[0])
//	}
//}
//
//func updateIndex(recallList []interface{}, forwardCols []string, invertedCols []string, url string) {
//	fs := &robot.FeishuRobot{Url: url}
//	if len(recallList) != 0 {
//		RecIndex.BuildIndex(recallList, forwardCols, invertedCols)
//		message := fmt.Sprintf(
//			"%s, %s, Updated index success, Fwi-size = %d Ivr-size cols = %d",
//			util.GetNowTimeFmt(),
//			os.Getenv("REC_ENV"),
//			RecIndex.getFwiIndexSize(),
//			RecIndex.getIvrIndexColsSize())
//		fs.Send2robot(message)
//	} else {
//		message := fmt.Sprintf(
//			"%s, %s, Updated index fail, please check! old Fwi-size = %d old Ivr-size cols = %d",
//			util.GetNowTimeFmt(),
//			os.Getenv("REC_ENV"),
//			RecIndex.getFwiIndexSize(),
//			RecIndex.getIvrIndexColsSize())
//		fs.Send2robot(message)
//	}
//}
