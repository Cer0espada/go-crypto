package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Take the data from the block

//create a counter (nonce) which starts at 0

//create a hash of the data plus the counter

//check the hash to see if it meets a set of requirments

// Requirements:
// The First few bytes must contain 0s in order to be valid -- hashcash

const Diffuculty = 12 // slowly increment difficulty over time

type ProofOfWork struct{
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork{ // allow us to take a pointer to a block then a pointer to a proof of work
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Diffuculty))

	pow := &ProofOfWork{b,target}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte{
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Diffuculty)),
		},
		[]byte{},
	)

	return data
}

// want to have the block rate -- the amount of blocks being created at a given time, the same
// keep the time to sign a block down 

func (pow *ProofOfWork) Run() (int, []byte){
	var intHash big.Int
	var hash [32]byte
	nonce := 0
	for nonce < math.MaxInt64{
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1{
			break
		}else{
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool{
	var intHash big.Int
	
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func ToHex (num int64) []byte{ // creates a buffer and decode it into bytes
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil{
		log.Panic(err)
	}

	return buff.Bytes()
}