package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/cafebazaar/booker-resources/proto"
)

type resourcesServer struct{}

func (s *resourcesServer) GetCategories(ctx context.Context, in *pb.CategoriesGetRequest) (*pb.CategoriesGetReply, error) {
	resp := &pb.CategoriesGetReply{
		ReplyProperties: pb.ReplyPropertiesTemplate(),
	}
	return resp, nil
}

func (s *resourcesServer) GetCategoryItems(ctx context.Context, in *pb.CategoryItemsGetRequest) (*pb.CategoryItemsGetReply, error) {
	resp := &pb.CategoryItemsGetReply{
		ReplyProperties: pb.ReplyPropertiesTemplate(),
	}
	return resp, nil

}

func (s *resourcesServer) GetItem(ctx context.Context, in *pb.ItemGetRequest) (*pb.ItemGetReply, error) {
	resp := &pb.ItemGetReply{
		ReplyProperties: pb.ReplyPropertiesTemplate(),
	}
	return resp, nil
}

func (s *resourcesServer) PostItem(ctx context.Context, in *pb.ItemPostRequest) (*pb.ItemPostReply, error) {
	resp := &pb.ItemPostReply{
		ReplyProperties: pb.ReplyPropertiesTemplate(),
	}
	return resp, nil

}

func RegisterServer(grpcServer *grpc.Server) {
	pb.RegisterResourcesServer(grpcServer, new(resourcesServer))
}
