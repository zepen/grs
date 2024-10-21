package internal

import (
	"context"
	"grs/apis"
	"grs/conf"
	"grs/internal/common/abtest"
	"grs/internal/common/ilog"
	"grs/internal/common/middleware"
	"grs/internal/dump"
	"grs/internal/index"
	"grs/internal/model"
	"grs/internal/monitor"
	"grs/internal/sort"
	"strconv"
	"strings"
	"time"
)

const limitSize = 10

type Recommend struct {
	Conf map[string]interface{}
}

func NewRecommend(ctx context.Context, conf_ *conf.Conf) (*Recommend, error) {
	// 构建中间件
	err := middleware.NewMiddleWare(conf_.Middleware)
	if err != nil {
		return nil, err
	}
	// 初始化日志
	ilog.InitLog(conf_.Logger)
	// 构建索引
	forwardColsStr := conf.GetStringRecommendConf(conf_.RecommendConf, "forward_cols")
	invertedColsStr := conf.GetStringRecommendConf(conf_.RecommendConf, "inverted_cols")
	index.NewIndex("init", forwardColsStr, invertedColsStr)
	fsRobotUrl := conf_.Robot.FeishuRobot.Url
	// 定时更新索引
	uiTime := conf.GetIntRecommendConf(conf_.RecommendConf, "update_index_time")
	go monitor.UpdateIndex(forwardColsStr, invertedColsStr, uiTime, fsRobotUrl)
	// 定时检查服务健康情况
	cshTime := conf.GetIntRecommendConf(conf_.RecommendConf, "check_server_health_time")
	go monitor.CheckServerHealth(fsRobotUrl, cshTime)
	// 定时删除清理用户浏览列表
	duvkTime := conf.GetIntRecommendConf(conf_.RecommendConf, "delete_user_view_key_time")
	go monitor.DeletedUserViewKeys(ctx, duvkTime, fsRobotUrl)
	// 创建推荐对象
	recommend := &Recommend{}
	recommend.Conf = conf_.RecommendConf
	return recommend, nil
}

func NewRecommendResp(ctx context.Context, reqContext *model.NoteReqContext, resp *apis.NoteResponse, noteList []*model.Note) {
	resp.NoteIds.Size = int32(len(noteList))
	ilog.Log.Infof("return noteList len = %d", len(noteList))
	noteListStr := make([]string, 0, 10)
	if len(noteList) > limitSize {
		noteList = noteList[:limitSize]
	}
	for _, v := range noteList {
		resp.NoteIds.RList = append(resp.NoteIds.RList, v.NoteId)
		noteListStr = append(noteListStr, strconv.FormatUint(v.NoteId, 10))
		newRecallTypeDump := dump.NewRecallTypeDump(reqContext, v)
		go dump.SendInfo2Kafka(newRecallTypeDump)
		ilog.Log.Infof("%d, %s, %s, %s, %d, %f, %f, %s, %f, %f, %d",
			v.NoteId, v.Tags.ParentTag, v.Tags.ChildTag, v.Tags.SexTag, v.PubTime, v.CoarseScore,
			v.MergeScore, v.RecallScore.Name, v.RecallScore.Score, v.RecallScore.Weight, v.RecallScore.Len)
	}
	dump.SendExploreNoteString2Redis(ctx, reqContext, strings.Join(noteListStr, ","))
}

func (r *Recommend) RecommendFlow(ctx context.Context, req *apis.UserRequest) (*apis.NoteResponse, error) {
	resp := &apis.NoteResponse{}
	resp.Version = req.Args["version"]
	resp.NoteIds = &apis.NoteList{}
	resp.NoteIds.TabName = req.Args["page_type"]
	userId, _ := strconv.Atoi(req.UserId)
	// 创建用户对象
	user := &model.User{}
	user.UserId = uint64(userId)
	user.UserSeqInfo = &model.UserSeqInfo{}
	// 获取redis实验配置信息, 选中实验组别
	whiteList := abtest.GetWhiteList(model.AbTestWhiteListKeyName)
	abTestConfig := abtest.GetExpConfig(model.AbTestKeyName)
	ek := &abtest.ExpGroup{}
	if whiteList != nil && abTestConfig != nil {
		// 判断user_id是否配置白名单
		if c, ok := whiteList.Wl[req.UserId]; ok {
			for _, expKey := range abTestConfig.ExpGroups {
				if c == expKey.Name {
					ek = expKey
				}
			}
		} else {
			hashNum := abtest.GenerateHashNum(100, int(abTestConfig.Seed), req.UserId)
			for _, expKey := range abTestConfig.ExpGroups {
				if hashNum >= expKey.Start && hashNum < expKey.End {
					ek = expKey
				}
			}
		}
	}
	// 获取用户redis信息
	user.FindUserViewSeq(ctx, req)
	user.FindUserIntentList(ctx, req)
	user.FindUserClickSeq(ctx, req)
	user.FindUserDisLikeNoteSeq(ctx, req)
	user.FindUserDisLikeAuthorSeq(ctx, req)
	user.FindUserDisLikeTagsSeq(ctx, req)
	user.FindUserFollowSeq(ctx, req)
	user.FindUserCFList(ctx, req)
	user.FindUserFromChannel(ctx, req)
	// 创建帖子对象
	noteSummary := &model.NoteSummary{}
	// 获取帖子redis信息
	noteSummary.FindNoteNewRankList(ctx)
	noteSummary.FindNoteHotRankList(ctx)
	// 构建请求对象
	reqC := &model.NoteReqContext{}
	reqC.ReqId = req.Args["req_id"]
	reqC.Version = req.Args["version"]
	reqC.UserId = uint64(userId)
	reqC.Sex = req.Args["sex"]
	reqC.IP = req.Args["ip"]
	reqC.OS = req.Args["os"]
	reqC.IsDebug = false
	reqC.UniqueId = req.Args["req_id"] + req.UserId //唯一id 用于dump
	reqC.Timestamp = time.Now().Unix() * 1000       //时间戳到毫秒
	reqC.User = user
	reqC.NoteSummary = noteSummary
	reqC.Conf = r.Conf
	reqC.AbTestConf = ek
	// 召回
	recallStart := time.Now()
	noteList := ExploreRecall(ctx, reqC)
	recallCost := time.Since(recallStart)
	ilog.Log.Infof("recall cost: %s", recallCost.String())
	// 融合召回
	mergeStart := time.Now()
	noteList = Merge(ctx, reqC, noteList)
	mergeCost := time.Since(mergeStart)
	ilog.Log.Infof("merge cost: %s", mergeCost.String())
	// 排序
	rankStart := time.Now()
	sort.CoarseWrap(ctx, reqC, noteList)
	rankCost := time.Since(rankStart)
	ilog.Log.Infof("rank cost: %s", rankCost.String())
	// 重排序
	reRankStart := time.Now()
	sort.ReRank(ctx, reqC, noteList)
	reRankCost := time.Since(reRankStart)
	ilog.Log.Infof("reRank cost: %s", reRankCost.String())
	// 构造返回列表
	newRecRespStart := time.Now()
	NewRecommendResp(ctx, reqC, resp, noteList)
	newRecRespCost := time.Since(newRecRespStart)
	ilog.Log.Infof("newRecRespStart cost: %s", newRecRespCost.String())
	return resp, nil
}
