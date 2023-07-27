package arc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	errors "github.com/bitcoin-sv/go-broadcast-client"
	"github.com/bitcoin-sv/go-broadcast-client/config"
	"github.com/bitcoin-sv/go-broadcast-client/internal/httpclient"
	"github.com/bitcoin-sv/go-broadcast-client/models"
)

var ErrMissingHash = fmt.Errorf("missing tx hash")
var ErrMissingStatus = fmt.Errorf("missing tx status")

func (a *ArcClient) QueryTransaction(ctx context.Context, txID string) (*models.QueryTxResponse, error) {

	if a == nil {
		return nil, errors.ErrClientUndefined
	}

	result, err := queryTransaction(ctx, a, txID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// queryTransaction will fire the HTTP request to retrieve the tx status and details
func queryTransaction(ctx context.Context, arc *ArcClient, txHash string) (*models.QueryTxResponse, error) {
	sb := strings.Builder{}
	sb.WriteString(arc.apiURL + config.ArcQueryTxRoute + txHash)

	payload := httpclient.NewPayload(
		httpclient.GET,
		sb.String(),
		arc.token,
		nil,
	)

	resp, err := arc.httpClient.DoRequest(
		ctx,
		payload,
	)
	if err != nil {
		return nil, err
	}

	model, err := decodeQueryTxBody(resp.Body)
	if err != nil {
		return nil, err
	}

	err = validateQueryTxResponse(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func decodeQueryTxBody(body io.ReadCloser) (*models.QueryTxResponse, error) {
	model := &models.QueryTxResponse{}
	err := json.NewDecoder(body).Decode(model)

	if err != nil || model == nil {
		return nil, err
	}

	return model, nil
}

func validateQueryTxResponse(model *models.QueryTxResponse) error {

	if model.BlockHash == "" {
		return ErrMissingHash
	}

	if model.TxStatus == "" {
		return ErrMissingStatus
	}

	return nil
}
