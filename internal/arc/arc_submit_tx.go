package arc

import (
	"context"
	"strings"

	"github.com/bitcoin-sv/go-broadcast-client/common"
	"github.com/bitcoin-sv/go-broadcast-client/config"
	"github.com/bitcoin-sv/go-broadcast-client/models"
)

func (a *ArcClient) SubmitTransaction(ctx context.Context, tx *common.Transaction) (*models.SubmitTxResponse, error) {
	return nil, nil
}

func submitTransaction(ctx context.Context, arc *ArcClient, tx *common.Transaction) (*models.SubmitTxResponse, error) {
	sb := strings.Builder{}
	sb.WriteString(arc.apiURL + config.ArcSubmitTxRoute)

	return nil, nil
}
