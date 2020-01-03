package data

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"runtime"
	"strings"

	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/models"
	"github.com/vntchain/vnt-explorer/tools/racer/pool"
	"path"
)

const (
	ActionInsert = 1
	ActionUpdate = 2

	NodeTask = "nodes"
)

var (
	BlockPool          = pool.New(runtime.NumCPU()*3, 50)
	BlockInsertPool    = pool.New(runtime.NumCPU()*3, 50)
	TxPool             = pool.New(runtime.NumCPU()*3, 6000)
	ReportPool         = pool.New(runtime.NumCPU()*3, 6000)
	AccountExtractPool = pool.New(runtime.NumCPU()*3, 6000)
	AccountPool        = pool.New(runtime.NumCPU()*3, 10000)
	WitnessesPool      = pool.New(runtime.NumCPU()*3, 100)
	NodePool           = pool.New(runtime.NumCPU()*3, 100)
	NodeInfoPool       = pool.New(runtime.NumCPU()*3, 100)
	LogoPool           = pool.New(runtime.NumCPU()*3, 100)
)

type BlockTask struct {
	pool.BasicTask
	BlockNumber int64
}

func (t *BlockTask) DoWork(workRoutine int) {
	t.PreDoWork(workRoutine)
	PersistBlock(t.BlockNumber)
}

func NewBlockTask(BlockNumber int64) *BlockTask {
	return &BlockTask{
		BasicTask: pool.BasicTask{
			Name: fmt.Sprintf("Block-%d", BlockNumber),
			Pool: BlockPool,
		},
		BlockNumber: BlockNumber,
	}
}

type BlockInsertTask struct {
	pool.BasicTask
	Block *models.Block
}

func (t *BlockInsertTask) DoWork(workRoutine int) {
	t.PreDoWork(workRoutine)
	beego.Debug("Will insert block:", t.Block.Number)
	err := t.Block.Insert()
	if err != nil {
		msg := fmt.Sprintf("Failed to insert or update block: %v, error: %s,", t.Block, err.Error())
		panic(msg)
	}
}

func NewBlockInsertTask(Block *models.Block) *BlockInsertTask {
	return &BlockInsertTask{
		BasicTask: pool.BasicTask{
			Name: fmt.Sprintf("Block-Insert-%d", Block.Number),
			Pool: BlockPool,
		},
		Block: Block,
	}
}

type TxTask struct {
	pool.BasicTask
	Tx *models.Transaction
}

func (t *TxTask) DoWork(workRoutine int) {
	t.PreDoWork(workRoutine)
	err := t.Tx.Insert()
	if err != nil {
		msg := fmt.Sprintf("Failed to insert transaction: %s", err.Error())
		panic(msg)
	}
}

func NewTxTask(Tx *models.Transaction) *TxTask {
	return &TxTask{
		BasicTask: pool.BasicTask{
			Name: fmt.Sprintf("Tx-%s", Tx.Hash),
			Pool: TxPool,
		},
		Tx: Tx,
	}
}

type ReportTask struct {
	pool.BasicTask
	Reports []*models.Report
}

func (t *ReportTask) DoWork(workRoutine int) {
	t.PreDoWork(workRoutine)
	o := orm.NewOrm()
	if _, err := o.InsertMulti(100, t.Reports); err != nil {
		msg := fmt.Sprintf("Failed to insert reports: %s", err.Error())
		panic(msg)
	}

}

func NewReportTask(txHash string, reports []*models.Report) *ReportTask {
	return &ReportTask{
		BasicTask: pool.BasicTask{
			Name: fmt.Sprintf("Report-%s", txHash),
			Pool: ReportPool,
		},
		Reports: reports,
	}
}

type ExtractAccountTask struct {
	pool.BasicTask
	Tx *models.Transaction
}

func (t *ExtractAccountTask) DoWork(workRoutine int) {
	t.PreDoWork(workRoutine)
	ExtractAcct(t.Tx)
}

func NewExtractAccountTask(Tx *models.Transaction) *ExtractAccountTask {
	return &ExtractAccountTask{
		pool.BasicTask{
			fmt.Sprintf("ext-account-%s", Tx.Hash),
			AccountExtractPool,
		},
		Tx,
	}
}

type AccountTask struct {
	pool.BasicTask
	Account *models.Account
	Action  int
}

func (t *AccountTask) DoWork(workRoutine int) {
	t.PreDoWork(workRoutine)
	switch t.Action {
	case ActionInsert:
		if err := t.Account.Insert(); err != nil {
			msg := fmt.Sprintf("Failed to insert account: %v, error: %s", t.Account, err.Error())
			beego.Error(msg)
			panic(msg)
		}
		acctCache.Set(t.Account.Address, t.Account)
		break
	case ActionUpdate:
		if err := t.Account.Update(); err != nil {
			msg := fmt.Sprintf("Failed to update account: %s, error: %s", t.Account.Address, err.Error())
			beego.Error(msg)
			panic(err)
		}
		acctCache.Set(t.Account.Address, t.Account)
		break
	default:

	}
}

func NewAccountTask(Account *models.Account, Action int) *AccountTask {
	return &AccountTask{
		pool.BasicTask{
			fmt.Sprintf("account-%s", Account.Address),
			AccountPool,
		},
		Account,
		Action,
	}
}

type WitnessesTask struct {
	pool.BasicTask
	Witnesses   []string
	BlockNumber uint64
}

func (t *WitnessesTask) DoWork(workRoutine int) {
	t.PreDoWork(workRoutine)
	PersistWitnesses(t.Witnesses, t.BlockNumber)
}

func NewWitnessesTask(Witnesses []string, BlockNumber uint64) *WitnessesTask {
	return &WitnessesTask{
		pool.BasicTask{
			"witnesses",
			AccountPool,
		},
		Witnesses,
		BlockNumber,
	}
}

type NodesTask struct {
	pool.BasicTask
}

func (t *NodesTask) DoWork(workRoutine int) {

	t.PreDoWork(workRoutine)
	witnesses := GetWitnesses(-1)
	witMap := make(map[string]int)
	for _, w := range witnesses {
		addr := fmt.Sprintf("%v", w)
		witMap[strings.ToLower(addr)] = 1
	}

	nodes := GetNodes()
	for _, node := range nodes {
		// fmt.Println("node address: %s", node.Address)
		if witMap[node.Address] == 1 {
			node.IsSuper = 1
		} else {
			node.IsSuper = 0
		}
		dbNode := &models.Node{}
		dbNode.Get(node.Address)

		// register account's Vname
		account := GetAccount(node.Address)
		if account != nil && account.Vname != node.Vname {
			account.Vname = node.Vname
			updateAcc(account)
		}

		// new node or node's home update, or node's location is unknown
		// try to get nodeInfo otherwise copy the old data
		if dbNode == nil {
			PostNodeInfoTask(NewNodeInfoTask(node))
		} else if dbNode.Home != node.Home ||
			(dbNode.Latitude == 360 && dbNode.Longitude == 360) ||
			dbNode.Logo == "" ||
			(dbNode.Status == 0 && node.Status == 1) {
			node.IsAlive = dbNode.IsAlive
			PostNodeInfoTask(NewNodeInfoTask(node))
		} else {
			node.Longitude = dbNode.Longitude
			node.Latitude = dbNode.Latitude
			node.City = dbNode.City
			node.Logo = dbNode.Logo
			node.IsAlive = dbNode.IsAlive
		}

		// if logo file doesn't exist, try to download it
		if len(node.Logo) > 0 {
			logoUrlList := strings.Split(node.Logo, ";")
			for _, logoUrl := range logoUrlList {
				imgName := path.Base(logoUrl)
				imgPath := path.Join(common.ImagePath, node.Address, imgName)
				if exists, _, _ := FileExists(imgPath); !exists {
					if logoUrl != "" {
						PostLogoTask(NewLogoTask(logoUrl, node.Address))
					}
				}
			}
		}

		if err := node.Insert(); err != nil {
			msg := fmt.Sprintf("Failed to insert node: %s", err.Error())
			panic(msg)
		}
	}
}

func NewNodesTask() *NodesTask {
	return &NodesTask{
		pool.BasicTask{
			Name: NodeTask,
			Pool: AccountPool,
		},
	}
}

type NodeInfoTask struct {
	pool.BasicTask
	Node *models.Node
}

func (t *NodeInfoTask) DoWork(workRoutine int) {
	t.PreDoWork(workRoutine)

	if len(t.Node.Home) == 0 {
		return
	}

	// nodeInfo := GetBpInfo(t.Node.Home + "/bp.json")
	nodeInfo := GetBpInfo("https://" + t.Node.Home + "/bp.json")
	beego.Info("Get nodeInfo ", t.Node.Home, nodeInfo)
	if nodeInfo != nil {
		t.Node.Latitude = nodeInfo.Location.Latitude
		t.Node.Longitude = nodeInfo.Location.Longitude
		t.Node.City = nodeInfo.Location.Name
		logoUrlList := []string{
			nodeInfo.Branding.Logo_256,
			nodeInfo.Branding.Logo_1024,
			nodeInfo.Branding.Logo_Svg,
		}
		nodeLogoList := []string{"", "", ""}
		for i, url := range logoUrlList {
			if url != "" {
				nodeLogoList[i] = url
				PostLogoTask(NewLogoTask(url, t.Node.Address))
			}
		}

		t.Node.Logo = strings.Join(nodeLogoList, ";")
		if err := t.Node.Insert(); err != nil {
			msg := fmt.Sprintf("Failed to insert node: %s", err.Error())
			beego.Error(msg)
		}
	}
}

func NewNodeInfoTask(Node *models.Node) *NodeInfoTask {
	return &NodeInfoTask{
		pool.BasicTask{
			"nodeInfo",
			NodeInfoPool,
		},
		Node,
	}
}

type LogoTask struct {
	pool.BasicTask
	imgUrl  string
	address string
}

func (t *LogoTask) DoWork(workRoutine int) {
	t.PreDoWork(workRoutine)
	GetLogo(t.imgUrl, t.address)
}

func NewLogoTask(imgUrl, address string) *LogoTask {
	return &LogoTask{
		pool.BasicTask{
			"logo",
			LogoPool,
		},
		imgUrl,
		address,
	}
}

func PostBlockTask(task *BlockTask) {
	err := BlockPool.PostWork("block", task)
	if err != nil {
		beego.Error("区块线程池满载！")
		panic("")
	}
}

func PostInsertBlockTask(task *BlockInsertTask) {
	err := BlockInsertPool.PostWork("block", task)
	if err != nil {
		beego.Error("区块插入线程池满载！")
		panic("")
	}
}

func PostTxTask(task *TxTask) {
	err := TxPool.PostWork("tx", task)
	if err != nil {
		beego.Error("交易线程池满载！")
		panic("")
	}
}

func PostReportTasks(task *ReportTask) {
	if err := ReportPool.PostWork("report", task); err != nil {
		beego.Error("交易线程池满载！")
		panic("")
	}
}

func PostExtractAccountTask(task *ExtractAccountTask) {
	err := AccountExtractPool.PostWork("ext-account", task)
	if err != nil {
		beego.Error("帐户线程池满载！")
		panic("")
	}
}

func PostAccountTask(task *AccountTask) {
	err := AccountPool.PostWork("account", task)
	if err != nil {
		beego.Error("帐户线程更新插入池满载！")
		panic("")
	}
}

func PostWitnessesTask(task *WitnessesTask) {
	err := AccountPool.PostWork("witnesses", task)
	if err != nil {
		beego.Error("Witnesses池满载！")
		panic("")
	}
}

func PostNodesTask(task *NodesTask) {
	err := NodePool.PostWork(NodeTask, task)
	if err != nil {
		beego.Error("Nodes池满载！")
		panic("")
	}
}

func PostNodeInfoTask(task *NodeInfoTask) {
	err := NodeInfoPool.PostWork("nodeInfo", task)
	if err != nil {
		beego.Error("NodeInfo池满载！")
	}
}

func PostLogoTask(task *LogoTask) {
	err := NodeInfoPool.PostWork("logo", task)
	if err != nil {
		beego.Error("Logo池满载！")
	}
}
