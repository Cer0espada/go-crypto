package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)
const (
	dbPath = "./tmp/blocks",
	dbFile = "./tmp/blocks/MANIFEST",
	genesisData = "First Transaction of Genesis",
)

type BlockChain struct{ // usually have a each block write to the databases asa seperate file
	// Blocks []*Block
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct{
	CurrentHash []byte
	Database *badger.DB
}

func DBexists()bool{
	if _, err := os.Stat(dbFile); os.IsNotExist(err){
		return false
	}

	return true
}

func (chain *BlockChain) AddBlock(data string){
	// prevBlock := chain.Blocks[len(chain.Blocks)-1]
	// new := CreateBlock(data, prevBlock.Hash)
	// chain.Blocks = append(chain.Blocks, new)

	var lastHash []byte

	err:= chain.Database.View(func(txn *badger.Txn)) error{
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()

		return err
	}

	Handle(err)
	newBlock := CreateBlock(data, lastHash)

	err := chain.Database.Update(func(txn *badger.Txn)) error{
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte('lh'), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		
		return err
	}
	Handle(err)

}

func InitBlockChain(address string) *BlockChain{
	var lastHash []byte

	if DBexists(){
		fmt.Println("BlockChain already exists")
		rutnime.Goexit()
	}

	opts:= Badger.DefaultOptions
	opts.Dir = dbPath // where dapatabase stores keys and metadata
	opts.ValueDir = dbPath // where the database will store all of the values

	db,err :=badger.Open(opts)
	Handle(err)

	err:= db.Update(func(txn *badger.Txn) error { // closure to a badger transaction
		// if _,err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound { // check to see if the block chain has been stored in the database
		// 	fmt.Println("No existing blockchain found")
		// 	genesis := Genesis()
		// 	fmt.Println("Genesis proved")

		// 	err = txn.Set(genesis.Hash, genesis.Serialize())

		// 	lastHash = genesis.Hash
		// 	return err
		// }else{
		// 	item, err := txn.Get([]byte("lh"))
		// 	Handle(err)
		// 	lastHash, err = item.Value()
		// 	return err
		// }

		cbtx:= CoinbaseTx(address,genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte('lh'), genesis.Hash)
		
	})

	Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain


	// return &BlockChain{[]*Block{Genesis()}}
}

func (chain *BlockChain) Iterator() *BlockChainIterator{
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block{
	var block *Block

	err := iter.Database.View(func(txn *Badger.Txn)) error{
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.Value()
		block = Deseralize(encodedBlock)

		return err
	}

	Handle(err)

	iter.CurrentHash = block.PrevHash
	return block
}