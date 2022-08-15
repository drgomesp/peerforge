package abci

import (
	"encoding/json"

	peerforge "github.com/drgomesp/peerforge/pkg"
	"github.com/tendermint/tendermint/abci/types"
)

const (
	_ = iota
	ErrCodeInvalidTx
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
	var tx struct {
		Events []*peerforge.Event `json:"events"`
	}
	if err := json.Unmarshal(req.GetTx(), &tx); err != nil {
		return types.ResponseCheckTx{
			Code: ErrCodeInvalidTx,
			Info: err.Error(),
		}
	}

	return types.ResponseCheckTx{
		Code: types.CodeTypeOK,
	}
}

func (a Application) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	data := string(req.GetTx())
	var tx struct {
		Events []*peerforge.Event `json:"events"`
	}
	if err := json.Unmarshal(req.GetTx(), &tx); err != nil {
		return types.ResponseDeliverTx{
			Code: ErrCodeInvalidTx,
			Info: err.Error(),
		}
	}

	_ = data

	return types.ResponseDeliverTx{Code: types.CodeTypeOK, Data: req.GetTx()}
}
