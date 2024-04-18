package service

import (
	"context"

	"github.com/94peter/serial_number/model"

	"github.com/94peter/serial_number/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewGcp(cfg *model.Config) pb.GcpServiceServer {
	return &gcp{
		serialMgr: cfg.NewSerial(),
	}
}

type gcp struct {
	pb.UnimplementedGcpServiceServer

	serialMgr model.SerialMgr
}

func (g *gcp) CreatePrefix(ctx context.Context, req *pb.CreatePrefixRequest) (*emptypb.Empty, error) {
	err := g.serialMgr.CreateSerial(req.Prefix, req.StartNumber)
	switch err {
	case nil:
		return &emptypb.Empty{}, nil
	case model.Err_PrefixExist:
		return nil, status.Error(codes.AlreadyExists, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}

// 更新序號
func (g *gcp) UpdateStartNumber(ctx context.Context, req *pb.CreatePrefixRequest) (*emptypb.Empty, error) {
	err := g.serialMgr.UpdateSerial(req.Prefix, req.StartNumber)
	switch err {
	case nil:
		return &emptypb.Empty{}, nil
	case model.Err_PrefixNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}

// 取得序號
func (g *gcp) GetSerialNumber(ctx context.Context, req *pb.GetSerialNumberRequest) (*pb.SerialNumberRespose, error) {
	serial, err := g.serialMgr.GetSerial(req.Prefix)
	switch err {
	case nil:
		return &pb.SerialNumberRespose{
			Prefix:       req.Prefix,
			SerialNumber: serial,
		}, nil
	case model.Err_PrefixNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}

// 清除前綴
func (g *gcp) ClearPrefix(ctx context.Context, req *pb.GetSerialNumberRequest) (*emptypb.Empty, error) {
	err := g.serialMgr.ClearSerial(req.Prefix)
	switch err {
	case nil:
		return &emptypb.Empty{}, nil
	case model.Err_PrefixNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}
