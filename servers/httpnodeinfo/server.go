package httpnodeinfo

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	chain "github.com/wuyazero/Elastos.ELA/blockchain"
	"github.com/wuyazero/Elastos.ELA/config"
	"github.com/wuyazero/Elastos.ELA/servers"
)

type Info struct {
	NodeVersion   string
	BlockHeight   uint32
	NeighborCnt   int
	Neighbors     []NgbNodeInfo
	HttpRestPort  int
	HttpWsPort    int
	HttpJsonPort  int
	HttpLocalPort int
	NodePort      uint16
	NodeId        string
}

type NgbNodeInfo struct {
	NgbId   string
	NbrAddr string
}

var templates = template.Must(template.New("info").Parse(page))

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var ngbrNodersInfo []NgbNodeInfo
	var node = servers.ServerNode

	neighbors := node.GetNeighborNodes()

	for i := 0; i < len(neighbors); i++ {
		ngbrNodersInfo = append(ngbrNodersInfo, NgbNodeInfo{
			NgbId:   fmt.Sprintf("0x%x", neighbors[i].ID()),
			NbrAddr: neighbors[i].Addr(),
		})
	}

	pageInfo := &Info{
		BlockHeight:  chain.DefaultLedger.Blockchain.BlockHeight,
		NeighborCnt:  len(neighbors),
		Neighbors:    ngbrNodersInfo,
		HttpRestPort: config.Parameters.HttpRestPort,
		HttpWsPort:   config.Parameters.HttpWsPort,
		HttpJsonPort: config.Parameters.HttpJsonPort,
		NodePort:     config.Parameters.NodePort,
		NodeId:       fmt.Sprintf("0x%x", node.ID()),
	}

	err := templates.ExecuteTemplate(w, "info", pageInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StartServer() {
	http.HandleFunc("/info", viewHandler)
	http.ListenAndServe(":"+strconv.Itoa(int(config.Parameters.HttpInfoPort)), nil)
}
