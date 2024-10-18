package index

import (
	"gitlab.com/cher8/lion/common/ilog"
	"strings"
)

const findCapacity = 1000

type Index struct {
	Fwi *ForwardIndex
	Ivr *InvertedIndex
}

type ForwardIndex struct {
	Id2Cols map[string]map[string]string
}

type InvertedIndex struct {
	Cols2Id map[string]map[string]string
}

func (in *Index) getFwiIndexSize() int {
	return len(in.Fwi.Id2Cols)
}

func (in *Index) getIvrIndexColsSize() int {
	return len(in.Ivr.Cols2Id)
}

func (in *Index) getIvrIndexSize(colName string) int {
	return len(in.Ivr.Cols2Id[colName])
}

func (in *Index) BuildIndex(data []interface{}, forwardCols []string, invertedCols []string) {
	// 填入正排
	id2Cols := make(map[string]map[string]string)
	for _, v := range data {
		_id := v.(map[string]interface{})["_id"].(string)
		sourceMap := v.(map[string]interface{})["_source"].(map[string]interface{})
		fwiD := make(map[string]string)
		for _, col := range forwardCols {
			if sourceMap[col] != nil {
				fwiD[col] = sourceMap[col].(string)
			}
		}
		id2Cols[_id] = fwiD
	}
	in.Fwi.Id2Cols = id2Cols
	ilog.Log.Infof("rec-forward-index size = %d", in.getFwiIndexSize())
	// 填入倒排
	cols2Id := make(map[string]map[string]string)
	for _, col := range invertedCols {
		ivrD := make(map[string]string)
		for _, v := range data {
			_id := v.(map[string]interface{})["_id"].(string)
			sourceMap := v.(map[string]interface{})["_source"].(map[string]interface{})
			if sourceMap[col] != nil {
				valueList := strings.Split(sourceMap[col].(string), ",")
				for _, value := range valueList {
					if iv, ok := ivrD[value]; ok {
						ivrD[value] = iv + "," + _id
					} else {
						ivrD[value] = _id
					}
				}
				cols2Id[col] = ivrD
			}
		}
	}
	in.Ivr.Cols2Id = cols2Id
	for _, col := range invertedCols {
		ilog.Log.Infof("rec-inverted-index cols size = %d,  %s2Id size = %d",
			in.getIvrIndexColsSize(), col, in.getIvrIndexSize(col))
	}
}

func (in *Index) FindForwardIds(colName string, ids string) string {
	if in.Fwi != nil {
		findRes := make([]string, 0, findCapacity)
		idsList := strings.Split(ids, ",")
		for _, id := range idsList {
			if i2v, ok := in.Fwi.Id2Cols[id]; ok {
				if v, vOk := i2v[colName]; vOk {
					if v != "" {
						findRes = append(findRes, v)
					}
				}
			}
		}
		return strings.Join(findRes, ",")
	}
	return ""
}

func (in *Index) FindForwardIdsList(colName string, idsList []string) string {
	if in.Fwi != nil {
		findRes := make([]string, 0, findCapacity)
		for _, id := range idsList {
			if i2v, ok := in.Fwi.Id2Cols[id]; ok {
				if v, vOk := i2v[colName]; vOk {
					if v != "" {
						findRes = append(findRes, v)
					}
				}
			}
		}
		return strings.Join(findRes, ",")
	}
	return ""
}

func (in *Index) FindInvertedIds(colName string, values string) string {
	if in.Ivr != nil {
		findRes := make([]string, 0, findCapacity)
		valuesList := strings.Split(values, ",")
		if v2i, ok := in.Ivr.Cols2Id[colName]; ok {
			for _, value := range valuesList {
				if id, vOk := v2i[value]; vOk {
					if id != "" {
						findRes = append(findRes, id)
					}
				}
			}
		}
		return strings.Join(findRes, ",")
	}
	return ""
}

func (in *Index) FindInvertedIdsList(colName string, valuesList []string) string {
	if in.Ivr != nil {
		findRes := make([]string, 0, findCapacity)
		if v2i, ok := in.Ivr.Cols2Id[colName]; ok {
			for _, value := range valuesList {
				if id, vOk := v2i[value]; vOk {
					if id != "" {
						findRes = append(findRes, id)
					}
				}
			}
		}
		return strings.Join(findRes, ",")
	}
	return ""
}
