package arc

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	arc_utils "github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/utils"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
)

// SubmitTxRequest is the request body for the submit tx endpoint.
type SubmitTxRequest struct {
	RawTx string `json:"rawTx"`
}

// ErrSubmitTxMarshal is the error returned when the submit tx request body cannot be marshalled.
var ErrSubmitTxMarshal = errors.New("error while marshalling submit tx body")

// SumitTransaction will fire the HTTP request to submit the tx to the network.
func (a *ArcClient) SubmitTransaction(ctx context.Context, tx *broadcast.Transaction) (*broadcast.SubmitTxResponse, error) {
	if a == nil {
		return nil, broadcast.ErrClientUndefined
	}

	result, err := submitTransaction(ctx, a, tx)
	if err != nil {
		return nil, err
	}

	if err := validateSubmitTxResponse(result); err != nil {
		return nil, err
	}

	return result, nil
}

// SubmitBatchTransactions will fire the HTTP request to submit the txs to the network.
func (a *ArcClient) SubmitBatchTransactions(ctx context.Context, txs []*broadcast.Transaction) ([]*broadcast.SubmitTxResponse, error) {
	if a == nil {
		return nil, broadcast.ErrClientUndefined
	}

	if len(txs) == 0 {
		return nil, errors.New("invalid request, no transactions to submit")
	}

	result, err := submitBatchTransactions(ctx, a, txs)
	if err != nil {
		return nil, err
	}

	if err := validateBatchResponse(result); err != nil {
		return nil, err
	}

	return result, nil
}

func submitTransaction(ctx context.Context, arc *ArcClient, tx *broadcast.Transaction) (*broadcast.SubmitTxResponse, error) {
	url := arc.apiURL + arcSubmitTxRoute
	data, err := createSubmitTxBody(tx)
	if err != nil {
		return nil, err
	}

	pld := httpclient.NewPayload(
		httpclient.POST,
		url,
		arc.token,
		data,
	)
	appendSubmitTxHeaders(&pld, tx)

	resp, err := arc.HTTPClient.DoRequest(
		ctx,
		pld,
	)
	if err != nil {
		return nil, arc_utils.HandleHttpError(err)
	}

	model := broadcast.SubmitTxResponse{}
	err = arc_utils.DecodeResponseBody(resp.Body, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func submitBatchTransactions(ctx context.Context, arc *ArcClient, txs []*broadcast.Transaction) ([]*broadcast.SubmitTxResponse, error) {
	url := arc.apiURL + arcSubmitBatchTxsRoute
	data, err := createSubmitBatchTxsBody(txs)
	if err != nil {
		return nil, err
	}

	pld := httpclient.NewPayload(
		httpclient.POST,
		url,
		arc.token,
		data,
	)
	appendSubmitTxHeaders(&pld, txs[0])

	resp, err := arc.HTTPClient.DoRequest(
		ctx,
		pld,
	)
	if err != nil {
		return nil, arc_utils.HandleHttpError(err)
	}

	var model []*broadcast.SubmitTxResponse
	err = arc_utils.DecodeResponseBody(resp.Body, &model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func createSubmitTxBody(tx *broadcast.Transaction) ([]byte, error) {
	body := &SubmitTxRequest{tx.RawTx}
	data, err := json.Marshal(body)

	if err != nil {
		return nil, ErrSubmitTxMarshal
	}

	return data, nil
}

func createSubmitBatchTxsBody(txs []*broadcast.Transaction) ([]byte, error) {
	rawTxs := make([]*SubmitTxRequest, 0, len(txs))
	for _, tx := range txs {
		rawTxs = append(rawTxs, &SubmitTxRequest{RawTx: tx.RawTx})
	}

	data, err := json.Marshal(rawTxs)
	if err != nil {
		return nil, ErrSubmitTxMarshal
	}

	return data, nil
}

func appendSubmitTxHeaders(pld *httpclient.HTTPRequest, tx *broadcast.Transaction) {
	if tx.MerkleProof {
		pld.AddHeader("X-MerkleProof", "true")
	}

	if tx.CallBackURL != "" {
		pld.AddHeader("X-CallbackUrl", tx.CallBackURL)
	}

	if tx.CallBackToken != "" {
		pld.AddHeader("X-CallbackToken", tx.CallBackToken)
	}

	if statusCode, ok := broadcast.MapTxStatusToInt(tx.WaitForStatus); ok {
		pld.AddHeader("X-WaitForStatus", strconv.Itoa(statusCode))
	}
}

func validateBatchResponse(model []*broadcast.SubmitTxResponse) error {
	for _, tx := range model {
		if err := validateSubmitTxResponse(tx); err != nil {
			return err
		}
	}

	return nil
}

func validateSubmitTxResponse(model *broadcast.SubmitTxResponse) error {
	if model.TxStatus == "" {
		return broadcast.ErrMissingStatus
	}

	return nil
}
