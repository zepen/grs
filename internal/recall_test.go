package internal

import (
	"fmt"
	"recommend-server/internal/model"
	"testing"
	"time"
)

func addWeightTest(note *model.Note, count *int) {
	arr := make([]string, 0)
	note.RecallScore.Score = note.RecallScore.Score + 0.1
	*count++
	arr = append(arr, "A")
	fmt.Printf("arr len = %d\n", len(arr))
}

func TestRecall(t *testing.T) {
	startTime := time.Now()
	arr := make([]*model.Note, 0, 10)
	for i := 0; i < 5; i++ {
		note := &model.Note{NoteId: uint64(i), Tags: nil}
		arr = append(arr, note)
	}
	fmt.Println(arr)
	fmt.Println(arr[:6])
	fmt.Printf("add cost time: %s\n", time.Since(startTime))
	//arrStr := strings.Join(arr, ",")
	//startTime = time.Now()
	//splitArr := strings.Split(arrStr, ",")
	//fmt.Printf("arr len = %d\n", len(splitArr))
	//fmt.Printf("cost time: %s\n", time.Since(startTime))
}

func TestRecallList(t *testing.T) {
	fmt.Printf("%d\n", time.Now().Unix())
	count := 0
	note := &model.Note{}
	recallScore := &model.RecallScore{}
	recallScore.Score = 0.0
	note.RecallScore = recallScore
	for i := 0; i < 10; i++ {
		fmt.Printf("Score: %f, Count: %d\n", note.RecallScore.Score, count)
		addWeightTest(note, &count)
		fmt.Printf("Add weight Score: %f, Count: %d\n", note.RecallScore.Score, count)
	}
}
