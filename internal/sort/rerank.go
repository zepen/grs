package sort

import (
	"context"
	"gitlab.com/cher8/lion/common/ilog"
	"recommend-server/conf"
	"recommend-server/internal/model"
	"strconv"
)

var rerankMap map[string]func(req *model.NoteReqContext, note []*model.Note, window int)

func init() {
	rerankMap = make(map[string]func(req *model.NoteReqContext, note []*model.Note, window int))
	rerankMap["default"] = parentTagShuffle // 笔记父类标签打散
}

func ReRank(ctx context.Context, reqContext *model.NoteReqContext, noteList []*model.Note) {
	funcName := "default"
	window := 3
	if reqContext.AbTestConf.Parameter != nil {
		funcName = reqContext.AbTestConf.Parameter.ReRankConfig.FuncName
		window = int(conf.GetFloatRecommendConf(reqContext.AbTestConf.Parameter.ReRankConfig.Param, "window"))
	}
	if f, ok := rerankMap[funcName]; ok {
		f(reqContext, noteList, window)
	} else {
		//TODO 找不到重排函数 告警
	}
}

func parentTagShuffle(req *model.NoteReqContext, note []*model.Note, window int) {
	ret := make([]*model.Note, 0) //结果
	if len(note) > maxCoarseLimitDefault {
		note = note[:maxCoarseLimitDefault] // 截取头部
	}
	if len(note) == 0 || len(note) <= window {
		ilog.Log.Warnf("No noteList re-rank!")
	}
	origin := make([]*model.Note, len(note))
	copy(origin, note) //拷贝下 一会需要删除
	tagItem := make([]int, 0, window)
	for len(origin) > 0 {
		if len(tagItem) >= window {
			tagItem = tagItem[1:]
		}
		indexMax := len(origin) - 1
		itemCategoryId := 0
		for i := 0; i <= indexMax; i++ {
			itemCategoryId, _ = strconv.Atoi(origin[i].Tags.ParentTag)
			if itemCategoryId <= 0 {
				continue
			}
			find := false
			// 遍历窗口中笔记
			for _, categoryId := range tagItem {
				if itemCategoryId == categoryId {
					find = true
					break
				}
			}
			// 如果没找到，最大数组下标赋值为当前下标
			if !find {
				indexMax = i
				break
			}
		}
		//结果里面加入数据
		ret = append(ret, origin[indexMax])
		tagItem = append(tagItem, itemCategoryId)
		origin = removeIndexInSlice(origin, indexMax)
	}
	//for _, v := range ret[:100] {
	//	fmt.Println(v.CategoryId)
	//}
	// 打散列表替换原列表
	copy(note, ret)
}

func removeIndexInSlice(origin []*model.Note, index int) []*model.Note {
	if index >= 0 && index < len(origin) {
		origin = append(origin[:index], origin[index+1:]...)
	}
	return origin
}
