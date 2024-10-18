package recall

import (
	"fmt"
	"testing"
)

func TestRecallFunc(t *testing.T) {
	fmt.Println(
		RcMap["random_recall"].Weight,
	)
	RcMap["random_recall"].UpdateWeight(0.5)
	fmt.Println(
		RcMap["random_recall"].Weight,
	)
	noteList := make([]int, 0, 1000)
	fmt.Println(len(noteList))
	for i := 0; i < 1101; i++ {
		noteList = append(noteList, 100)
	}
	fmt.Println(len(noteList))
}
