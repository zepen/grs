package internal

import (
	"context"
	"grs/conf"
	"grs/internal/model"
	"math"
	"time"
)

var decayMap map[string]func(a float32, k float32, t float32) float32

func init() {
	decayMap = make(map[string]func(a float32, k float32, t float32) float32)
	decayMap["linear_decay"] = linearDecay
	decayMap["exp_decay"] = expDecay
	decayMap["default"] = expDecay
}

func linearDecay(a float32, k float32, t float32) float32 {
	score := a - k*t
	if score > 0 {
		return score
	} else {
		return 0
	}
}

func expDecay(a float32, k float32, t float32) float32 {
	score := a * float32(math.Exp(-float64(k*t)))
	if score > 0.01 {
		return score
	} else {
		return 0
	}
}

func Merge(ctx context.Context, reqContext *model.NoteReqContext, noteList []*model.Note) []*model.Note {
	decayBase := float32(0.5)
	decayRate := float32(0.01)
	decayFuncName := "default"
	if reqContext.AbTestConf.Parameter != nil {
		decayBase = conf.GetFloatRecommendConf(reqContext.AbTestConf.Parameter.MergeConfig.Param, "decay_base")
		decayRate = conf.GetFloatRecommendConf(reqContext.AbTestConf.Parameter.MergeConfig.Param, "decay_rate")
		if _, ok := decayMap[reqContext.AbTestConf.Parameter.MergeConfig.FuncName]; ok {
			decayFuncName = reqContext.AbTestConf.Parameter.MergeConfig.FuncName
		}
	}
	for _, v := range noteList {
		d := int((time.Now().UnixMilli() - v.PubTime) / 86400000)
		timeScore := decayMap[decayFuncName](decayBase, decayRate, float32(d))
		v.MergeScore = (v.RecallScore.Score*timeScore + 0.5*timeScore) * v.RecallScore.Weight
		//ilog.Log.Infof("id %d, pubTime %d, timeScore %f, mergeScore %f, rw %f", v.NoteId, v.PubTime, timeScore, v.MergeScore, v.RecallScore.Weight)
	}
	return noteList
}
