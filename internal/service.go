package internal

import (
	"context"
	"fmt"
	"gitlab.com/cher8/lion/common/sd"
	"recommend-server/apis"
	"recommend-server/conf"
	"time"
)

type Service struct {
	apis.UnimplementedRecommenderServer
	sd     *sd.SDServer
	conf   *conf.Conf
	Engine *Recommend
}

func (s *Service) Init(conf *conf.Conf) {
	ip := sd.IP()
	start := time.Now()
	s.conf = conf
	recommend, err := NewRecommend(context.Background(), conf)
	s.Engine = recommend
	if err != nil {
		panic("create recommend is fail!")
	}
	sdServer, err := sd.NewServiceDiscover(conf.SD)
	if err != nil {
		panic("create sdServer is fail!")
	}
	s.sd = sdServer
	cost := time.Since(start)
	fmt.Printf("new-service-cost:%v, start ip:%v\n", cost.String(), ip)
	fmt.Printf(LogoShow)
}

func (s *Service) RecommendServer(ctx context.Context, req *apis.UserRequest) (*apis.NoteResponse, error) {
	if pType, ok := req.Args["page_type"]; ok {
		if pType == "explore" {
			return s.ExploreApi(ctx, req)
		} else if pType == "following" {
			return s.FollowingApi(ctx, req)
		} else {
			return &apis.NoteResponse{}, nil
		}
	} else {
		return &apis.NoteResponse{}, nil
	}
}

func (s *Service) ExploreApi(ctx context.Context, req *apis.UserRequest) (*apis.NoteResponse, error) {
	if s.Engine != nil {
		resp, err := s.Engine.RecommendFlow(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	} else {
		return &apis.NoteResponse{}, nil
	}
}

func (s *Service) FollowingApi(ctx context.Context, req *apis.UserRequest) (*apis.NoteResponse, error) {
	return &apis.NoteResponse{}, nil
}

func (s *Service) Close() {
	fmt.Printf("%v:close", s.conf.ServiceName)
}
