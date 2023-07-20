package broadcast

import (
	"context"
)

type BestQuoter interface {
	// BestQuote(ctx context.Context, feeCategory, feeType string) (*FeeQuoteResponse, error)
	BestQuote(ctx context.Context, feeCategory, feeType string) error
}

// type FastestQuoter interface {
// 	// FastestQuote(ctx context.Context, timeout time.Duration) (*FeeQuoteResponse, error)
// 	FastestQuote(ctx context.Context, timeout time.Duration) error
// }

// type FeeQuoter interface {
// 	// FeeQuote(ctx context.Context, miner *Miner) (*FeeQuoteResponse, error)
// 	FeeQuote(ctx context.Context) error
// }

// type PolicyQuoter interface {
// 	// PolicyQuote(ctx context.Context, miner *Miner) (*PolicyQuoteResponse, error)
// 	PolicyQuote(ctx context.Context) error
// }

// type TransactionQuerier interface {
// 	// // QueryTransaction(ctx context.Context, miner *Miner, txID string, opts ...QueryTransactionOptFunc) (*QueryTransactionResponse, error)
// 	QueryTransaction(ctx context.Context, txID string) error
// }

// type TransactionSubmitter interface {
// 	// SubmitTransaction(ctx context.Context, miner *Miner, tx *Transaction) (*SubmitTransactionResponse, error)
// 	SubmitTransaction(ctx context.Context) error
// }

// type TransactionsSubmitter interface {
// 	// SubmitTransactions(ctx context.Context, miner *Miner, txs []Transaction) (*SubmitTransactionsResponse, error)
// 	SubmitTransactions(ctx context.Context) error
// }
