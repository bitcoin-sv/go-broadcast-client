package arc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	arc_utils "github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/utils"
	"github.com/bitcoin-sv/go-broadcast-client/httpclient"
	"github.com/libsv/go-bt/v2"
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

	sumbiter := getSubmiter(options)
	result, err := sumbiter.submit(ctx, a, tx, options)
	if err != nil {
		return nil, arc_utils.WithCause(errors.New("SubmitTransaction: submitting failed"), err)
	}

	if err := validateSubmitTxResponse(result); err != nil {
		return nil, arc_utils.WithCause(errors.New("SubmitTransaction: validation of submit tx response failed"), err)
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
		return nil, broadcast.ErrClientUndefined
	}

	if len(txs) == 0 {
		err := errors.New("invalid request, no transactions to submit")
		return nil, arc_utils.WithCause(errors.New("SubmitBatchTransactions: bad request"), err)
	}

	options := &broadcast.TransactionOpts{}
	for _, opt := range opts {
		opt(options)
	}

	sumbiter := getSubmiter(options)
	result, err := sumbiter.batchSubmit(ctx, a, txs, options)
	if err != nil {
		return nil, arc_utils.WithCause(errors.New("SubmitBatchTransactions: submitting failed"), err)
	}

	if err := validateBatchResponse(result); err != nil {
		return nil, arc_utils.WithCause(errors.New("SubmitBatchTransactions: validation of batch submit tx response failed"), err)
	}

	response := &broadcast.SubmitBatchTxResponse{
		BaseResponse: broadcast.BaseResponse{Miner: a.apiURL},
		Transactions: result,
	}

	a.Logger.Debug().Msgf("Got submit batch txs response from miner: %s", response.Miner)
	return response, nil
}

type submitFormatStrategy interface {
	submit(ctx context.Context, arc *ArcClient, tx *broadcast.Transaction, opts *broadcast.TransactionOpts) (*broadcast.SubmittedTx, error)
	batchSubmit(ctx context.Context, arc *ArcClient, txs []*broadcast.Transaction, opts *broadcast.TransactionOpts) ([]*broadcast.SubmittedTx, error)
}

func getSubmiter(opts *broadcast.TransactionOpts) submitFormatStrategy {
	switch opts.TransactionFormat {
	case broadcast.RawTxFormat:
		return &rawtxSubmitStrategy{}
	case broadcast.EfFormat:
		return &efSubmitStrategy{}
	case broadcast.BeefFormat:
		return &beefSubmitStrategy{}
	default:
		return &efSubmitStrategy{}
	}
}

type beefSubmitStrategy struct{}

func (s *beefSubmitStrategy) submit(ctx context.Context, arc *ArcClient, tx *broadcast.Transaction, opts *broadcast.TransactionOpts) (*broadcast.SubmittedTx, error) {
	return nil, errors.New("submitting transactions in BEEF format is unimplemented yet...")
}

func (s *beefSubmitStrategy) batchSubmit(ctx context.Context, arc *ArcClient, txs []*broadcast.Transaction, opts *broadcast.TransactionOpts) ([]*broadcast.SubmittedTx, error) {
	return nil, errors.New("submitting transactions in BEEF format is unimplemented yet...")
}

type efSubmitStrategy struct{}

func (s *efSubmitStrategy) submit(ctx context.Context, arc *ArcClient, tx *broadcast.Transaction, opts *broadcast.TransactionOpts) (*broadcast.SubmittedTx, error) {
	body := &SubmitTxRequest{RawTx: tx.Hex}
	return submit(ctx, arc, body, opts)
}

func (s *efSubmitStrategy) batchSubmit(ctx context.Context, arc *ArcClient, txs []*broadcast.Transaction, opts *broadcast.TransactionOpts) ([]*broadcast.SubmittedTx, error) {
	body := convertToSubmittedTxs(txs)
	return batchSubmit(ctx, arc, body, opts)
}

type rawtxSubmitStrategy struct{}

func (s *rawtxSubmitStrategy) submit(ctx context.Context, arc *ArcClient, tx *broadcast.Transaction, opts *broadcast.TransactionOpts) (*broadcast.SubmittedTx, error) {
	body := &SubmitTxRequest{RawTx: tx.Hex}
	// send raw tx directly to arc
	res, rawTxSubmitErr := submit(ctx, arc, body, opts)

	if rawTxSubmitErr != nil {
		if arcErr, ok := rawTxSubmitErr.(broadcast.ArcError); ok && arcErr.Status == 460 { // no extended format error
			return s.submitAsEf(ctx, arc, tx, opts)
		}
		return nil, rawTxSubmitErr
	}

	return res, nil
}

func (s *rawtxSubmitStrategy) batchSubmit(ctx context.Context, arc *ArcClient, txs []*broadcast.Transaction, opts *broadcast.TransactionOpts) ([]*broadcast.SubmittedTx, error) {
	body := convertToSubmittedTxs(txs)
	// send raw tx directly to arc
	res, rawTxSubmitErr := batchSubmit(ctx, arc, body, opts)

	if rawTxSubmitErr != nil {
		if arcErr, ok := rawTxSubmitErr.(broadcast.ArcError); ok {
			if arcErr.Status == 460 { // no extended format error
				return s.batchSubmitAsEf(ctx, arc, txs, opts)
			}
		}
		return nil, rawTxSubmitErr
	}

	return res, nil
}

func (s *rawtxSubmitStrategy) submitAsEf(ctx context.Context, arc *ArcClient, tx *broadcast.Transaction, opts *broadcast.TransactionOpts) (*broadcast.SubmittedTx, error) {
	efBody, convertErr := s.convertToEf(ctx, arc, tx)
	if convertErr != nil {
		return nil, convertErr
	}

	return submit(ctx, arc, efBody, opts)
}

func (s *rawtxSubmitStrategy) batchSubmitAsEf(ctx context.Context, arc *ArcClient, txs []*broadcast.Transaction, opts *broadcast.TransactionOpts) ([]*broadcast.SubmittedTx, error) {
	efBody := make([]*SubmitTxRequest, 0, len(txs))
	for _, tx := range txs {
		efTx, convertErr := s.convertToEf(ctx, arc, tx)
		if convertErr != nil {
			return nil, convertErr
		}

		efBody = append(efBody, efTx)
	}

	return batchSubmit(ctx, arc, efBody, opts)
}

func (s *rawtxSubmitStrategy) convertToEf(ctx context.Context, arc *ArcClient, tx *broadcast.Transaction) (*SubmitTxRequest, error) {
	transaction, err := bt.NewTxFromString(tx.Hex)
	if err != nil {
		return nil, arc_utils.WithCause(errors.New("rawTxRequest: bt.NewTxFromString failed"), err)
	}

	for _, input := range transaction.Inputs {
		if err = updateUtxoWithMissingData(arc, input); err != nil {
			return nil, arc_utils.WithCause(errors.New("rawTxRequest: updateUtxoWithMissingData() failed"), err)
		}
	}

	request := &SubmitTxRequest{
		RawTx: hex.EncodeToString(transaction.ExtendedBytes()),
	}
	return request, nil
}

func convertToSubmittedTxs(txs []*broadcast.Transaction) []*SubmitTxRequest {
	body := make([]*SubmitTxRequest, 0, len(txs))
	for _, tx := range txs {
		requestTx := &SubmitTxRequest{RawTx: tx.Hex}
		body = append(body, requestTx)
	}

	return body
}

func submit(ctx context.Context, arc *ArcClient, body *SubmitTxRequest, opts *broadcast.TransactionOpts) (*broadcast.SubmittedTx, error) {
	url := arc.apiURL + arcSubmitTxRoute

	data, err := json.Marshal(body)
	if err != nil {
		return nil, ErrSubmitTxMarshal
	}

	pld := httpclient.NewPayload(
		httpclient.POST,
		url,
		arc.token,
		data,
	)
	appendSubmitTxHeaders(&pld, opts, arc.deploymentID)

	return httpclient.RequestModel(
		ctx,
		arc.HTTPClient.DoRequest,
		pld,
		decodeSubmitResponseBody,
		parseArcError,
	)
}

func batchSubmit(ctx context.Context, arc *ArcClient, body []*SubmitTxRequest, opts *broadcast.TransactionOpts) ([]*broadcast.SubmittedTx, error) {
	url := arc.apiURL + arcSubmitBatchTxsRoute

	data, err := json.Marshal(body)
	if err != nil {
		return nil, ErrSubmitTxMarshal
	}

	pld := httpclient.NewPayload(
		httpclient.POST,
		url,
		arc.token,
		data,
	)
	appendSubmitTxHeaders(&pld, opts, arc.deploymentID)

	return httpclient.RequestModel(
		ctx,
		arc.HTTPClient.DoRequest,
		pld,
		decodeSubmitBatchResponseBody,
		parseArcError,
	)
}

func appendSubmitTxHeaders(pld *httpclient.HTTPRequest, opts *broadcast.TransactionOpts, deploymentID string) {
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

	if deploymentID != "" {
		pld.AddHeader(XDeploymentIDHeader, deploymentID)
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

func updateUtxoWithMissingData(arc *ArcClient, input *bt.Input) error {
	txid := input.PreviousTxIDStr()
	pld := httpclient.NewPayload(
		httpclient.GET,
		fmt.Sprintf("https://junglebus.gorillapool.io/v1/transaction/get/%s", txid),
		"",
		nil,
	)

	tx, err := httpclient.RequestModel(
		context.Background(),
		arc.HTTPClient.DoRequest,
		pld,
		decodeJunblebusTransaction,
		nil,
	)

	if err != nil {
		return fmt.Errorf("junglebus request failed: %w", err)
	}

	if len(tx.Transaction) == 0 {
		return errors.New("junglebus responded with empty tx.Transaction[]")
	}

	actualTx, err := bt.NewTxFromBytes(tx.Transaction)
	if err != nil {
		return arc_utils.WithCause(errors.New("converting junglebusTransaction.Transaction to bt.Tx failed"), err)
	}

	o := actualTx.Outputs[input.PreviousTxOutIndex]
	input.PreviousTxScript = o.LockingScript
	input.PreviousTxSatoshis = o.Satoshis
	return nil
}

func decodeJunblebusTransaction(resp *http.Response) (*junglebusTransaction, error) {
	tx := &junglebusTransaction{}
	err := arc_utils.DecodeResponseBody(resp.Body, &tx)
	if err != nil {
		return nil, fmt.Errorf("decodeJunblebusTransaction: %w", err)
	}

	fmt.Println(tx)
	return tx, nil
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
