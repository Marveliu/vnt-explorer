package common

const (
	DefaultPageSize = 100
	DefaultOffset   = 0
	DefaultOrder = "desc"
)

const (
	Rpc_BlockNumber      = "core_blockNumber"
	Rpc_GetBlockByNumber = "core_getBlockByNumber"
	Rpc_GetTxByHash      = "core_getTransactionByHash"
	Rpc_GetTxReceipt     = "core_getTransactionReceipt"
	Rpc_GetBlance        = "core_getBalance"
	Rpc_Call             = "core_call"
	Rpc_GetAllCandidates = "core_getAllCandidates"
)

const (
	H_ContentType = "application/json; charset=utf-8"
)
