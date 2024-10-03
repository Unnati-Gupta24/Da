package da

import (
    "bytes"
    "fmt"
    "log"

    "github.com/btcsuite/btcd/wire"
)

type RawBlockProcessor struct{}

func (proc *RawBlockProcessor) Process(msg [][]byte, protocolId string) ([]byte, error) {
    topic := string(msg[0])
    serializedBlock := msg[1]

    fmt.Printf("Topic: %s\n", topic)

    parsedBlock, err := parseBlock(serializedBlock)
    if err != nil {
        log.Printf("Failed to parse block: %v", err)
        return nil, err
    }

    printBlock(parsedBlock)
    
    // Additional processing logic can be added here.

    return nil, nil // Return appropriate value based on processing.
}

// parseBlock parses a serialized Bitcoin block.
func parseBlock(data []byte) (*wire.MsgBlock, error) {
    var block wire.MsgBlock
    err := block.Deserialize(bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    return &block, nil
}

// printBlock prints the details of a Bitcoin block.
func printBlock(block *wire.MsgBlock) {
    fmt.Println("Block Details:")
    fmt.Printf("  Block Header:\n")
    fmt.Printf("    Version: %d\n", block.Header.Version)
    fmt.Printf("    Previous Block: %s\n", block.Header.PrevBlock)
    fmt.Printf("    Merkle Root: %s\n", block.Header.MerkleRoot)
    fmt.Printf("    Timestamp: %s\n", block.Header.Timestamp)
}
