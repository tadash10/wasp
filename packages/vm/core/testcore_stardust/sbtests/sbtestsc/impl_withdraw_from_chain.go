package sbtestsc

import (
	"math"

	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/kv/kvdecoder"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
)

// withdrawFromChain withdraws all the available balance existing on the target chain
func withdrawFromChain(ctx iscp.Sandbox) dict.Dict {
	ctx.Log().Infof(FuncWithdrawFromChain.Name)
	params := kvdecoder.New(ctx.Params(), ctx.Log())
	targetChain := params.MustGetChainID(ParamChainID)
	iotasToWithdrawal := params.MustGetUint64(ParamIotasToWithdrawal)
	// gasBudget := params.MustGetUint64(ParamGasBudgetToSend)

	availableIotas := ctx.AllowanceAvailable().Iotas

	request := iscp.RequestParameters{
		TargetAddress: targetChain.AsAddress(),
		Assets:        iscp.NewAssetsIotas(availableIotas),
		Metadata: &iscp.SendMetadata{
			TargetContract: accounts.Contract.Hname(),
			EntryPoint:     accounts.FuncWithdraw.Hname(),
			// GasBudget:      gasBudget,
			GasBudget: math.MaxUint64,
			Allowance: iscp.NewAssetsIotas(iotasToWithdrawal),
		},
	}
	requiredDustDeposit := ctx.EstimateRequiredDustDeposit(request)
	if availableIotas < requiredDustDeposit {
		ctx.Log().Panicf("no enough iotas sent to cover dust deposit")
	}
	ctx.TransferAllowedFunds(ctx.AccountID())
	ctx.Send(request)

	ctx.Log().Infof("%s: success", FuncWithdrawFromChain)
	return nil
}
