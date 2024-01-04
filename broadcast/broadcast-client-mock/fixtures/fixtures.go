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
	Timestamp                                    = "2024-01-04T16:14:41.409761246Z"
	TimestampSecondary                           = "2023-09-05T17:05:29.736256927Z"
	TxResponseStatus                             = 200
	TxResponseTitle                              = "OK"
	TxStatus                  broadcast.TxStatus = "MINED"
	TxBlockHash                                  = "00000000000000000896d2b93efa4476c4bd47ed7a554aeac6b38044745a6257"
	TxBlockHashSecondary                         = "0000000000000000045c969f3acd5db37896aba95f91389f2d191496bf15584b"
	TxBlockHeight             int64              = 825599
	TxBlockHeightSecondary    int64              = 799439
	TxExtraInfo                                  = ""
	TxMerklePath                                 = "feff980c000a02fd180202b7fa685c09cefa4944c00a253fa2003cad0da4afe402dd9a4f268ee01e9e1ba1fd1902004249f4326179d24add2753da7094774fc32323456fe6d1c355f9d2c1d5cd984401fd0d0100d453d1a607bb11fc761b2bd71b6ff7f81275e5a33f1153ed518c180b128539c301870058bab5901cfcd0fc1ca8c6ea1cc70980f21210fdd1fdebe257d58706f6a9b73e0142006e8d8f5f8a4899f631e859eb1736dea7f2d468a54a16affd0347ce15583c75d4012000800ba451f3de17d05140a4bb6aa1da8e5fde8e18a4dd5e9a694fae710ba7fb66011100d28ac45b55d2034ea6a6d1f5a4b5fbacb4902e80606be52d9b78db92d9860173"
	TxMerklePathSecondary                        = "fd13400f3e40f8eedd367dbe5dee632fb577e513fcc8d73cd1f6032ee6f710e95c5907d71f67fbc7027c4652a67752c45b90057f3758f6582aaad728c5ea7a29b29f53798f6ce2a931f86c84f665b31fa5247669dbf7501cba7c562c4f11a58654fc1a8a57f78cffbced9bcc18925ed46f1002f22a4f56edf29c403dc9ca5846994a89f54038e8c68b898d17e505dd1e861b1dd20ffd8736dfc3c503706ee27ec1b852cbcdb33a564c0e8f0959cf1b758bdf0164f318de57c1e674d11e5b6b0c5336a662edfeb161ced6b16916f0e4f709ba9498b3d0440dbff65ca82c39f64e3487d3d9e7e46c8251eb2e443202a19b060ddf3dcbd6414371ff0544de2370d1948ba8c5a9ef29350e6475061fa82428b1b199f7d5e7a2ad28ccba7c802a5ba38b038ec3e4710e1483440e507f7e188d35693082d441d48e4be16e5da0f6e3203b93f295334a08ead059024a0750cbc88ee77a6b3101203ff8eb17353b58f5d03568084fd22b22467e500f72fa368021c8c8f0fa559b41ec3c01c6c355a9186a7d8bda348bf5ee7035d66e8b32ed57ebfcff619e893ea9a9873bd63c8247f4467d91e16fa256fc586a56e652850ee65b69bad583c72860f30ac04f6f8d71acc4840853353b9ae3c1da33107bb631e55d0c6b474d42c95f42611a50e21ff4271350fc4868"
	TxId                                         = "680d975a403fd9ec90f613e87d17802c029d2d930df1c8373cdcdda2f536a1c0"
	TxIdSecondary                                = "469dd0f63fe4b3fe972dc72d28057e931abd69d8dfc85bf6e623efa41d50ef73"
)
