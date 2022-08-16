package abci

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
	peerforge "github.com/drgomesp/peerforge/pkg"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/abci/types"
)

const (
	_ = iota
	ErrCodeInvalidTx
)

var _ types.Application = &Application{}

type Application struct {
	*types.BaseApplication
	db           *badger.DB
	pendingBlock *badger.Txn
}

func NewApplication(db *badger.DB) *Application {
	return &Application{
		BaseApplication: types.NewBaseApplication(),
		db:              db,
	}
}

func (a *Application) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	a.pendingBlock = a.db.NewTransaction(true)
	return types.ResponseBeginBlock{}
}

func (a *Application) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
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

func (a *Application) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
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

	for _, event := range tx.Events {
		if err := a.pendingBlock.Set([]byte(event.ID), []byte(data)); err != nil {
			log.Err(err).Msgf("Error reading database, unable to verify tx")
		}

	}

	return types.ResponseDeliverTx{Code: types.CodeTypeOK, Data: req.GetTx()}
}

func (a *Application) Commit() types.ResponseCommit {
	if err := a.pendingBlock.Commit(); err != nil {
		log.Err(err).Msgf("Error writing to database, unable to commit block")
	}

	return types.ResponseCommit{Data: []byte{}}
}

func (a *Application) Query(req types.RequestQuery) types.ResponseQuery {
	resp := types.ResponseQuery{Key: req.Data}

	dbErr := a.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(req.Data)
		if err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
			resp.Log = "key does not exist"
			return nil
		}

		return item.Value(func(val []byte) error {
			resp.Log = "exists"
			resp.Value = val
			return nil
		})
	})

	if dbErr != nil {
		log.Err(dbErr).Msgf("Error reading database, unable to verify tx")
	}

	return resp
}
