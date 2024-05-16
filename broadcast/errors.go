// Package broadcast provides a set of functions to broadcast or query transactions to the Bitcoin SV network using the client provided.
package broadcast

import (
	"errors"
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

// ErrStrategyUnkown is returned when the strategy provided is unknown.
// Example:
//
// func NewBroadcaster(strategy Strategy, factories ...BroadcastFactory) broadcast.Client
// Calling NewBroadcaster we need to provide a strategy, if the strategy is unknown (we don't have an implementation for that) we return ErrStrategyUnkown.
var ErrStrategyUnkown = errors.New("unknown strategy")

// ErrNoMinerResponse is returned when no response is received from any miner.
var ErrNoMinerResponse = errors.New("failed to get reponse from any miner")

// ArcFailure is the interface for the error returned by the ArcClient.
type ArcFailure interface {
	error
	Details() *FailureResponse
}

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

// Details returns the details of the error it's the implementation of the ArcFailure interface.
func (failure *FailureResponse) Details() *FailureResponse {
	return failure
}

// Error returns the error string it's the implementation of the error interface.
func (err ArcError) Error() string {
	return "arc error"
}

// FailureResponse is the response returned by the ArcClient when the request fails.
type FailureResponse struct {
	Description      string
	ArcErrorResponse *ArcError
}

// Error returns the error string it's the implementation of the error interface.
func (failure *FailureResponse) Error() string {
	sb := strings.Builder{}
	sb.WriteString(failure.Description)

	if failure.ArcErrorResponse != nil {
		sb.WriteString(", ")
		sb.WriteString(failure.ArcErrorResponse.Error())
	}

	return sb.String()
}

// Failure returns a new FailureResponse with the description and the error.
func Failure(description string, err error) *FailureResponse {
	var arcErr ArcError
	if errors.As(err, &arcErr) {
		return &FailureResponse{
			Description:      description,
			ArcErrorResponse: &arcErr,
		}
	}

	return &FailureResponse{Description: utils.WithCause(errors.New(description), err).Error()}
}
