package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index     int
	Timestamp string
	Data      string
	PrevHash  string
	Hash      string
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Data + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateBlock(prevBlock Block, data string) Block {
	newBlock := Block{
		Index:     prevBlock.Index + 1,
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
	}
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

func GenesisBlock() Block {
	genesis := Block{
		Index:     0,
		Timestamp: time.Now().String(),
		Data:      "Genesis Block",
		PrevHash:  "",
	}
	genesis.Hash = calculateHash(genesis)
	return genesis
}
