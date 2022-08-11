package abci

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/tendermint/tendermint/abci/types"
)

var _ types.Application = Application{}

type Application struct {
	*types.BaseApplication
}

func NewApplication() *Application {
	return &Application{
		types.NewBaseApplication(),
	}
}

func (a Application) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	tx := req.GetTx()

	_ = tx
	return types.ResponseCheckTx{
		Code: types.CodeTypeOK,
	}
}

func (a Application) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	tx := req.GetTx()

	spew.Dump(string(tx))
	return types.ResponseDeliverTx{Code: types.CodeTypeOK, Data: tx}
}
