package fixtures

import "github.com/bitcoin-sv/go-broadcast-client/broadcast"

const (
	ProviderMain                                 = "https://mocked_arc_api_url.com/arc"
	ProviderSecondary                            = "https://secondary_mocked_arc_api_url.com/arc"
	MaxScriptSizePolicy       int64              = 100000000
	MaxTxSigOpsCountPolicy    int64              = 4294967295
	MaxTxSizePolicy           int64              = 100000000
	MaxTxSizePolicySecondary  int64              = 220000000
	MiningFeeBytes            int64              = 1000
	SatoshisPerBytes          int64              = 1
	SatoshisPerBytesSecondary int64              = 2
	Timestamp                                    = "2023-09-05T17:03:49.537230128Z"
	TimestampSecondary                           = "2023-09-05T17:05:29.736256927Z"
	TxResponseStatus                             = 200
	TxResponseTitle                              = "OK"
	TxStatus                  broadcast.TxStatus = "SEEN_ON_NETWORK"
	TxBlockHash                                  = "0000000000000000019a575e0ea4d9bbe251dd24c473a0d8407935973151f282"
	TxBlockHashSecondary                         = "0000000000000000045c969f3acd5db37896aba95f91389f2d191496bf15584b"
	TxBlockHeight             int64              = 800182
	TxBlockHeightSecondary    int64              = 799439
	TxExtraInfo                                  = ""
)
