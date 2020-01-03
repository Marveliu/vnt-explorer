package common

const (
	DefaultPageSize        = 100
	DefaultOffset          = 0
	DefaultOrder           = "desc"
	DefaultHydrantCount    = 100
	DefaultHydrantInterval = 3600
	DefaultNodeInterval    = 300
	DefaultHydrantChainId  = 1333
	DefaultGasLimit        = 90000
	DefaultGasPrice        = 500000000000
	DefaultNodeStatus      = -1
	// VntTotal               = "10000000000000000000000000000"
	VntTotal   = "1000000000000000000000000000000000000000000000000000000000"
	VntDecimal = 18
	ImagePath  = "static/image/"
)

const (
	RpcBlockNumber        = "core_blockNumber"
	RpcGetBlockByNumber   = "core_getBlockByNumber"
	RpcGetTxByHash        = "core_getTransactionByHash"
	RpcGetTxReceipt       = "core_getTransactionReceipt"
	RpcGetBalance         = "core_getBalance"
	RpcCall               = "core_call"
	RpcGetAllCandidates   = "core_getAllCandidates"
	RpcSendRawTransaction = "core_sendRawTransaction"
	GetTransactionCount   = "core_getTransactionCount"
)

const (
	HContentType = "application/json; charset=utf-8"
)
