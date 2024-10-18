package sort

//粗排
import (
	"context"
	"gitlab.com/cher8/lion/common/ilog"
	"google.golang.org/grpc"
	"math/rand"
	"recommend-server/internal/model"
	"recommend-server/internal/sort/estimate"
	"sort"
	"time"
)

const (
	maxCoarseLimitDefault = 100 //最大上限
)

// 粗排函数
var rankMap map[string]func(req *model.NoteReqContext, noteList []*model.Note)

func init() {
	rankMap = make(map[string]func(req *model.NoteReqContext, noteList []*model.Note))
	rankMap["default"] = coarseLinerScore //基准
	rankMap["lr"] = estimateScore
}

// CoarseWrap 粗排闭包
func CoarseWrap(ctx context.Context, reqContext *model.NoteReqContext, noteList []*model.Note) {
	funcName := "default"
	if reqContext.AbTestConf.Parameter != nil {
		funcName = reqContext.AbTestConf.Parameter.RankConfig.FuncName
	}
	if f, ok := rankMap[funcName]; ok {
		f(reqContext, noteList)
	} else {
		//TODO 找不到粗排函数 告警
	}
}

// 线性打分
func coarseLinerScore(req *model.NoteReqContext, noteList []*model.Note) {
	rand.Seed(time.Now().UnixNano())
	randFloat := rand.Float64() // 生成 0 到 1 之间的随机浮点数
	for _, v := range noteList {
		v.CoarseScore = v.MergeScore * float32(randFloat)
	}
	CoarserScoreSortFunc := func(i, j int) bool {
		return noteList[i].CoarseScore > noteList[j].CoarseScore
	}
	sort.Slice(noteList, CoarserScoreSortFunc)
	//for _, v := range noteList {
	//	fmt.Println(v.NoteId, v.Tags.ChildTag, v.CoarseScore, v.MergeScore, v.RecallScore.Name, v.RecallScore.Score, v.RecallScore.Weight, v.RecallScore.Len)
	//}
	//fmt.Println(len(noteList))
}

func newRequestFeatures(req *model.NoteReqContext, noteList []*model.Note) *estimate.EstimateRequest {
	// 创建请求对象
	er := &estimate.EstimateRequest{}
	uf := make(map[string]string)
	uf_ := &estimate.UserFeatures{
		UserId:   req.UserId,
		Features: uf,
	}
	nfs := make([]*estimate.NoteFeatures, 0, maxCoarseLimitDefault)
	for _, note := range noteList {
		nf := make(map[string]string)
		nf_ := &estimate.NoteFeatures{
			NoteId:   note.NoteId,
			Features: nf,
		}
		nfs = append(nfs, nf_)
	}
	er.Uf = uf_
	er.Nf = nfs
	return er
}

// 调用排序模型
func estimateScore(req *model.NoteReqContext, noteList []*model.Note) {
	if c, ok := req.AbTestConf.Parameter.RankConfig.Param["estimate_addr"]; ok {
		conn, err := grpc.Dial(c.(string), grpc.WithInsecure())
		if err != nil {
			ilog.Log.Warnf("did not connect: %v", err)
		} else {
			client := estimate.NewEstimatorClient(conn)
			coarseLinerScore(req, noteList)
			estReq := newRequestFeatures(req, noteList)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			r, err2 := client.EstimatorResp(ctx, estReq)
			if err2 != nil {
				ilog.Log.Warnf("could not conn estimate: %v", err2)
			} else {
				for _, note := range noteList {
					if score, ok2 := r.Outputs[note.NoteId]; ok2 {
						note.CoarseScore = score
					} else {
						note.CoarseScore = 0.0
					}
				}
				CoarserScoreSortFunc := func(i, j int) bool {
					return noteList[i].CoarseScore > noteList[j].CoarseScore
				}
				sort.Slice(noteList, CoarserScoreSortFunc)
			}
		}
		defer conn.Close()
	} else {
		coarseLinerScore(req, noteList)
	}
}
