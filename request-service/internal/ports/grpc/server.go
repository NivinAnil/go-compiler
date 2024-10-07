package grpc

import (
	"context"
	"go-compiler/common/pkg/utils/logger"
	"go-compiler/request-service/internal/domain/dto/request"
	"go-compiler/request-service/internal/domain/services/interfaces"
	pb "go-compiler/request-service/proto"
)

type GRPCServer struct {
	pb.UnimplementedRequestServiceServer
	ExecutionService interfaces.IExecutionService
}

func NewGRPCServer(es interfaces.IExecutionService) *GRPCServer {
	return &GRPCServer{
		ExecutionService: es,
	}
}

func (s *GRPCServer) ProcessRequest(ctx context.Context, req *pb.ExecutionRequest) (*pb.ExecutionResponse, error) {
	log := logger.GetLogger(ctx)
	methodName := "ProcessRequest"
	log.Info("Entering", "methodName", methodName)

	payload := request.NewExecutionRequest{
		Id:          req.Id,
		Code:        req.Code,
		LanguageId:  req.LanguageId,
		ConnctionId: req.ConnectionId,
		StdIn:       req.Stdin,
	}

	err := s.ExecutionService.ProcessRequest(ctx, payload)
	if err != nil {
		log.Error("Error in processing request", "error", err.Error())
		return nil, err
	}

	return &pb.ExecutionResponse{
		Message: "Request processed successfully",
	}, nil
}