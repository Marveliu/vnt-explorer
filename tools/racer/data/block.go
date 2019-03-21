package data

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/astaxie/beego"
	"github.com/bluele/gcache"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/common/utils"
	"github.com/vntchain/vnt-explorer/models"
	"github.com/vntchain/vnt-explorer/tools/racer/token"
)

var acctCache = gcache.New(10000).LRU().Build()

const (
	ACC_TYPE_NULL     = 0
	ACC_TYPE_NORMAL   = 1
	ACC_TYPE_CONTRACT = 2
	ACC_TYPE_TOKEN    = 3
)

func GetLocalHeight() (int64, *models.Block) {
	b := &models.Block{}
	count, err := b.Count()
	if err != nil {
		msg := fmt.Sprintf("Failed to get block count: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	block, err := b.Last()
	if err != nil {
		msg := fmt.Sprintf("Failed to get last block: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	if block == nil && count > 0 {
		msg := fmt.Sprintf("Block data in db not matched! count %d not equal to lastest block number %d, please check you local database.", count, 0)
		beego.Error(msg)
		panic(msg)
	}

	var bNumber uint64

	if block == nil {
		bNumber = 0
	} else {
		bNumber = block.Number
	}

	if bNumber != uint64(count) {
		msg := fmt.Sprintf("Block data in db not matched! count %d not equal to lastest block number %d, please check you local database.", count, bNumber)
		beego.Error(msg)
		panic(msg)
	}

	return count, block
}

func GetRemoteHeight() int64 {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_BlockNumber

	err, resp := utils.CallRpc(rpc)
	if err != nil {
		panic(err.Error())
	}

	beego.Info("Response body", resp)

	blockNumber := utils.Hex(resp.Result.(string)).ToInt64()

	return blockNumber
}

func GetBlock(number int64) (*models.Block, []interface{}, []interface{}) {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetBlockByNumber

	hex := utils.Encode(big.NewInt(number).Bytes())
	if strings.HasPrefix(hex, "0x0") {
		hex = "0x" + hex[3:]
	}

	rpc.Params = append(rpc.Params, hex, false)

	err, resp := utils.CallRpc(rpc)
	if err != nil {
		panic(err.Error())
	}

	blockMap := resp.Result.(map[string]interface{})

	beego.Info("BlockMap: ", blockMap)

	bNumber := utils.Hex(blockMap["number"].(string)).ToUint64()

	timestamp := utils.Hex(blockMap["timestamp"].(string)).ToUint64()

	size := utils.Hex(blockMap["size"].(string)).ToUint64()

	gasUsed := utils.Hex(blockMap["gasUsed"].(string)).ToUint64()

	gasLimit := utils.Hex(blockMap["gasLimit"].(string)).ToUint64()

	b := &models.Block{
		Number:     bNumber,
		TimeStamp:  timestamp,
		Hash:       blockMap["hash"].(string),
		ParentHash: blockMap["parentHash"].(string),
		Producer:   blockMap["producer"].(string),
		Size:       fmt.Sprintf("%d", size),
		GasUsed:    gasUsed,
		GasLimit:   gasLimit,
		ExtraData:  blockMap["extraData"].(string),
	}

	var txs, witnesses []interface{}
	var ok bool
	txIs := blockMap["transactions"].([]interface{})
	beego.Info("txs: ", txIs)
	if txs, ok = blockMap["transactions"].([]interface{}); !ok {
		txs = make([]interface{}, 0)
	}

	if witnesses, ok = blockMap["witnesses"].([]interface{}); !ok {
		witnesses = make([]interface{}, 0)
	}

	return b, txs, witnesses
}

func GetTx(txHash string) *models.Transaction {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetTxByHash

	rpc.Params = append(rpc.Params, txHash)

	err, resp := utils.CallRpc(rpc)
	if err != nil {
		panic(err.Error())
	}

	txMap := resp.Result.(map[string]interface{})
	beego.Info("Transaction: ", txMap)

	rpc.Method = common.Rpc_GetTxReceipt

	err, resp = utils.CallRpc(rpc)
	if err != nil {
		panic(err.Error())
	}

	receiptMap := resp.Result.(map[string]interface{})
	beego.Info("Transaction: ", receiptMap)

	tx := &models.Transaction{
		Hash:        txMap["hash"].(string),
		From:        txMap["from"].(string),
		Value:       utils.Hex(txMap["value"].(string)).ToString(),
		GasLimit:    utils.Hex(txMap["gas"].(string)).ToUint64(),
		GasPrice:    utils.Hex(txMap["gasPrice"].(string)).ToString(),
		GasUsed:     utils.Hex(receiptMap["gasUsed"].(string)).ToUint64(),
		Nonce:       utils.Hex(txMap["nonce"].(string)).ToUint64(),
		Index:       utils.Hex(txMap["transactionIndex"].(string)).ToInt(),
		Input:       txMap["input"].(string),
		Status:      utils.Hex(receiptMap["status"].(string)).ToInt(),
		BlockNumber: utils.Hex(txMap["blockNumber"].(string)).ToUint64(),
	}

	var to string
	var ok bool
	if to, ok = txMap["to"].(string); !ok {
		to = ""

		beego.Info("This is a transaction of contract creation.")
		if contractAddr, ok := receiptMap["contractAddress"].(string); ok {
			tx.ContractAddr = contractAddr
		}
		tx.To = nil
	} else {
		tx.To = &models.Account{Address: to}
	}

	return tx
}

// Extract Account from a transaction
func ExtractAcct(tx *models.Transaction) {

	if tx.Status == 0 {
		return
	}
	from := tx.From
	to := tx.To
	contractAddr := tx.ContractAddr

	if a := GetAccount(from); a == nil {
		beego.Info("Block:", tx.BlockNumber, ", will insert normal account:", from)
		NewAccount(from, tx, ACC_TYPE_NORMAL, 1)
	} else {
		beego.Info("Block:", tx.BlockNumber, ", will update normal account:", from)
		UpdateAccount(a, tx, ACC_TYPE_NORMAL)
	}

	if to != nil && to.Address != "" {
		if a := GetAccount(to.Address); a == nil {
			beego.Info("Block:", tx.BlockNumber, ", will insert normal account:", to)
			NewAccount(to.Address, tx, ACC_TYPE_NORMAL, 1)
		} else {
			if a.IsToken {
				beego.Info("Block:", tx.BlockNumber, ", will update token account:", to)
				UpdateAccount(a, tx, ACC_TYPE_TOKEN)

				// Update the tx
				tx.IsToken = true
				err := tx.Update()
				if err != nil {
					msg := fmt.Sprintf("Failed to update transaction: %s, error: %s", tx.Hash, err.Error())
					beego.Error(msg)
					panic(msg)
				}
			} else if a.IsContract {
				beego.Info("Block:", tx.BlockNumber, ", will update contract account:", to)
				UpdateAccount(a, tx, ACC_TYPE_CONTRACT)
			} else {
				beego.Info("Block:", tx.BlockNumber, ", will update normal account:", from)
				UpdateAccount(a, tx, ACC_TYPE_NORMAL)
			}
		}
	} else if contractAddr != "" { // this case is for contract creation
		if a := GetAccount(contractAddr); a == nil {
			// new contract account
			beego.Info("Block:", tx.BlockNumber, ", will insert contract account:", contractAddr)
			NewAccount(contractAddr, tx, ACC_TYPE_CONTRACT, 1)
		} else if !a.IsContract {
			// this account already exists as a normal account,
			// will change it to a contract account
			//a.IsContract = true
			beego.Info("Block:", tx.BlockNumber, ", will update contract account:", contractAddr)
			UpdateAccount(a, tx, ACC_TYPE_CONTRACT)
		}
	}
	return
}

func GetBalance(addr string, blockNumber uint64) string {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetBlance

	rpc.Params = append(rpc.Params, addr)
	rpc.Params = append(rpc.Params, utils.EncodeUint64(blockNumber))

	err, resp := utils.CallRpc(rpc)
	if err != nil {
		panic(err.Error())
	}

	balance := utils.Hex(resp.Result.(string)).ToString()
	return balance
}

func IsToken(addr string, tx *models.Transaction) (bool, *token.Erc20) {
	totalSupply := token.GetTotalSupply(addr, tx.BlockNumber)
	tokenName := token.GetTokenName(addr, tx.BlockNumber)
	decimals := token.GetDecimals(addr, tx.BlockNumber)
	symbol := token.GetSymbol(addr, tx.BlockNumber)

	if totalSupply != nil && decimals != nil && symbol != "" && tokenName != "" {
		erc20 := &token.Erc20{
			Address:     addr,
			TokenName:   tokenName,
			TotalSupply: totalSupply,
			Symbol:      symbol,
			Decimals:    decimals,
		}

		// Update the tx
		tx.IsToken = true
		err := tx.Update()
		if err != nil {
			msg := fmt.Sprintf("Failed to update transaction: %s, error: %s", tx.Hash, err.Error())
			beego.Error(msg)
			panic(msg)
		}

		return true, erc20
	}

	return false, nil
}

// Insert a new Account, in this case, tye _type only could be "normal" or "contract"
func NewAccount(addr string, tx *models.Transaction, _type int, txCount uint64) {
	a := &models.Account{
		Address:        addr,
		Vname:          addr, //todo: get vname
		Balance:        "0",
		TxCount:        txCount,
		FirstBlock:     tx.BlockNumber,
		LastBlock:      tx.BlockNumber,
		TokenAmount:    "0",
		TokenAcctCount: "0",
		InitTx:         tx.Hash,
		LastTx:         tx.Hash,
	}

	if _type == ACC_TYPE_CONTRACT {
		a.IsContract = true
		a.ContractName = "" // TODO: extract contract name from contract code
		a.ContractOwner = tx.From

		beego.Info("######### a.ContractOwner:", a.ContractOwner)

		if ok, erc20 := IsToken(addr, tx); ok {
			a.IsToken = true

			// TODO: get token detail by calling the contract
			a.TokenType = token.TOKEN_ERC20
			a.ContractName = erc20.TokenName
			a.TokenAmount = erc20.TotalSupply.String()
			a.TokenSymbol = erc20.Symbol
			a.TokenDecimals = erc20.Decimals.Uint64()
			a.TokenAcctCount = "0"
			a.TokenLogo = ""
		}
	}

	a.Balance = GetBalance(addr, tx.BlockNumber)
	InsertAcc(addr, a)
}

func UpdateAccount(account *models.Account, tx *models.Transaction, _type int) {

	account.Balance = GetBalance(account.Address, tx.BlockNumber)
	account.LastBlock = tx.BlockNumber

	if account.LastTx != tx.Hash {
		account.LastTx = tx.Hash
		account.TxCount += 1
	}

	retAddrs := make([]string, 0)

	if _type == ACC_TYPE_CONTRACT {
		// if already exists as a normal account,
		// then now it turns out a new contract account
		if !account.IsContract {
			if ok, erc20 := IsToken(account.Address, tx); ok {
				account.IsToken = true
				account.TokenType = token.TOKEN_ERC20
				account.ContractName = erc20.TokenName
				account.TokenAmount = erc20.TotalSupply.String()
				account.TokenSymbol = erc20.Symbol
				account.TokenDecimals = erc20.Decimals.Uint64()
				account.TokenAcctCount = "0"
				account.TokenLogo = ""
			}
		}
		account.IsContract = true
		account.ContractOwner = tx.From
	} else if _type == ACC_TYPE_TOKEN {
		tx.IsToken = true
		retAddrs = token.UpdateTokenBalance(account, tx)
	}

	UpdateAcc(account.Address, account)
	// Save the accounts in token transfer
	for _, a := range retAddrs {
		if acct := GetAccount(a); acct != nil {
			acct.Balance = GetBalance(a, tx.BlockNumber)
			acct.LastBlock = tx.BlockNumber
			if acct.LastTx != tx.Hash {
				acct.LastTx = tx.Hash
				acct.TxCount += 1
			}
			UpdateAcc(a, acct)
		} else {
			NewAccount(a, tx, ACC_TYPE_NORMAL, 1)
			beego.Info("Inserted accounts: ", a)
		}
	}
}

func PersistWitnesses(accts []string, blockNumber uint64) {
	beego.Info("Will persist witnesses accounts: ", accts)
	for _, a := range accts {
		if acct := GetAccount(a); acct != nil {
			acct.Balance = GetBalance(a, blockNumber)
			acct.LastBlock = blockNumber
			UpdateAcc(a, acct)
		} else {
			NewAccount(a, &models.Transaction{BlockNumber: blockNumber}, ACC_TYPE_NORMAL, 0)
			beego.Info("Inserted witness account: ", a)
		}
	}
}

func GetAccount(addr string) *models.Account {
	addr = strings.ToLower(addr)
	if _type, err := acctCache.Get(addr); err == nil && _type != nil {
		beego.Info("Address hit in cache:", addr)
		return _type.(*models.Account)
	} else {
		beego.Info("Address not hit in cache:", addr)
		a := &models.Account{}
		a, err := a.Get(addr)
		if err != nil {
			beego.Info("Address not hit in db:", addr)
			return nil
		}
		beego.Info("Address hit in db:", addr)
		acctCache.Set(addr, a)
		return a
	}
}

// insert into db and cache
func InsertAcc(addr string, acct *models.Account) {
	addr = strings.ToLower(addr)
	if err := acct.Insert(); err != nil {
		msg := fmt.Sprintf("Failed to insert account: %v, error: %s", acct, err.Error())
		beego.Error(msg)
		panic(msg)
	}
	acctCache.Set(addr, acct)
}

// update db and cache
func UpdateAcc(addr string, acct *models.Account) {
	addr = strings.ToLower(addr)
	if err := acct.Update(); err != nil {
		msg := fmt.Sprintf("Failed to update account: %s, error: %s", addr, err.Error())
		beego.Error(msg)
		panic(err)
	}
	acctCache.Set(addr, acct)
}
