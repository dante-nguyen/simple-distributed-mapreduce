package master

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/logx"
	rpcv1 "github.com/nlduy0310/simple-distributed-mapreduce/rpc/v1"
)

var (
	errInvalidConfig = errors.New("invalid config")
)

type Service struct {
	rpcv1.UnimplementedMasterServiceServer
	Config Config
	reg    *registry
}

func (s *Service) Register(ctx context.Context, req *rpcv1.RegisterRequest) (*rpcv1.RegisterResponse, error) {
	if err := s.reg.register(ctx, req.Name, req.Address); err != nil {
		logx.Err(fmt.Sprintf("register worker %q at %q", req.Name, req.Address), err)
		return nil, err
	}

	log.Printf("successfully registered worker %q at %q", req.Name, req.Address)
	return &rpcv1.RegisterResponse{Ok: true}, nil
}

func (s *Service) Heartbeat(ctx context.Context, req *rpcv1.HeartbeatRequest) (*rpcv1.HeartbeatResponse, error) {
	ts := time.Now()
	if err := s.reg.recordHeartbeat(ctx, req.Name, ts); err != nil {
		return nil, err
	}

	// log.Printf("received heartbeat from %q at %s", req.Name, ts)
	return &rpcv1.HeartbeatResponse{Ok: true}, nil
}

func NewService(cfg Config) (*Service, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, errx.Chain(errInvalidConfig, err)
	}

	return &Service{
		Config: cfg,
		reg:    newRegistry(),
	}, nil
}

// PeriodicHealthcheck periodically check latest worker heartbeats and remove them if necessary
func (s *Service) PeriodicHealthcheck(ctx context.Context, interval, timeout, healthy time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.doHealthcheck(ctx, time.Now(), timeout, healthy)
			if err == nil {
				continue
			} else if !errx.OneOf(err, context.Canceled, context.DeadlineExceeded) {
				logx.Err("healthcheck", err)
			} else if ctxErr := ctx.Err(); ctxErr != nil { // parent stopped
				return ctxErr
			}
			// else it's a contention timeout, so retry
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Service) doHealthcheck(parent context.Context, now time.Time, timeout, healthy time.Duration) error {
	ctx, timeoutHealthcheck := context.WithTimeout(parent, timeout)
	defer timeoutHealthcheck()

	names, err := s.reg.names(ctx)
	if err != nil {
		return err
	}

	for _, name := range names {
		lastHeartbeat, err := s.reg.lastHeartbeat(ctx, name)
		if err != nil {
			return err
		}

		diff := now.Sub(lastHeartbeat)
		if diff <= healthy {
			continue
		}

		if err = s.reg.remove(ctx, name); err != nil {
			return err
		}
	}

	return nil
}
