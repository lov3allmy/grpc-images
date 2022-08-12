package internal

import (
	"context"
	"encoding/base64"
	pb "github.com/lov3allmy/tages/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Server struct {
	pb.UnimplementedImageStorageServiceServer
	repo *Repository
}

func NewServer(repo *Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {

	name := req.GetName()
	data := req.GetData()
	if len(data) > 10240 {
		return nil, status.Error(codes.InvalidArgument, "the image size exceeds 10 MB")
	}

	var encodedData []byte
	base64.StdEncoding.Encode(encodedData, data)

	id, err := s.repo.UploadImage(ctx, name, encodedData, time.Now())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot write changes to DB: %v", err)
	}

	res := &pb.UploadImageResponse{Id: id}

	return res, nil
}

func (s *Server) UpdateImage(ctx context.Context, req *pb.UpdateImageRequest) (*pb.UpdateImageResponse, error) {

	id := req.GetId()
	name := req.GetName()
	data := req.GetData()
	if len(data) > 10240 {
		return nil, status.Error(codes.InvalidArgument, "the image size exceeds 10 MB")
	}

	var encodedData []byte
	base64.StdEncoding.Encode(encodedData, data)

	err := s.repo.UpdateImage(ctx, id, name, data, time.Now())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot write changes to DB: %v", err)
	}

	res := &pb.UpdateImageResponse{}

	return res, nil
}

func (s *Server) DownloadImage(ctx context.Context, req *pb.DownloadImageRequest) (*pb.DownloadImageResponse, error) {

	id := req.GetId()

	encodedData, err := s.repo.DownloadImage(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot read image data from DB: %v", err)
	}

	var data []byte
	_, err = base64.StdEncoding.Decode(data, encodedData)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot decode image data: %v", err)
	}

	res := &pb.DownloadImageResponse{Data: data}

	return res, nil
}

func (s *Server) GetImagesList(ctx context.Context, _ *pb.GetImagesListRequest) (*pb.GetImagesListResponse, error) {

	list, err := s.repo.GetImagesList(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot read images info from DB: %v", err)
	}

	var resList []*pb.GetImagesListResponse_ImageInfo

	for _, imageInfo := range list {
		resList = append(resList, &pb.GetImagesListResponse_ImageInfo{
			Name:       imageInfo.name,
			CreatedAt:  timestamppb.New(imageInfo.createdAt),
			ModifiedAt: timestamppb.New(imageInfo.modifiedAt),
		})
	}

	res := &pb.GetImagesListResponse{ImageInfo: resList}

	return res, nil
}
