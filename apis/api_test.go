package apis

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

func TestRecommendApi(t *testing.T) {
	address := "127.0.0.1:10010"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := NewRecommenderClient(conn)
	s := time.Now()
	userReq := &UserRequest{}
	userReq.UserId = "1635093838685372418"
	args := make(map[string]string)
	args["page_type"] = "explore"
	userReq.Args = args
	r, err := client.RecommendServer(context.Background(), userReq)
	if err != nil {
		log.Fatalf("recommend server is fail: %v", err)
	}
	log.Println(r.NoteIds)
	log.Printf("cost time: %s\n\n", time.Since(s))
}
