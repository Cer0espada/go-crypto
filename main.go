package main

import (
	"flag"
	"fmt"
	"go-crypto/blockchain"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct{
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage(){
	fmt.Println("Usage:")
	fmt.Println(" add -block Block_Data - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs(){
	if len(os.Args) < 2{
		cli.printUsage()
		runtime.Goexit() // exits the application by shutting down the goroutine
		//badger has to properly garbage collect before closing, if it doesnt it could potentially corrupt the data
	}
}

func (cli *CommandLine) addBlock(data string){
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block !")
}

func (cli *CommandLine) printChain(){
	iter := cli.blockchain.Iterator()

	for{
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0{
			break
		}
	}
}

func (cli *CommandLine) run(){
	cli.validateArgs()
	
	addBlockCmd := flag.NewFlagSet("Add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1]{
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	
	case "print":
		err:= printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}
	if addBlockCmd.Parsed(){
		if *addBlockData ==""{
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed(){
		cli.printChain()
	}
}
func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.run()

	// chain.AddBlock("First Block after Genesis")
	// chain.AddBlock("Second Block after Genesis")
	// chain.AddBlock("Third Block after Genesis")

	// for _, block := range chain.Blocks{
		
	// }
}

//concesus algorithms / proof of work algorithms -- forcing the network to do work to add work to the chain 
// bitcoin -- get fees by powering the network, makes the blocks more secure,proof of work validation, when a user does work to sign of block they show proof, the work is hard but the proof that the work was done is easy