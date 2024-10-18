package es

import (
	"fmt"
	"gitlab.com/cher8/lion/common/ilog"
	"strings"
)

const MaxDocSize = 1000000

func AllRecallQueryBody(sourceCols string) *strings.Reader {
	dsl := fmt.Sprintf(
		`{
                  "size": %d,
				  "query": {
					"bool": {
					  "must": [
						{
						   "term": {
						   "deleted": false
						  }
						},
                        {
                           "term": {
                           "disabled": 0
                          }
                        },
						{
						  "terms": {
							"note_type": [0, 1]
						  }
						}
					  ]
					}
				  }, "_source": [%s]
				}`, MaxDocSize, sourceCols)
	ilog.Log.Infof("all recall %s", dsl)
	return strings.NewReader(dsl)
}

func CreateTimeRecallQueryBody(recallSize int, sourceCols string) *strings.Reader {
	dsl := fmt.Sprintf(
		`{
				  "size": %d,
				  "sort": [
					{"create_time": {"order": "desc"}}
				  ],
				  "query": {
					"match_all": {
					}
				  }, "_source": [%s]
				}`, recallSize, sourceCols)
	ilog.Log.Infof("create_time recall %s", dsl)
	return strings.NewReader(dsl)
}

func LabelRecallQueryBody(labels string, maxNotes int) *strings.Reader {
	dsl := fmt.Sprintf(
		`{
				  "from": 0, 
					"size": %d,
					"query": {
					  "terms": {
						"category_id": [%s]
					  }
				  }, "_source": "id"}`, maxNotes, labels)
	ilog.Log.Infof("label recall %s", dsl)
	return strings.NewReader(dsl)
}

func I2iRecallQueryBody(ids string, col string, maxNotes int) *strings.Reader {
	dsl := fmt.Sprintf(
		`{
                  "from": 0,
                  "size": %d,
		          "query": {
                     "ids": {
                        "values": [%s]
                     }
	              }, "_source": ["%s"]}`, maxNotes, ids, col)
	ilog.Log.Infof("i2i recall %s", dsl)
	return strings.NewReader(dsl)
}

func U2u2iRecallQueryBody(terms string, col string, maxNotes int) *strings.Reader {
	dsl := fmt.Sprintf(
		`{
                  "from": 0,
                  "size": %d,
                  "query": {
                  "terms": {
                     "user_id": [%s]
                     }
                  }, "_source": "%s"}}`, maxNotes, terms, col)
	ilog.Log.Infof("u2u2i recall %s", dsl)
	return strings.NewReader(dsl)
}

func RandomRecallQueryBody(size int) *strings.Reader {
	dsl := fmt.Sprintf(
		`{
                  "from": 0, 
                  "size": %d,
				  "query": {
					"function_score": {
					  "query": {
						"match_all": {}
					  },
					  "random_score": {}, 
					  "boost_mode": "replace"
					}
				  }, "_source": ["id", "category_id"]
				}`, size)
	ilog.Log.Infof("random recall %s", dsl)
	return strings.NewReader(dsl)
}
