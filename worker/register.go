package worker

import (
	"context"
	"time"

	"github.com/nlduy0310/simple-distributed-mapreduce/errorsx"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

func (s *Server) initialRegister() {
	s.logger.Debugf("attempting to register this as %s", s.id)

	timeoutDuration := time.Duration(s.cfg.RegisterTimeoutSec) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	req := constructRegisterPayload(s)
	_, err := s.masterServiceClient.RegisterWorker(ctx, req)
	if err != nil {
		s.logger.Fatal(errorsx.WrapAsMessage("unable to register to master", err))
	}
}

func constructRegisterPayload(s *Server) *rpcv1.RegisterWorkerRequest {
	return &rpcv1.RegisterWorkerRequest{Id: s.id}
}
