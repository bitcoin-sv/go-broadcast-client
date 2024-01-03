package arc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
		a.Logger.Error().Msgf("Failed to submit tx: %s", broadcast.ErrClientUndefined.Error())
		return nil, broadcast.ErrClientUndefined
	}

	options := &broadcast.TransactionOpts{}
	for _, opt := range opts {
		opt(options)
	}

	result, err := submitTransaction(ctx, a, tx, options)
	if err != nil {
		a.Logger.Error().Msgf("Failed to submit tx: %s", err.Error())
		return nil, err
	}

	if err := validateSubmitTxResponse(result); err != nil {
		a.Logger.Error().Msgf("Failed to validate submit tx response: %s", err.Error())
		return nil, err
	}

	response := &broadcast.SubmitTxResponse{
		BaseResponse: broadcast.BaseResponse{Miner: a.apiURL},
		SubmittedTx:  result,
	}

	a.Logger.Debug().Msgf("Got submit tx response from miner: %s", response.Miner)

	return response, nil
}

func (a *ArcClient) SubmitBatchTransactions(ctx context.Context, txs []*broadcast.Transaction, opts ...broadcast.TransactionOptFunc) (*broadcast.SubmitBatchTxResponse, error) {
	if a == nil {
		a.Logger.Error().Msgf("Failed to submit batch txs: %s", broadcast.ErrClientUndefined.Error())
		return nil, broadcast.ErrClientUndefined
	}

	if len(txs) == 0 {
		err := errors.New("invalid request, no transactions to submit")
		a.Logger.Error().Msgf("Failed to submit batch txs: %s", err.Error())
		return nil, err
	}

	options := &broadcast.TransactionOpts{}
	for _, opt := range opts {
		opt(options)
	}

	result, err := submitBatchTransactions(ctx, a, txs, options)
	if err != nil {
		a.Logger.Error().Msgf("Failed to submit batch txs: %s", err.Error())
		return nil, err
	}

	if err := validateBatchResponse(result); err != nil {
		a.Logger.Error().Msgf("Failed to validate batch txs response: %s", err.Error())
		return nil, err
	}

	response := &broadcast.SubmitBatchTxResponse{
		BaseResponse: broadcast.BaseResponse{Miner: a.apiURL},
		Transactions: result,
	}

	a.Logger.Debug().Msgf("Got submit batch txs response from miner: %s", response.Miner)
	return response, nil
}

func submitTransaction(ctx context.Context, arc *ArcClient, tx *broadcast.Transaction, opts *broadcast.TransactionOpts) (*broadcast.SubmittedTx, error) {
	url := arc.apiURL + arcSubmitTxRoute
	data, err := createSubmitTxBody(arc, tx, opts.TransactionFormat)
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

	model, err := decodeSubmitResponseBody(resp)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func submitBatchTransactions(ctx context.Context, arc *ArcClient, txs []*broadcast.Transaction, opts *broadcast.TransactionOpts) ([]*broadcast.SubmittedTx, error) {
	url := arc.apiURL + arcSubmitBatchTxsRoute
	data, err := createSubmitBatchTxsBody(arc, txs, opts.TransactionFormat)
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

	model, err := decodeSubmitBatchResponseBody(resp)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func createSubmitTxBody(arc *ArcClient, tx *broadcast.Transaction, txFormat broadcast.TransactionFormat) ([]byte, error) {
	body := &SubmitTxRequest{
		RawTx: tx.Hex,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, ErrSubmitTxMarshal
	}

	return data, nil
}

func createSubmitBatchTxsBody(arc *ArcClient, txs []*broadcast.Transaction, txFormat broadcast.TransactionFormat) ([]byte, error) {
	rawTxs := make([]*SubmitTxRequest, 0, len(txs))
	for _, tx := range txs {
		requestTx := &SubmitTxRequest{
			RawTx: tx.Hex,
		}
		rawTxs = append(rawTxs, requestTx)
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

func decodeSubmitResponseBody(resp *http.Response) (*broadcast.SubmittedTx, error) {
	base := broadcast.BaseSubmitTxResponse{}
	err := arc_utils.DecodeResponseBody(resp.Body, &base)
	if err != nil {
		return nil, broadcast.ErrUnableToDecodeMerklePath
	}

	model := &broadcast.SubmittedTx{
		BaseSubmitTxResponse: base,
		MerklePath:           base.MerklePath,
	}

	return model, nil
}

func decodeSubmitBatchResponseBody(resp *http.Response) ([]*broadcast.SubmittedTx, error) {
	base := make([]broadcast.BaseSubmitTxResponse, 0)
	err := arc_utils.DecodeResponseBody(resp.Body, &base)
	if err != nil {
		return nil, err
	}

	model := make([]*broadcast.SubmittedTx, 0)
	for _, tx := range base {
		model = append(model, &broadcast.SubmittedTx{
			BaseSubmitTxResponse: tx,
			MerklePath:           tx.MerklePath,
		})
	}

	return model, nil
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
