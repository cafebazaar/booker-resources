package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/cafebazaar/booker-resources/proto"
)

type resourcesServer struct{}

func (s *resourcesServer) GetCategories(ctx context.Context, in *pb.CategoriesGetRequest) (*pb.CategoriesGetReply, error) {
	categoryNames, err := getCategories()
	if err != nil {
		//TODO: fix error codes
		return &pb.CategoriesGetReply{ReplyProperties: &pb.ReplyProperties{StatusCode: pb.ReplyProperties_SERVER_ERROR}}, err
	}
	categoriesList := make([]*pb.Category, 0)
	for _, v := range categoryNames {
		categoriesList = append(categoriesList, &pb.Category{Name: v})
	}
	resp := &pb.CategoriesGetReply{
		ReplyProperties: pb.ReplyPropertiesTemplate(),
		Categories:      categoriesList,
	}
	return resp, nil
}

func (s *resourcesServer) GetCategoryItems(ctx context.Context, in *pb.CategoryItemsGetRequest) (*pb.CategoryItemsGetReply, error) {
	items, err := getCategoryItems(in.Name)
	respItems := make([]*pb.Item, 0)
	for _, item := range items {
		respItems = append(respItems, &pb.Item{Name: item.Name, Spec: item.Spec})
	}
	if err != nil {
		//TODO: fix error codes
		return &pb.CategoryItemsGetReply{ReplyProperties: &pb.ReplyProperties{StatusCode: pb.ReplyProperties_SERVER_ERROR}}, err
	}
	resp := &pb.CategoryItemsGetReply{
		ReplyProperties: pb.ReplyPropertiesTemplate(),
		Items:           respItems,
	}
	return resp, nil

}

func (s *resourcesServer) GetItem(ctx context.Context, in *pb.ItemGetRequest) (*pb.ItemGetReply, error) {
	println(":", in.CategoryName, in.ItemName)
	item, err := getItem(in.CategoryName, in.ItemName)
	if err != nil || item == nil {
		//TODO: fix error codes
		return &pb.ItemGetReply{ReplyProperties: &pb.ReplyProperties{StatusCode: pb.ReplyProperties_SERVER_ERROR}}, err
	}
	resp := &pb.ItemGetReply{
		ReplyProperties: pb.ReplyPropertiesTemplate(),
		Item:            &pb.Item{Name: item.Name, Spec: item.Spec},
	}
	return resp, nil
}

func (s *resourcesServer) PostItem(ctx context.Context, in *pb.ItemPostRequest) (*pb.ItemPostReply, error) {
	item, err := createItem(in.CategoryName, in.GetItem().Name, in.GetItem().GetSpec())
	if err != nil {
		return &pb.ItemPostReply{ReplyProperties: &pb.ReplyProperties{StatusCode: pb.ReplyProperties_SERVER_ERROR}}, err
	}
	resp := &pb.ItemPostReply{
		ReplyProperties: pb.ReplyPropertiesTemplate(),
		CategoryName:    item.Name,
	}
	return resp, nil

}

func RegisterServer(grpcServer *grpc.Server) {
	pb.RegisterResourcesServer(grpcServer, new(resourcesServer))
}
