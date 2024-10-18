package sort

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRemoveIndexInSlice(t *testing.T) {
	origin := []int{1, 2, 3, 4, 5}
	for removeIndex := 0; removeIndex < len(origin); removeIndex++ {
		fmt.Println("remove before:", origin)
		fmt.Println("remove index:", removeIndex)
		if removeIndex >= 0 && removeIndex < len(origin) {
			origin = append(origin[:removeIndex], origin[removeIndex+1:]...)
		}
		fmt.Println("remove after", origin)
		fmt.Println(strings.Repeat("-", 20))
	}
}

func TestCategoryIdShuffle(t *testing.T) {
	starTime := time.Now()
	note := []int{1, 1, 1, 1, 8, 8, 8, 8, 11, 11, 11, 11, 8, 8, 8, 8}
	window := 3
	ret := make([]int, 0) //结果
	if len(note) == 0 || len(note) <= window {
		fmt.Println("source:", note)
	}
	origin := make([]int, len(note))
	copy(origin, note) //拷贝下 一会需要删除
	categoryIdItem := make([]int, 0, window)
	for len(origin) > 0 {
		if len(categoryIdItem) >= window {
			categoryIdItem = categoryIdItem[1:]
		}
		indexMax := len(origin) - 1
		itemCategoryId := 0
		for i := 0; i <= indexMax; i++ {
			itemCategoryId = origin[i]
			if itemCategoryId <= 0 {
				continue
			}
			find := false
			// 遍历窗口中笔记
			for _, categoryId := range categoryIdItem {
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
		categoryIdItem = append(categoryIdItem, itemCategoryId)
		if indexMax >= 0 && indexMax < len(origin) {
			origin = append(origin[:indexMax], origin[indexMax+1:]...)
		}
	}
	fmt.Println("cost: ", time.Since(starTime))
	fmt.Println("res:", ret[:10])
}

func TestNoteReplace(t *testing.T) {
	note := []int{8, 4, 8, 8, 8, 4, 8, 8, 8, 8}
	shuffleNote := []int{8, 4, 14, 8, 4, 14, 8, 4, 10, 8}
	fmt.Println("source note: ", note)
	note = shuffleNote
	fmt.Println("replace note: ", shuffleNote)
}
