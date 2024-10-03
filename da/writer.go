package da

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/zeromq/goczmq.v4"

	"github.com/Layer-Edge/bitcoin-da/config"
	"github.com/Layer-Edge/bitcoin-da/db"
)

var (
	BtcCliPath     = ""
	BashScriptPath = ""
)

type HashBlockProcessor struct {
	layerEdgeClient *ethclient.Client
}

func (proc *HashBlockProcessor) Process(msg [][]byte, protocolId string) ([]byte, error) {
	topic := string(msg[0])
	fmt.Printf("Topic: %s\n", topic)

	layerEdgeHeader, err := proc.layerEdgeClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println("Error getting layerEdgeHeader: ", err)
		return nil, err
	}
	dhash := layerEdgeHeader.Hash()
	log.Println("Latest LayerEdge Block Hash:", dhash.Hex())

	data := append([]byte(protocolId), dhash.Bytes()...)
	hash, err := CallScriptWithData(hex.EncodeToString(data))
	return hash, err
}

func CallScriptWithData(data string) ([]byte, error) {
	cmd := exec.Command(BashScriptPath+"/op_return_transaction.sh", data)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "BTC_CLI_PATH="+BtcCliPath)
	return cmd.Output()
}

// HashBlockSubscriber subscribes to hash blocks and processes them.
func HashBlockSubscriber(cfg *config.Config) {
	if err := db.InitDB("da.db"); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	channeler := goczmq.NewSubChanneler(cfg.ZmqEndpointHashBlock, "hashblock")
	BashScriptPath = cfg.BashScriptPath
	BtcCliPath = cfg.BtcCliPath

	if channeler == nil {
		log.Fatal("Error creating channeler", channeler)
	}
	defer channeler.Destroy()

	layerEdgeClient, err := ethclient.Dial(cfg.LayerEdgeRPC.HTTP)
	if err != nil {
		log.Fatal("Error creating layerEdgeClient: ", err)
	}

	hashProcessor := &HashBlockProcessor{layerEdgeClient: layerEdgeClient}
	channelReader := &ZmqChannelReader{}

	if !channelReader.Subscribe(cfg.ZmqEndpointHashBlock, "hashblock") {
		log.Fatal("Failed to subscribe to hashblock")
	}

	counter := 0

	for {
		select {
		case msg, ok := <-channelReader.channeler.RecvChan:
			if !ok || (counter%cfg.WriteIntervalBlock) != 0 || len(msg) != 3 {
				continue // Skip processing based on conditions.
			}

			hash, err := hashProcessor.Process(msg, cfg.ProtocolId)
			if err != nil {
				log.Println("Error writing -> ", err)
				continue
			}
			counter++
			log.Println("Relayer Write Done -> ", strings.ReplaceAll(string(hash[:]), "\n", ""))

			if err := db.InsertTxnHash(string(hash)); err != nil {
				log.Println("Error inserting transaction hash into DB:", err)
			}
		}
	}
}
