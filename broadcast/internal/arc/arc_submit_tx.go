package arc

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
)

var ErrSubmitTxMarshal = errors.New("error while marshalling submit tx body")

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
		return nil, err
	}

	model, err := decodeSubmitTxBody(resp.Body)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func createSubmitTxBody(tx *broadcast.Transaction) ([]byte, error) {
	body := map[string]string{
		"rawTx": tx.RawTx,
	}
	data, err := json.Marshal(body)
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

func decodeSubmitTxBody(body io.ReadCloser) (*broadcast.SubmitTxResponse, error) {
	model := broadcast.SubmitTxResponse{}
	err := json.NewDecoder(body).Decode(&model)

	if err != nil {
		return nil, broadcast.ErrUnableToDecodeResponse
	}

	return &model, nil
}

func validateSubmitTxResponse(model *broadcast.SubmitTxResponse) error {
	if model.TxStatus == "" {
		return broadcast.ErrMissingStatus
	}

	return nil
}
