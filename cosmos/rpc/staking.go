package rpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"pkg.berachain.dev/polaris/eth/polar"
)

type StakingAPIHandler struct {
	backend polar.Backend
	qc      func(height int64, prove bool) (sdk.Context, error)
	qs      stakingtypes.QueryServer
}

func NewStakingAPIHandler(be polar.Backend, qc func(height int64, prove bool) (sdk.Context, error), qs stakingtypes.QueryServer) *StakingAPIHandler {
	return &StakingAPIHandler{
		backend: be,
		qc:      qc,
		qs:      qs,
	}
}

func (h StakingAPIHandler) Validators(page string) (*stakingtypes.QueryValidatorsResponse, error) {
	blockNum := h.backend.CurrentBlock().Number.Int64()
	ctx, err := h.qc(blockNum, false)
	if err != nil {
		return nil, err
	}
	return h.qs.Validators(ctx, &stakingtypes.QueryValidatorsRequest{})
}
