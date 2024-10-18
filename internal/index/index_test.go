package index

import (
	"fmt"
	"testing"
	"time"
)

func TestIndex(t *testing.T) {
	stTime := time.Now()
	arr := make([]string, 0, 100000)
	for i := 0; i < 40000; i++ {
		arr = append(arr, "A")
	}
	fmt.Println(time.Since(stTime).Microseconds())
}
