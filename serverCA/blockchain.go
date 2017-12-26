package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fatih/color"
)

type Block struct {
	Hash         string    `json:"hash"`
	Height       int64     `json:"height"`
	Date         time.Time `json:"date"`
	PreviousHash string    `json:"previoushash"`
	NextHash     string    `json:"nexthash"`
	Data         []Address `json:"data"`
	Emitter      string    `json:"emitter"` //the ID of the peer that has emmited the block
}

type Blockchain struct {
	GenesisBlock string    `json:"genesisblock"`
	LastUpdate   time.Time `json:"lastupdate"`
	Blocks       []Block   `json:"blocks"`
}

var blockchain Blockchain

func (bc *Blockchain) getBlockByHash(hash string) (Block, error) {
	for _, block := range bc.Blocks {
		if block.Hash == hash {
			return block, nil
		}
	}
	var b Block
	return b, errors.New("Block Hash not found")
}

func (bc *Blockchain) blockExists(block Block) bool {
	for _, b := range bc.Blocks {
		if b.Hash == block.Hash {
			return true
		}
	}
	return false
}

func (bc *Blockchain) addBlock(block Block) error {
	if blockchain.blockExists(block) {
		return errors.New("[Error adding Block]: Block already exists in the Blockchain")
	}
	if len(bc.Blocks) > 0 {
		bc.Blocks[len(bc.Blocks)-1].NextHash = block.Hash
	} else {
		bc.GenesisBlock = block.Hash
	}
	bc.Blocks = append(bc.Blocks, block)

	return nil
}

func reconstructBlockchainFromBlock(urlAPI string, h string) {
	color.Yellow("reconstructing the blockchain from last block in memory")
	var block Block
	var err error

	block, err = blockchain.getBlockByHash(h)
	if err != nil {
		fmt.Println("reconstructBlockFromBlock: block with " + h + " not found. Getting genesis block")
	}

	if h == "" {
		//no genesis block yet
		color.Green(urlAPI + "/blocks/genesis")
		res, err := http.Get(urlAPI + "/blocks/genesis")
		check(err)
		body, err := ioutil.ReadAll(res.Body)
		check(err)
		err = json.Unmarshal(body, &block)
		check(err)
		color.Yellow("[New Block]: " + block.Hash)
		err = blockchain.addBlock(block)
		check(err)
	} else {
		block.NextHash = h
	}

	for block.NextHash != "" && block.Hash != "" {
		res, err := http.Get(urlAPI + "/blocks/next/" + block.Hash)
		check(err)
		body, err := ioutil.ReadAll(res.Body)
		check(err)
		err = json.Unmarshal(body, &block)
		check(err)
		if block.Hash != "" {
			color.Yellow("[New Block]: " + block.Hash)
			err = blockchain.addBlock(block)
			check(err)
		}
	}
	blockchain.print()
}

func (bc *Blockchain) print() {
	color.Green("Printing Blockchain stored in memory")
	color.Green("Genesis Block: " + bc.GenesisBlock)
	for _, b := range bc.Blocks {
		color.Green("Block height:")
		fmt.Println(b.Height)
		color.Green("Hash: " + b.Hash)
		color.Green("Date: " + b.Date.String())
		color.Green("---")
	}
}
