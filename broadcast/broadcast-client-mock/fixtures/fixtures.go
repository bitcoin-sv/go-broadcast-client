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
	TxMerklePath                                 = "fd123511fc17127984cea2846b6a6b3965cbd523ea92736c11022372dc922d0013ee15f577bb7d832ce6102e22a15c6e4f4e4214e246d769f04baaf458e3407f5117de8eceb78782a1f90bfaa103915dcc3aa940c7e1476398db47dc62833f3149c51593cf9eeaec36f86c871d64e849a486b97eee1ee863c18aa28e599e3b30853a8563a87441202a97f707d0614cfc024f73ba4a2848ec63e19bc5fb168e994c2c2c105b5951dee4775f20879cb63b3b525370b35200e9f91e9e0fa234dfae5cc4792d17ab08f5cb29aa97e99742dc1e48a9d7e648be33db2d75b806519d88c62a57b39ef900cf68c9d765a1a3b17895d05077ab33280fdb4c09f818d31065b1325b6eb0366c4fd48c2840c1e3a5eb64a9596ffe53bbb315b44784882fb727a98e791d4c8a0bc1d5a88fd85365afdf83da86216cf81b587ee960b5c2a4174dfb26fa2606addf5e51edea58f76cc415f21876c9fa4551891aab537f4e463f1ded910d6a27ff27531c32b4dad15360b35d4d6da348c0c311aaf412bd19a92e2253e90c67104c8f0c0167cf43032f4ffa616cb8c5c28733a3faef5afb83d3c621b252112cbf35859c7b3b4caf820ae43f7054db87da2b7bc735f41b0969b4caadd1cc78be66fb367990103b5682f33605628dcc9af941bf555c64e1c48d3d0955f2c53f7ed7611d976d3effb053eb13b85ca2bff429516aa21227c39d6b6f4106f513069a559805977dee6462d354522d6e42f30d60744bd0f08471577ea5e1f63f319261"
	TxId                                         = "680d975a403fd9ec90f613e87d17802c029d2d930df1c8373cdcdda2f536a1c0"
	TxIdSecondary                                = "d4f4bcb8decdbc75a3dcc54580c76457ee663a28b843956733570b38f7c73b5a"
)
