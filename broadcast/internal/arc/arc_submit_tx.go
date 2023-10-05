package arc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/GorillaPool/go-junglebus"
	"github.com/libsv/go-bt/v2"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	arc_utils "github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/utils"
	"github.com/bitcoin-sv/go-broadcast-client/httpclient"
)

type SubmitTxRequest struct {
	RawTx string `json:"rawTx"`
}

var ErrSubmitTxMarshal = errors.New("error while marshalling submit tx body")

func (a *ArcClient) SubmitTransaction(ctx context.Context, tx *broadcast.Transaction, opts ...broadcast.TransactionOptFunc) (*broadcast.SubmitTxResponse, error) {
	if a == nil {
		return nil, broadcast.ErrClientUndefined
	}

	options := &broadcast.TransactionOpts{}
	for _, opt := range opts {
		opt(options)
	}

	result, err := submitTransaction(ctx, a, tx, options)
	if err != nil {
		return nil, err
	}

	if err := validateSubmitTxResponse(result); err != nil {
		return nil, err
	}

	response := &broadcast.SubmitTxResponse{
		BaseResponse: broadcast.BaseResponse{Miner: a.apiURL},
		SubmittedTx:  result,
	}

	return response, nil
}

func (a *ArcClient) SubmitBatchTransactions(ctx context.Context, txs []*broadcast.Transaction, opts ...broadcast.TransactionOptFunc) (*broadcast.SubmitBatchTxResponse, error) {
	if a == nil {
		return nil, broadcast.ErrClientUndefined
	}

	if len(txs) == 0 {
		return nil, errors.New("invalid request, no transactions to submit")
	}

	options := &broadcast.TransactionOpts{}
	for _, opt := range opts {
		opt(options)
	}

	result, err := submitBatchTransactions(ctx, a, txs, options)
	if err != nil {
		return nil, err
	}

	if err := validateBatchResponse(result); err != nil {
		return nil, err
	}

	response := broadcast.SubmitBatchTxResponse{
		BaseResponse: broadcast.BaseResponse{Miner: a.apiURL},
		Transactions: result,
	}

	return &response, nil
}

func submitTransaction(ctx context.Context, arc *ArcClient, tx *broadcast.Transaction, opts *broadcast.TransactionOpts) (*broadcast.SubmittedTx, error) {
	url := arc.apiURL + arcSubmitTxRoute
	data, err := createSubmitTxBody(tx, opts.TransactionFormat)
	if err != nil {
		return nil, err
	}

	pld := httpclient.NewPayload(
		httpclient.POST,
		url,
		arc.token,
		data,
	)
	appendSubmitTxHeaders(&pld, opts)

	resp, err := arc.HTTPClient.DoRequest(
		ctx,
		pld,
	)
	if err != nil {
		return nil, arc_utils.HandleHttpError(err)
	}

	model := broadcast.SubmittedTx{}
	err = arc_utils.DecodeResponseBody(resp.Body, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func submitBatchTransactions(ctx context.Context, arc *ArcClient, txs []*broadcast.Transaction, opts *broadcast.TransactionOpts) ([]*broadcast.SubmittedTx, error) {
	url := arc.apiURL + arcSubmitBatchTxsRoute
	data, err := createSubmitBatchTxsBody(txs, opts.TransactionFormat)
	if err != nil {
		return nil, err
	}

	pld := httpclient.NewPayload(
		httpclient.POST,
		url,
		arc.token,
		data,
	)
	appendSubmitTxHeaders(&pld, opts)

	resp, err := arc.HTTPClient.DoRequest(
		ctx,
		pld,
	)
	if err != nil {
		return nil, arc_utils.HandleHttpError(err)
	}

	var model []*broadcast.SubmittedTx
	err = arc_utils.DecodeResponseBody(resp.Body, &model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func createSubmitTxBody(tx *broadcast.Transaction, txFormat broadcast.TransactionFormat) ([]byte, error) {
	body := &SubmitTxRequest{tx.Hex}

	if txFormat == broadcast.RawTxFormat {
		if err := convertToEfTransaction(body); err != nil {
			return nil, fmt.Errorf("Conversion to EF format failed: %s", err.Error())
		}
	} else if txFormat == broadcast.BeefFormat {
		// To be implemented
		return nil, fmt.Errorf("Submitting transactions in BEEF format is unimplemented yet...")
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, ErrSubmitTxMarshal
	}

	return data, nil
}

func createSubmitBatchTxsBody(txs []*broadcast.Transaction, txFormat broadcast.TransactionFormat) ([]byte, error) {
	rawTxs := make([]*SubmitTxRequest, 0, len(txs))
	for _, tx := range txs {
		rawTxs = append(rawTxs, &SubmitTxRequest{RawTx: tx.Hex})
	}

	if txFormat == broadcast.RawTxFormat {
		if err := convertBatchToEfTransaction(rawTxs); err != nil {
			return nil, fmt.Errorf("Conversion to EF format failed for one or more transactions with error: %s", err.Error())
		}
	} else if txFormat == broadcast.BeefFormat {
		// To be implemented
		return nil, fmt.Errorf("Submitting transactions in BEEF format is unimplemented yet...")
	}

	data, err := json.Marshal(rawTxs)
	if err != nil {
		return nil, ErrSubmitTxMarshal
	}

	return data, nil
}

func appendSubmitTxHeaders(pld *httpclient.HTTPRequest, opts *broadcast.TransactionOpts) {
	if opts == nil {
		return
	}

	if opts.MerkleProof {
		pld.AddHeader("X-MerkleProof", "true")
	}

	if opts.CallbackURL != "" {
		pld.AddHeader("X-CallbackUrl", opts.CallbackURL)
	}

	if opts.CallbackToken != "" {
		pld.AddHeader("X-CallbackToken", opts.CallbackToken)
	}

	if statusCode, ok := broadcast.MapTxStatusToInt(opts.WaitForStatus); ok {
		pld.AddHeader("X-WaitForStatus", strconv.Itoa(statusCode))
	}
}

func validateBatchResponse(model []*broadcast.SubmittedTx) error {
	for _, tx := range model {
		if err := validateSubmitTxResponse(tx); err != nil {
			return err
		}
	}

	return nil
}

func validateSubmitTxResponse(model *broadcast.SubmittedTx) error {
	if model.TxStatus == "" {
		return broadcast.ErrMissingStatus
	}

	return nil
}

func convertToEfTransaction(tx *SubmitTxRequest) error {
	junglebusClient, err := junglebus.New(
		junglebus.WithHTTP("https://junglebus.gorillapool.io"),
	)
	if err != nil {
		return err
	}

	transaction, err := bt.NewTxFromString(tx.RawTx)
	if err != nil {
		return err
	}

	for _, input := range transaction.Inputs {
		if err = updateUtxoWithMissingData(junglebusClient, input); err != nil {
			return err
		}
	}

	tx.RawTx = hex.EncodeToString(transaction.ExtendedBytes())
	return nil
}

func updateUtxoWithMissingData(jbc *junglebus.Client, input *bt.Input) error {
	txid := input.PreviousTxIDStr()

	tx, err := jbc.GetTransaction(context.Background(), txid)
	if err != nil {
		return err
	}

	actualTx, err := bt.NewTxFromBytes(tx.Transaction)
	if err != nil {
		return err
	}

	o := actualTx.Outputs[input.PreviousTxOutIndex]
	input.PreviousTxScript = o.LockingScript
	input.PreviousTxSatoshis = o.Satoshis
	return nil
}

func convertBatchToEfTransaction(rawTxs []*SubmitTxRequest) error {
	for _, rawTx := range rawTxs {
		if err := convertToEfTransaction(rawTx); err != nil {
			return err
		}
	}
	return nil
}
