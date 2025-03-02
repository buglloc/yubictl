package yubictl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type Yubikey struct {
	id        string
	serial    uint32
	pingTick  time.Duration
	httpc     *resty.Client
	ctx       context.Context
	cancelCtx context.CancelFunc
	closed    chan struct{}
}

func (y *Yubikey) ID() string {
	return y.id
}

func (y *Yubikey) Serial() uint32 {
	return y.serial
}

func (y *Yubikey) Touch(ctx context.Context) error {
	var serviceErr ServiceError
	rsp, err := y.httpc.R().
		SetContext(ctx).
		SetError(&serviceErr).
		SetBody(TouchReq{
			ID: y.id,
		}).
		ForceContentType("application/json").
		Get("/v1/touch")

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if !rsp.IsSuccess() {
		if serviceErr.Code != ServiceErrorCodeNone {
			return &serviceErr
		}

		return fmt.Errorf("request failed: non-200 status code: %s", rsp.Status())
	}

	return nil
}

func (y *Yubikey) Reboot(ctx context.Context) error {
	var serviceErr ServiceError
	rsp, err := y.httpc.R().
		SetContext(ctx).
		SetError(&serviceErr).
		SetBody(RebootReq{
			ID: y.id,
		}).
		ForceContentType("application/json").
		Get("/v1/reboot")

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if !rsp.IsSuccess() {
		if serviceErr.Code != ServiceErrorCodeNone {
			return &serviceErr
		}

		return fmt.Errorf("request failed: non-200 status code: %s", rsp.Status())
	}

	return nil
}

func (y *Yubikey) Ping(ctx context.Context) error {
	var serviceErr ServiceError
	rsp, err := y.httpc.R().
		SetContext(ctx).
		SetError(&serviceErr).
		SetBody(PingReq{
			ID: y.id,
		}).
		ForceContentType("application/json").
		Get("/v1/ping")

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if !rsp.IsSuccess() {
		if serviceErr.Code != ServiceErrorCodeNone {
			return &serviceErr
		}

		return fmt.Errorf("request failed: non-200 status code: %s", rsp.Status())
	}

	return nil
}

func (y *Yubikey) Close(ctx context.Context) error {
	return y.Release(ctx)
}

func (y *Yubikey) Release(ctx context.Context) error {
	y.cancelCtx()

	var serviceErr ServiceError
	rsp, err := y.httpc.R().
		SetContext(ctx).
		SetError(&serviceErr).
		SetBody(RebootReq{
			ID: y.id,
		}).
		ForceContentType("application/json").
		Get("/v1/release")

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if !rsp.IsSuccess() {
		if serviceErr.Code != ServiceErrorCodeNone {
			return &serviceErr
		}

		return fmt.Errorf("request failed: non-200 status code: %s", rsp.Status())
	}

	return nil
}

func (y *Yubikey) pingLoop() {
	defer close(y.closed)

	for {
		tonextTick := time.Until(
			time.Now().Add(y.pingTick).Truncate(y.pingTick),
		)
		t := time.NewTimer(tonextTick)

		select {
		case <-y.ctx.Done():
			t.Stop()
			return
		case <-t.C:
			if err := y.Ping(y.ctx); err != nil {
				log.Error().
					Err(err).
					Str("session_id", y.id).
					Uint32("yk_serial", y.serial).
					Msg("yubikey ping failed")
			}
		}
	}
}
