package arc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bitcoin-sv/go-broadcast-client/common"
	"github.com/bitcoin-sv/go-broadcast-client/config"
	"github.com/bitcoin-sv/go-broadcast-client/internal/decoder"
	"github.com/bitcoin-sv/go-broadcast-client/internal/httpclient"
	"github.com/bitcoin-sv/go-broadcast-client/models"
	"github.com/bitcoin-sv/go-broadcast-client/shared"
)

var ErrSubmitTxMarshal = fmt.Errorf("error while marshalling submit tx body")

func (a *ArcClient) SubmitTransaction(ctx context.Context, tx *common.Transaction) (*models.SubmitTxResponse, error) {
	if a == nil {
		return nil, shared.ErrClientUndefined
	}

	result, err := submitTransaction(ctx, a, tx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func submitTransaction(ctx context.Context, arc *ArcClient, tx *common.Transaction) (*models.SubmitTxResponse, error) {
	sb := strings.Builder{}
	sb.WriteString(arc.apiURL + config.ArcSubmitTxRoute)

	data, err := createSubmitTxBody(tx)
	if err != nil {
		return nil, err
	}

	pld := httpclient.NewPayload(
		httpclient.POST,
		sb.String(),
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

	if err := validateSubmitTxResponse(model); err != nil {
		return nil, err
	}

	return model, nil
}

func createSubmitTxBody(tx *common.Transaction) ([]byte, error) {
	body := map[string]string{
		"rawtx": tx.RawTx,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, ErrSubmitTxMarshal
	}

	return data, nil
}

func appendSubmitTxHeaders(pld *httpclient.HTTPPayload, tx *common.Transaction) {
	if tx.MerkleProof {
		pld.AddHeader("X-MerkleProof", "true")
	}

	if tx.CallBackURL != "" {
		pld.AddHeader("X-CallbackUrl", tx.CallBackURL)
	}

	if tx.CallBackToken != "" {
		pld.AddHeader("X-CallbackToken", tx.CallBackToken)
	}

	if statusCode, ok := common.MapTxStatusToInt(tx.WaitForStatus); ok {
		pld.AddHeader("X-WaitForStatus", strconv.Itoa(statusCode))
	}
}

func decodeSubmitTxBody(body io.ReadCloser) (*models.SubmitTxResponse, error) {
	result, err := decoder.NewDecoder[models.SubmitTxResponse](body).Result()

	if err != nil || &result == nil {
		return nil, shared.ErrUnableToDecodeResponse
	}

	return &result, nil
}

func validateSubmitTxResponse(model *models.SubmitTxResponse) error {
	if model.TxStatus == "" {
		return shared.ErrMissingStatus
	}

	return nil
}
