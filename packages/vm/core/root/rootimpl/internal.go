package rootimpl

import (
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/assert"
	"github.com/iotaledger/wasp/packages/iscp/coreutil"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/core/root"
)

func mustStoreContract(ctx iscp.Sandbox, i *coreutil.ContractInfo, a assert.Assert) {
	rec := root.NewContractRecord(i, &iscp.NilAgentID)
	ctx.Log().Debugf("mustStoreAndInitCoreContract: '%s', hname = %s", i.Name, i.Hname())
	mustStoreContractRecord(ctx, rec, a)
}

func mustStoreAndInitCoreContract(ctx iscp.Sandbox, i *coreutil.ContractInfo, a assert.Assert, params ...dict.Dict) {
	mustStoreContract(ctx, i, a)
	var p dict.Dict
	if len(params) == 1 {
		p = params[0]
	}
	_, err := ctx.Call(iscp.Hn(i.Name), iscp.EntryPointInit, p, nil)
	a.RequireNoError(err)
}

func mustStoreContractRecord(ctx iscp.Sandbox, rec *root.ContractRecord, a assert.Assert) {
	hname := rec.Hname()
	contractRegistry := collections.NewMap(ctx.State(), root.VarContractRegistry)
	a.Require(!contractRegistry.MustHasAt(hname.Bytes()), "contract '%s'/%s already exist", rec.Name, hname.String())
	contractRegistry.MustSetAt(hname.Bytes(), rec.Bytes())
}

// isAuthorizedToDeploy checks if caller is authorized to deploy smart contract
func isAuthorizedToDeploy(ctx iscp.Sandbox) bool {
	caller := ctx.Caller()
	if caller.Equals(ctx.ChainOwnerID()) {
		// chain owner is always authorized
		return true
	}
	if caller.Address().Equals(ctx.ChainID().AsAddress()) {
		// smart contract from the same chain is always authorize
		return true
	}

	return collections.NewMap(ctx.State(), root.VarDeployPermissions).MustHasAt(caller.Bytes())
}
