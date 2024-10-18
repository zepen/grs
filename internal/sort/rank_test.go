package sort

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"recommend-server/internal/sort/estimate"
	"testing"
	"time"
)

func Slice(a []int) {
	a = a[:10]
	fmt.Println(len(a))
}

func TestSlice(t *testing.T) {
	a := make([]int, 0, 100)
	i := 0
	for {
		a = append(a, i)
		if i > 10 {
			break
		}
		i++
	}
	fmt.Println(len(a))
	Slice(a)
	fmt.Println(len(a))
	//fmt.Println(len(a[:10]))
	//fmt.Println(a[:10])
	//fmt.Println(len(a[10:]))
	//fmt.Println(a[10:])
	//fmt.Println(a[:100])

}

func TestEstimateScore(t *testing.T) {
	conn, err := grpc.Dial("localhost:10013", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := estimate.NewEstimatorClient(conn)
	//log.Printf("conn cost time: %s\n\n", time.Since(s))
	for j := 0; j < 10000; j++ {
		s := time.Now()
		esReq := &estimate.EstimateRequest{}
		uf := make(map[string]string)
		uf_ := &estimate.UserFeatures{
			UserId:   111,
			Features: uf,
		}
		nfs := make([]*estimate.NoteFeatures, 0, 100)
		for i := 0; i < 1000; i++ {
			nf := make(map[string]string)
			nf_ := &estimate.NoteFeatures{
				NoteId:   1122,
				Features: nf,
			}
			nfs = append(nfs, nf_)
		}
		esReq.Uf = uf_
		esReq.Nf = nfs
		r, err := client.EstimatorResp(context.Background(), esReq)
		if err != nil {
			log.Fatalf("could not estimate: %v", err)
		}
		log.Println(r.Outputs)
		log.Printf("cost time: %s\n\n", time.Since(s))
	}
}
