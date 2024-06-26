// Package broadcast provides a set of functions to broadcast or query transactions to the Bitcoin SV network using the client provided.
package broadcast

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/utils"
)

// ErrClientUndefined is returned when the client is undefined.
// Example:
//
//	func (a *ArcClient) QueryTransaction(ctx context.Context, txID string) (*broadcast.QueryTxResponse, error) {
//		if a == nil {
//			return nil, broadcast.ErrClientUndefined
//		}
//
// It should be returned for all defined clients in the future.
var ErrClientUndefined = errors.New("client is undefined")

// ErrAllBroadcastersFailed is returned when all configured broadcasters failed to query or broadcast the transaction.
var ErrAllBroadcastersFailed = errors.New("all broadcasters failed")

// ErrUnableToDecodeResponse is returned when the http response cannot be decoded.
var ErrUnableToDecodeResponse = errors.New("unable to decode response")

// ErrMissingStatus is returned when the tx status is missing.
var ErrMissingStatus = errors.New("missing tx status")

// ErrStrategyUnknown is returned when the strategy provided is unknown.
// Example:
//
// func NewBroadcaster(strategy Strategy, factories ...BroadcastFactory) broadcast.Client
// Calling NewBroadcaster we need to provide a strategy, if the strategy is unknown (we don't have an implementation for that) we return ErrStrategyUnknown.
var ErrStrategyUnknown = errors.New("unknown strategy")

// ErrNoMinerResponse is returned when no response is received from any miner.
var ErrNoMinerResponse = errors.New("failed to get reponse from any miner")

// ArcError is general type for the error returned by the ArcClient.
type ArcError struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Detail    string `json:"detail"`
	Instance  string `json:"instance,omitempty"`
	Txid      string `json:"txid,omitempty"`
	ExtraInfo string `json:"extraInfo,omitempty"`
}

// IsRejectedTransaction returns true if the transaction is in rejected status.
func (err *ArcError) IsRejectedTransaction() bool {
	const RejectedStatus = 109
	return err.Status == RejectedStatus
}

// Error returns the error string it's the implementation of the error interface.
func (err *ArcError) Error() string {
	sb := strings.Builder{}

	sb.WriteString("arc error: {")
	sb.WriteString(fmt.Sprintf("type: %s, title: %s, status: %d, detail: %s",
		err.Type, err.Title, err.Status, err.Detail))

	if err.Instance != "" {
		sb.Write([]byte(fmt.Sprintf(", instance: %s", err.Instance)))
	}

	if err.Txid != "" {
		sb.Write([]byte(fmt.Sprintf(", txid: %s", err.Txid)))
	}

	if err.ExtraInfo != "" {
		sb.Write([]byte(fmt.Sprintf(", extraInfo: %s", err.ExtraInfo)))
	}

	sb.WriteString("}")
	return sb.String()
}

func (err *ArcError) Is(target error) bool {
	var arcError *ArcError
	return errors.As(target, &arcError)
}

// Failure returns a new FailureResponse with the description and the error.
func Failure(description string, err error) error {
	return utils.WithCause(errors.New(description), err)
}
