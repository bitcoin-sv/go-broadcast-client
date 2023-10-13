package arc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

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

	model := broadcast.SubmittedTx{}
	err = arc_utils.DecodeResponseBody(resp.Body, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
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

	var model []*broadcast.SubmittedTx
	err = arc_utils.DecodeResponseBody(resp.Body, &model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func createSubmitTxBody(arc *ArcClient, tx *broadcast.Transaction, txFormat broadcast.TransactionFormat) ([]byte, error) {
	body, err := formatTxRequest(arc, tx, txFormat)
	if err != nil {
		return nil, err
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
		requestTx, err := formatTxRequest(arc, tx, txFormat)
		if err != nil {
			return nil, err
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

func formatTxRequest(arc *ArcClient, tx *broadcast.Transaction, txFormat broadcast.TransactionFormat) (*SubmitTxRequest, error) {
	var (
		body *SubmitTxRequest
		err  error
	)

	switch txFormat {
	case broadcast.RawTxFormat:
		body, err = rawTxRequest(arc, tx.Hex)
	case broadcast.EfFormat:
		body, err = efTxRequest(tx.Hex)
	case broadcast.BeefFormat:
		body, err = beefTxRequest(tx.Hex)
	default:
		err = fmt.Errorf("unknown transaction format")
	}

	if err != nil {
		return nil, err
	}

	return body, nil
}

func efTxRequest(rawTx string) (*SubmitTxRequest, error) {
	request := &SubmitTxRequest{RawTx: rawTx}

	return request, nil
}

func beefTxRequest(rawTx string) (*SubmitTxRequest, error) {
	return nil, fmt.Errorf("Submitting transactions in BEEF format is unimplemented yet...")
}

func rawTxRequest(arc *ArcClient, rawTx string) (*SubmitTxRequest, error) {
	transaction, err := bt.NewTxFromString(rawTx)
	if err != nil {
		return nil, err
	}

	for _, input := range transaction.Inputs {
		if err = updateUtxoWithMissingData(arc, input); err != nil {
			return nil, err
		}
	}

	request := &SubmitTxRequest{
		RawTx: hex.EncodeToString(transaction.ExtendedBytes()),
	}
	return request, nil
}

func updateUtxoWithMissingData(arc *ArcClient, input *bt.Input) error {
	txid := input.PreviousTxIDStr()

	pld := httpclient.NewPayload(
		httpclient.GET,
		fmt.Sprintf("https://junglebus.gorillapool.io/v1/transaction/get/%s", txid),
		"",
		nil,
	)
	resp, err := arc.HTTPClient.DoRequest(context.Background(), pld)
	if err != nil {
		return err
	}

	var tx *junglebusTransaction
	err = arc_utils.DecodeResponseBody(resp.Body, &tx)
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

type junglebusTransaction struct {
	ID          string `json:"id"`
	Transaction []byte `json:"transaction"`
	BlockHash   string `json:"block_hash"`
	BlockHeight uint32 `json:"block_height"`
	BlockTime   uint32 `json:"block_time"`
	BlockIndex  uint64 `json:"block_index"`

	// index data
	// input/output types are
	// p2pkh, p2sh, token-stas, opreturn, tokenized, metanet, bitcom, run, map, bap, non-standard etc.
	Addresses   []string `json:"addresses"`
	Inputs      []string `json:"inputs"`
	Outputs     []string `json:"outputs"`
	InputTypes  []string `json:"input_types"`
	OutputTypes []string `json:"output_types"`
	Contexts    []string `json:"contexts"`     // optional contexts of output types, only for known protocols
	SubContexts []string `json:"sub_contexts"` // optional sub-contexts of output types, only for known protocols
	Data        []string `json:"data"`         // optional data of output types, only for known protocols

	// the merkle proof in binary
	MerkleProof []byte `json:"merkle_proof"`
}
