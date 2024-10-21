package recall

import (
	"grs/internal/common/ilog"
	"grs/internal/index"
	"grs/internal/model"
	"strings"
	"time"
)

const RcCapacity = 1000
const Default = "new_content_recall:0.5,user_follow_recall:0.5,best_content_recall:0.5,random_recall:0.4"

type Recall struct {
	Name   string
	Col    string
	Weight float32
	Func   func(colName string, count int, reqContext *model.NoteReqContext) ([]map[string]interface{}, int)
}

func (rc *Recall) UpdateCol(col string) {
	rc.Name = col
}

func (rc *Recall) UpdateWeight(weight float32) {
	rc.Weight = weight
}

var RcMap = map[string]*Recall{
	"i2i_content_recall":  {"i2i_content_recall", "similar_note_ids", 0.0, i2iRecallList},
	"new_content_recall":  {"new_content_recall", "", 0.0, newRecallList},
	"hot_content_recall":  {"hot_content_recall", "", 0.0, hotRecallList},
	"user_follow_recall":  {"user_follow_recall", "user_id", 0.0, userFollowRecallList},
	"best_content_recall": {"best_content_recall", "group_id", 0.0, bestRecallList},
	"user_tags_recall":    {"user_tags_recall", "tags", 0.0, userTagsRecallList},
	"user_cf_recall":      {"user_cf_recall", "", 0.0, userCFRecallList},
	"random_recall":       {"random_recall", "", 0.0, RandomRecallList},
}

func NewNoteMap(id string) map[string]interface{} {
	noteMap := make(map[string]interface{})
	noteMap["id"] = id
	noteMap["score"] = index.RecIndex.FindForwardIds("score", id)
	noteMap["tags"] = index.RecIndex.FindForwardIds("tags", id)
	noteMap["author"] = index.RecIndex.FindForwardIds("user_id", id)
	noteMap["group_id"] = index.RecIndex.FindForwardIds("group_id", id)
	noteMap["pub_time"] = index.RecIndex.FindForwardIds("pub_time", id)
	noteMap["channel_name"] = index.RecIndex.FindForwardIds("channel_name", id)
	return noteMap
}

func i2iRecallList(colName string, count int, reqContext *model.NoteReqContext) ([]map[string]interface{}, int) {
	ids := index.RecIndex.FindForwardIds(colName, reqContext.User.UserSeqInfo.ClickSeq)
	noteList := make([]map[string]interface{}, 0, RcCapacity)
	distinctMap := make(map[string]struct{})
	if len(ids) > 0 {
		similarNoteList := strings.Split(ids, ",")
		for _, noteId := range similarNoteList {
			if _, ok := distinctMap[noteId]; !ok {
				noteList = append(noteList, NewNoteMap(noteId))
			}
		}
		return noteList, len(noteList)
	}
	return noteList, 0
}

func newRecallList(colName string, count int, reqContext *model.NoteReqContext) ([]map[string]interface{}, int) {
	noteList := make([]map[string]interface{}, 0, RcCapacity)
	noteNewList := strings.Split(reqContext.NoteSummary.NoteNewRank, ",")
	if len(noteNewList) > 0 {
		for _, noteId := range noteNewList {
			noteList = append(noteList, NewNoteMap(noteId))
		}
		return noteList, len(noteList)
	}
	return noteList, 0
}

func hotRecallList(colName string, count int, reqContext *model.NoteReqContext) ([]map[string]interface{}, int) {
	noteList := make([]map[string]interface{}, 0, RcCapacity)
	noteHotList := strings.Split(reqContext.NoteSummary.NoteHotRank, ",")
	if len(noteHotList) > 0 {
		for _, noteId := range noteHotList {
			noteList = append(noteList, NewNoteMap(noteId))
		}
		return noteList, len(noteList)
	}
	return noteList, 0
}

func userCFRecallList(colName string, count int, reqContext *model.NoteReqContext) ([]map[string]interface{}, int) {
	noteList := make([]map[string]interface{}, 0, RcCapacity)
	userCFList := strings.Split(reqContext.User.UserSeqInfo.UserCFList, ",")
	if len(userCFList) > 0 {
		for _, noteId := range userCFList {
			noteList = append(noteList, NewNoteMap(noteId))
		}
		return noteList, len(noteList)
	}
	return noteList, 0
}

func userFollowRecallList(colName string, count int, reqContext *model.NoteReqContext) ([]map[string]interface{}, int) {
	ids := index.RecIndex.FindInvertedIds(colName, reqContext.User.UserSeqInfo.UserFollowSeq)
	noteList := make([]map[string]interface{}, 0, RcCapacity)
	if len(ids) > 0 {
		userFollowList := strings.Split(ids, ",")
		for _, noteId := range userFollowList {
			noteList = append(noteList, NewNoteMap(noteId))
		}
		return noteList, len(noteList)
	}
	return noteList, 0
}

func bestRecallList(colName string, count int, reqContext *model.NoteReqContext) ([]map[string]interface{}, int) {
	noteList := make([]map[string]interface{}, 0, RcCapacity)
	if bestGroupId, ok := reqContext.Conf["best_group_id"]; ok {
		ids := index.RecIndex.FindInvertedIds(colName, bestGroupId.(string))
		if len(ids) > 0 {
			userFollowList := strings.Split(ids, ",")
			for _, noteId := range userFollowList {
				noteList = append(noteList, NewNoteMap(noteId))
			}
			return noteList, len(noteList)
		}
		return noteList, 0
	}
	return noteList, 0
}

func userTagsRecallList(colName string, count int, reqContext *model.NoteReqContext) ([]map[string]interface{}, int) {
	ids := index.RecIndex.FindInvertedIds(colName, reqContext.User.UserSeqInfo.IntentList)
	noteList := make([]map[string]interface{}, 0, RcCapacity)
	distinctMap := make(map[string]struct{})
	if len(ids) > 0 {
		userTagsList := strings.Split(ids, ",")
		for _, noteId := range userTagsList {
			if _, ok := distinctMap[noteId]; !ok {
				noteList = append(noteList, NewNoteMap(noteId))
			}
		}
		return noteList, len(noteList)
	}
	return noteList, 0
}

func RandomRecallList(colName string, count int, reqContext *model.NoteReqContext) ([]map[string]interface{}, int) {
	noteList := make([]map[string]interface{}, 0, RcCapacity)
	i := 0
	for noteId := range index.RecIndex.Fwi.Id2Cols {
		noteList = append(noteList, NewNoteMap(noteId))
		if i == count {
			break
		}
		i++
	}
	return noteList, len(noteList)
}

func RFunc(recallName string, reqContext *model.NoteReqContext, colName string, count int,
	f func(colName string, count int, reqContext *model.NoteReqContext) ([]map[string]interface{}, int)) ([]map[string]interface{}, int) {
	listCl := make(chan []map[string]interface{}, 1)
	lenCl := make(chan int, 1)
	go func() {
		start := time.Now()
		_list, _len := f(colName, count, reqContext)
		listCl <- _list
		lenCl <- _len
		cost := time.Since(start)
		ilog.Log.Infof("%s cost time: %s", recallName, cost.String())
	}()
	_len := <-lenCl
	_list := <-listCl
	ilog.Log.Infof("%s Len = %d", recallName, _len)
	return _list, _len
}
