package grpc

import (
	libgrpc "github.com/ElfAstAhe/go-service-template/pkg/transport/grpc"
)

func mapToGrpcError(err error) error {
	return libgrpc.MapToGrpcError(err)
}
