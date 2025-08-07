package core

type Blockchain struct {
	Blocks []Block
}

func NewBlockchain() *Blockchain {
	genesis := GenesisBlock()
	return &Blockchain{
		Blocks: []Block{genesis},
	}
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := GenerateBlock(prevBlock, data)
	bc.Blocks = append(bc.Blocks, newBlock)
}
