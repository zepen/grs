package eshook

import (
	"context"
	"fmt"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/sirupsen/logrus"
	"recommend-server/internal/common/util"
	"strings"
)

type OpenSearchHook struct {
	Client *opensearch.Client
	Index  string
}

func (hook *OpenSearchHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *OpenSearchHook) Fire(entry *logrus.Entry) error {
	doc := make(map[string]interface{})
	for k, v := range entry.Data {
		doc[k] = v
	}
	doc["@timestamp"] = entry.Time.Format("2006-01-02T15:04:05.000Z")
	doc["message"] = entry.Message
	body := strings.NewReader(util.ToString(doc))
	docId := util.GetNowTimeFmt() + " " + util.RandomString(5)
	go func() {
		req := opensearchapi.IndexRequest{
			Index:      hook.Index + "_" + util.GenToday(),
			DocumentID: docId,
			Body:       body,
		}
		_, err := req.Do(context.Background(), hook.Client)
		if err != nil {
			fmt.Printf("failed to index document: %v", err)
		}
	}()
	return nil
}
