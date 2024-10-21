package internal

import (
	"context"
	"grs/internal/common/ilog"
	"grs/internal/common/middleware"
	"grs/internal/common/util"
	"grs/internal/model"
	"grs/internal/recall"
	"os"
	"strconv"
	"strings"
)

const MaxDoc = 10000
const MinDoc = 1000

type NoteIds struct {
	DistinctMap       map[uint64]struct{} `json:"distinct_map"`
	NoteList          []*model.Note       `json:"note_list"`
	UserIntentMap     map[string]int      `json:"user_intent_map"`
	UserDisLikeAuthor map[string]struct{} `json:"user_dis_like_author"`
	UserDisLikeTags   map[string]struct{} `json:"user_dis_like_tags"`
}

func NewTags(tag *model.Tags, tags interface{}) {
	if tags == "" {
		// 如果没有标签，随机生成一个父标签，用于打散排序
		tag.ParentTag = util.RandOneLabel("", 16)
	} else {
		tagsList := strings.Split(tags.(string), ",")
		if len(tagsList) == 1 {
			tag.ParentTag = tagsList[0]
		} else if len(tagsList) == 2 {
			tag.ParentTag = tagsList[0]
			tag.ChildTag = tagsList[1]
		} else if len(tagsList) == 3 {
			tag.ParentTag = tagsList[0]
			tag.ChildTag = tagsList[1]
			if tagsList[2] == "500" {
				tag.SexTag = "2"
			} else if tagsList[2] == "501" {
				tag.SexTag = "0"
			} else if tagsList[2] == "502" {
				tag.SexTag = "1"
			}
		}
	}
}

func NewNote(subMap map[string]interface{}, recallName string, recallTotal int) *model.Note {
	noteId, _ := strconv.ParseInt(subMap["id"].(string), 10, 64)
	_score, _ := strconv.ParseFloat(subMap["score"].(string), 64)
	score := float32(_score)
	noteTag := &model.Tags{}
	noteTag.ParentTag = ""
	noteTag.ChildTag = ""
	noteTag.SexTag = ""
	if c, ok := subMap["tags"]; ok {
		NewTags(noteTag, c)
	}
	author := subMap["author"].(string)
	pubTime, _ := strconv.ParseInt(subMap["pub_time"].(string), 10, 64)
	weight := float32(0)
	if r, ok := recall.RcMap[recallName]; ok {
		weight = r.Weight
	}
	channelName := ""
	if c, ok := subMap["channel_name"]; ok {
		channelName = c.(string)
	}
	recallScore := &model.RecallScore{}
	recallScore.Name = recallName
	recallScore.Len = recallTotal
	recallScore.Score = score
	recallScore.Weight = weight
	note := &model.Note{}
	note.NoteId = uint64(noteId)
	note.Author = author
	note.Tags = noteTag
	note.PubTime = pubTime
	note.ChannelName = channelName
	note.RecallScore = recallScore
	note.MergeScore = 0.0
	note.CoarseScore = 0.0
	return note
}

func (n *NoteIds) addIntentWeight(note *model.Note) {
	intentWeight := float32(0.0)
	ciw, cOk := n.UserIntentMap[note.Tags.ChildTag]
	piw, pOk := n.UserIntentMap[note.Tags.ParentTag]
	if cOk && !pOk { // 子标签命中，父标签未命中，取子标签
		intentWeight = float32(ciw)
	} else if (cOk && pOk) || (!cOk && pOk) { // 子父标签全部命中，取父标签权重
		intentWeight = float32(piw)
	}
	note.RecallScore.Score = note.RecallScore.Score * intentWeight
}

func (n *NoteIds) addNoteList(reqContext *model.NoteReqContext, recallName string, note *model.Note, count *int) {
	//去掉重复的笔记
	if _, ok := n.DistinctMap[note.NoteId]; !ok {
		// 过滤掉用户不喜欢作者
		if _, disLikeOk := n.UserDisLikeAuthor[note.Author]; !disLikeOk {
			// 命中用户不喜欢的子标签降权
			if _, disLikeTagOk := n.UserDisLikeTags[note.Tags.ChildTag]; disLikeTagOk {
				note.RecallScore.Score = note.RecallScore.Score * 0.01
			}
			// 如果学校地点加权
			if note.ChannelName == reqContext.User.UserSeqInfo.UserFromChannel {
				note.RecallScore.Score = note.RecallScore.Score * 10
			}
			// 如果帖子带性别标签，性别加权
			if note.Tags.SexTag != "" && reqContext.Sex != "" && reqContext.Sex == note.Tags.SexTag {
				n.addIntentWeight(note)
				note.RecallScore.Score = note.RecallScore.Score + 0.35
			} else {
				n.addIntentWeight(note)
			}
			// 如果随机召回，不喜欢子标签直接过滤
			if recallName == "random_recall" {
				if _, disLikeTagOk := n.UserDisLikeTags[note.Tags.ChildTag]; !disLikeTagOk {
					n.NoteList = append(n.NoteList, note)
					n.DistinctMap[note.NoteId] = struct{}{}
					*count++
				}
			} else {
				n.NoteList = append(n.NoteList, note)
				n.DistinctMap[note.NoteId] = struct{}{}
				*count++
			}
		}
	}
}

func (n *NoteIds) AddNoteList(ctx context.Context, reqContext *model.NoteReqContext,
	recallName string, recallTotal int, subHits []map[string]interface{}) {
	if recallTotal > 0 && subHits != nil {
		count := 0
		for _, v := range subHits {
			note := NewNote(v, recallName, recallTotal)
			n.addNoteList(reqContext, recallName, note, &count)
		}
		ilog.Log.Infof("%s add noteList and len = %d", recallName, count)
	}
}

func (n *NoteIds) CompletingRecallList(ctx context.Context, reqContext *model.NoteReqContext) {
	noteListLen := len(n.NoteList)
	if noteListLen < MinDoc {
		// 前面召回数量不够，随机补
		diffCount := MinDoc - noteListLen
		_list, _len := recall.RandomRecallList("", diffCount, reqContext)
		n.AddNoteList(ctx, reqContext, "random_recall", _len, _list)
	} else if noteListLen > MaxDoc {
		// 超过最大召回数，截断
		n.NoteList = n.NoteList[:MaxDoc]
	}
}

// 读取召回策略
func newRecallStrategy(reqContext *model.NoteReqContext) []*recall.Recall {
	recallList := make([]*recall.Recall, 0)
	recallStrategy := strings.Split(recall.Default, ",")
	if reqContext.AbTestConf.Parameter != nil {
		recallStrategy = strings.Split(reqContext.AbTestConf.Parameter.RecallConfig.RecallStrategy, ",")
	}
	for _, v := range recallStrategy {
		vSplit := strings.Split(v, ":")
		recallName := vSplit[0]
		recallWeight, _ := strconv.ParseFloat(vSplit[1], 64)
		ilog.Log.Infof("recall_name: %s, weight: %f", recallName, recallWeight)
		if rc, ok := recall.RcMap[recallName]; ok {
			rc.UpdateWeight(float32(recallWeight))
			// random_recall 为默认召回
			if recallName != "random_recall" {
				recallList = append(recallList, rc)
			}
		}
	}
	return recallList
}

func ExploreRecall(ctx context.Context, reqContext *model.NoteReqContext) []*model.Note {
	// 首页召回
	noteList := &NoteIds{
		DistinctMap:       make(map[uint64]struct{}),
		NoteList:          make([]*model.Note, 0, MaxDoc), //1w cap避免频繁扩容
		UserIntentMap:     make(map[string]int),
		UserDisLikeAuthor: make(map[string]struct{}),
		UserDisLikeTags:   make(map[string]struct{}),
	}
	// 添加用户浏览曝光 && 被封禁帖子列表
	if reqContext.User.UserSeqInfo != nil && reqContext.User.UserSeqInfo.ViewSeq != "" {
		viewSeq := strings.Split(reqContext.User.UserSeqInfo.ViewSeq, ",")
		for _, noteId := range viewSeq {
			if noteId != "" {
				uintNoteId, _ := strconv.ParseUint(noteId, 10, 64)
				noteList.DistinctMap[uintNoteId] = struct{}{}
			}
		}
		if os.Getenv("REC_ENV") == "prd" {
			fnList, fnErr := middleware.M.RedisClusterCli.SMembers(ctx, model.NoteForbiddenKeyName).Result()
			if fnErr != nil {
				ilog.Log.Warnf("key %s has error = %v\n", model.NoteForbiddenKeyName, fnErr)
			} else {
				for _, noteId := range fnList {
					noteId = strings.Replace(noteId, "\"", "", -1)
					if noteId != "" {
						uintNoteId, _ := strconv.ParseUint(noteId, 10, 64)
						noteList.DistinctMap[uintNoteId] = struct{}{}
					}
				}
			}
		} else {
			fnList, fnErr := middleware.M.RedisCli.SMembers(ctx, model.NoteForbiddenKeyName).Result()
			if fnErr != nil {
				ilog.Log.Warnf("key %s has error = %v\n", model.NoteForbiddenKeyName, fnErr)
			} else {
				for _, noteId := range fnList {
					noteId = strings.Replace(noteId, "\"", "", -1)
					if noteId != "" {
						uintNoteId, _ := strconv.ParseUint(noteId, 10, 64)
						noteList.DistinctMap[uintNoteId] = struct{}{}
					}
				}
			}
		}
	}
	// 添加用户dislikeNote列表
	if reqContext.User.UserSeqInfo != nil && reqContext.User.UserSeqInfo.DislikeNoteSeq != "" {
		dislikeNoteSeq := strings.Split(reqContext.User.UserSeqInfo.DislikeNoteSeq, ",")
		for _, noteId := range dislikeNoteSeq {
			if noteId != "" {
				uintNoteId, _ := strconv.ParseUint(noteId, 10, 64)
				noteList.DistinctMap[uintNoteId] = struct{}{}
			}
		}
	}
	// 添加用户dislike作者列表 && 被封禁作者列表
	if reqContext.User.UserSeqInfo != nil && reqContext.User.UserSeqInfo.DislikeAuthorSeq != "" {
		dislikeAuthorSeq := strings.Split(reqContext.User.UserSeqInfo.DislikeAuthorSeq, ",")
		for _, authorId := range dislikeAuthorSeq {
			if authorId != "" {
				noteList.UserDisLikeAuthor[authorId] = struct{}{}
			}
		}
		if os.Getenv("REC_ENV") == "prd" {
			fuList, fuErr := middleware.M.RedisClusterCli.SMembers(ctx, model.UserForbiddenKeyName).Result()
			if fuErr != nil {
				ilog.Log.Warnf("key %s has error = %v\n", model.UserForbiddenKeyName, fuErr)
			} else {
				for _, authorId := range fuList {
					authorId = strings.Replace(authorId, "\"", "", -1)
					if authorId != "" {
						noteList.UserDisLikeAuthor[authorId] = struct{}{}
					}
				}
			}
		} else {
			fuList, fuErr := middleware.M.RedisCli.SMembers(ctx, model.UserForbiddenKeyName).Result()
			if fuErr != nil {
				ilog.Log.Warnf("key %s has error = %v\n", model.UserForbiddenKeyName, fuErr)
			} else {
				for _, authorId := range fuList {
					authorId = strings.Replace(authorId, "\"", "", -1)
					if authorId != "" {
						noteList.UserDisLikeAuthor[authorId] = struct{}{}
					}
				}
			}
		}
	}
	// 添加用户dislikeTags列表
	if reqContext.User.UserSeqInfo != nil && reqContext.User.UserSeqInfo.DislikeTagsSeq != "" {
		dislikeTagsSeq := strings.Split(reqContext.User.UserSeqInfo.DislikeTagsSeq, ",")
		for _, tag := range dislikeTagsSeq {
			if tag != "" {
				noteList.UserDisLikeTags[tag] = struct{}{}
			}
		}
	}
	// 查找用户兴趣标签
	if reqContext.User.UserSeqInfo != nil && reqContext.User.UserSeqInfo.IntentList != "" {
		intentList := strings.Split(reqContext.User.UserSeqInfo.IntentList, ",")
		intentListCopy := make([]string, 0)
		for _, v := range intentList {
			splitV := strings.Split(v, ":")
			if len(splitV) > 1 {
				// 取最后一位作为权重
				weight, _ := strconv.Atoi(splitV[len(splitV)-1])
				noteList.UserIntentMap[splitV[0]] = weight
				intentListCopy = append(intentListCopy, splitV[0])
			} else {
				// 仅为了兼容无权重模式序列
				noteList.UserIntentMap[v] = 1
			}
		}
		// 重建用户兴趣标签序列
		reqContext.User.UserSeqInfo.IntentList = strings.Join(intentListCopy, ",")
	}
	recallList := newRecallStrategy(reqContext)
	// redis和openSearch开放状态
	if middleware.M.Conf.MiddleWareControl.IsRedis && middleware.M.Conf.MiddleWareControl.IsOpenSearch {
		for _, rc := range recallList {
			_list, _len := recall.RFunc(rc.Name, reqContext, rc.Col, 0, rc.Func)
			noteList.AddNoteList(ctx, reqContext, rc.Name, _len, _list)
		}
		// 构建最终召回列表
		noteList.CompletingRecallList(ctx, reqContext)
		return noteList.NoteList
	} else {
		// TODO
	}
	if reqContext.IsDebug {
		// TODO 上下问传透debug
	}
	return nil
}

func (r *Recommend) FollowingRecall(ctx context.Context, reqContext *model.NoteReqContext) error {
	// TODO 实现关注页召回逻辑
	return nil
}
