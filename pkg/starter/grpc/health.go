package grpc

import (
	"context"
	pb "google.golang.org/grpc/health/grpc_health_v1"
	"hidevops.io/hiboot/pkg/at"
)

// controller
type healthCheckService struct {
	at.HealthCheckService
	// declare HelloServiceClient
	healthClient pb.HealthClient
}

// Init inject helloServiceClient
func newHealthCheckService(healthClient pb.HealthClient) *healthCheckService {
	return &healthCheckService{
		healthClient: healthClient,
	}
}

// Status return health check display name grpc
func (c *healthCheckService) Name() (name string) {
	return Profile
}

// Status return grpc health check status as bool
func (c *healthCheckService) Status() (up bool) {
	resp, err := c.healthClient.Check(context.TODO(), &pb.HealthCheckRequest{})
	if err == nil {
		up = resp.Status == pb.HealthCheckResponse_SERVING
	}
	return
}
