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
	TxMerklePathSecondary                        = "fd13400f3e40f8eedd367dbe5dee632fb577e513fcc8d73cd1f6032ee6f710e95c5907d71f67fbc7027c4652a67752c45b90057f3758f6582aaad728c5ea7a29b29f53798f6ce2a931f86c84f665b31fa5247669dbf7501cba7c562c4f11a58654fc1a8a57f78cffbced9bcc18925ed46f1002f22a4f56edf29c403dc9ca5846994a89f54038e8c68b898d17e505dd1e861b1dd20ffd8736dfc3c503706ee27ec1b852cbcdb33a564c0e8f0959cf1b758bdf0164f318de57c1e674d11e5b6b0c5336a662edfeb161ced6b16916f0e4f709ba9498b3d0440dbff65ca82c39f64e3487d3d9e7e46c8251eb2e443202a19b060ddf3dcbd6414371ff0544de2370d1948ba8c5a9ef29350e6475061fa82428b1b199f7d5e7a2ad28ccba7c802a5ba38b038ec3e4710e1483440e507f7e188d35693082d441d48e4be16e5da0f6e3203b93f295334a08ead059024a0750cbc88ee77a6b3101203ff8eb17353b58f5d03568084fd22b22467e500f72fa368021c8c8f0fa559b41ec3c01c6c355a9186a7d8bda348bf5ee7035d66e8b32ed57ebfcff619e893ea9a9873bd63c8247f4467d91e16fa256fc586a56e652850ee65b69bad583c72860f30ac04f6f8d71acc4840853353b9ae3c1da33107bb631e55d0c6b474d42c95f42611a50e21ff4271350fc4868"
	TxId                                         = "680d975a403fd9ec90f613e87d17802c029d2d930df1c8373cdcdda2f536a1c0"
	TxIdSecondary                                = "469dd0f63fe4b3fe972dc72d28057e931abd69d8dfc85bf6e623efa41d50ef73"
)
