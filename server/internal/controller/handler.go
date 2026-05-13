package controller

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "dns-manager/proto/dns"
	"dns-manager/server/internal/domain"
	"dns-manager/server/internal/usecase"
)

type Handler struct {
	pb.UnimplementedDnsServiceServer
	uc *usecase.DnsUseCase
}

func NewHandler(uc *usecase.DnsUseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) AddDnsServer(ctx context.Context, req *pb.AddDnsServerRequest) (*pb.DnsResponse, error) {
	if err := h.uc.Add(ctx, req.Address); err != nil {
		return nil, mapError(err)
	}
	return &pb.DnsResponse{Success: true}, nil
}

func (h *Handler) RemoveDnsServer(ctx context.Context, req *pb.RemoveDnsServerRequest) (*pb.DnsResponse, error) {
	if err := h.uc.Remove(ctx, req.Address); err != nil {
		return nil, mapError(err)
	}
	return &pb.DnsResponse{Success: true}, nil
}

func (h *Handler) ListDnsServers(ctx context.Context, req *pb.ListDnsServersRequest) (*pb.ListDnsServersResponse, error) {
	servers, err := h.uc.List(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	return &pb.ListDnsServersResponse{Success: true, Servers: servers}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, "dns server already exists")
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, "dns server not found")
	case errors.Is(err, domain.ErrInvalidAddr):
		return status.Error(codes.InvalidArgument, "invalid address")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
