package arc

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/bitcoin-sv/go-broadcast-client/config"
	"github.com/bitcoin-sv/go-broadcast-client/internal/decoder"
	"github.com/bitcoin-sv/go-broadcast-client/internal/httpclient"
	"github.com/bitcoin-sv/go-broadcast-client/models"
	"github.com/bitcoin-sv/go-broadcast-client/shared"
)

var ErrMissingHash = fmt.Errorf("missing tx hash")

func (a *ArcClient) QueryTransaction(ctx context.Context, txID string) (*models.QueryTxResponse, error) {

	if a == nil {
		return nil, shared.ErrClientUndefined
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

	pld := httpclient.NewPayload(
		httpclient.GET,
		sb.String(),
		arc.token,
		nil,
	)

	resp, err := arc.httpClient.DoRequest(
		ctx,
		pld,
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
	result, err := decoder.NewDecoder[models.QueryTxResponse](body).Result()

	if err != nil || &result == nil {
		return nil, shared.ErrUnableToDecodeResponse
	}

	return &result, nil
}

func validateQueryTxResponse(model *models.QueryTxResponse) error {

	if model.BlockHash == "" {
		return ErrMissingHash
	}

	if model.TxStatus == "" {
		return shared.ErrMissingStatus
	}

	return nil
}
