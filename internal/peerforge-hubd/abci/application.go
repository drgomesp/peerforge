package abci

import (
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
	return types.ResponseCheckTx{
		Code: types.CodeTypeOK,
	}
}

func (a Application) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	return types.ResponseDeliverTx{Code: types.CodeTypeOK}
}
