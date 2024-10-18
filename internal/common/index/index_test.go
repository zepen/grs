package index

import (
	"fmt"
	"testing"
	"time"
)

func readAndWrite(i *IdDocIdMap, num uint32) {
	go func() {
		i.mux.RLock()
		fmt.Printf("read:%d\n", i.Data["a"])
		i.mux.RUnlock()

		i.mux.Lock()
		i.Data["a"] = num
		i.mux.Unlock()

		i.mux.RLock()
		fmt.Printf("write:%d\n", i.Data["a"])
		i.mux.RUnlock()
	}()
}

func a(x int) {
	fmt.Printf("x=%d, address=%d\n", x, &x)
}

func b(x *int) {
	fmt.Printf("x=%d, address=%d, p-address=%d\n", *x, x, &x)
}

func TestIndex(t *testing.T) {
	i := IdDocIdMap{}
	i.Data = make(map[string]uint32)
	i.Data["a"] = 1
	readAndWrite(&i, 1)
	readAndWrite(&i, 2)
	readAndWrite(&i, 3)
	readAndWrite(&i, 4)
	time.Sleep(time.Second * 3)
}
