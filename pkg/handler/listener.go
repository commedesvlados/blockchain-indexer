package handler

import (
	"context"
	"fmt"
	token "github.com/commedesvlados/blockchain-indexer/contracts"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"strings"
)

func (h *Handler) BlockListener(ctx context.Context) {
	fmt.Println("block listener start")

	headers := make(chan *types.Header)
	sub, err := h.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:

			block, err := h.client.BlockByHash(ctx, header.Hash())
			if err != nil {
				log.Fatalln(err)
			}

			averageBlockGasPrice := uint64(0)

			for _, tx := range block.Transactions() {
				blockUsedGas := uint64(0)
				blockUsedGas += tx.GasPrice().Uint64()
				averageBlockGasPrice = blockUsedGas / uint64(block.Transactions().Len())
			}

			// listen block
			fmt.Printf("\nblock hash: %s\n", block.Hash().Hex())
			fmt.Printf("block number: %d\n", block.Number().Uint64())
			fmt.Printf("average used gas price in block: %d\n", averageBlockGasPrice)

			// add to db
			err = h.repository.AddAverageGas(block.Hash().Hex(), block.Number().Uint64(), averageBlockGasPrice)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

// Transfer & Approval log struct
type TokenLog struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

func (h *Handler) ERC20Listener(ctx context.Context) {
	fmt.Println("ERC20 listener start")

	// TODO: current block, current block + n
	// TODO: past address USDT
	query := createFilterQuery("", 100000000, 10000020)

	logs, err := h.client.FilterLogs(ctx, query)
	if err != nil {
		log.Fatalln(err)
	}

	contractABI, err := abi.JSON(strings.NewReader(string(token.TokenMetaData.ABI)))
	if err != nil {
		log.Fatalln(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(logApprovalSig)

	for _, vLog := range logs {
		fmt.Printf("\nLog Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			eventName := "Transfer"

			fmt.Printf("Log Name: %s\n")

			tokenLogs := unpackTokenLogs(contractABI, eventName, vLog)

			err := h.repository.AddERC20Logs(
				vLog.BlockNumber, uint64(vLog.Index), vLog.BlockHash.Hex(), eventName, tokenLogs.From.Hex(), tokenLogs.To.Hex(), tokenLogs.Value.String())
			if err != nil {
				log.Fatalln(err)
			}

		case logApprovalSigHash.Hex():
			eventName := "Approval"

			fmt.Printf("Log Name: %s\n", eventName)

			tokenLogs := unpackTokenLogs(contractABI, eventName, vLog)

			err := h.repository.AddERC20Logs(
				vLog.BlockNumber, uint64(vLog.Index), vLog.BlockHash.Hex(), eventName, tokenLogs.From.Hex(), tokenLogs.To.Hex(), tokenLogs.Value.String())
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func createFilterQuery(address string, blockFrom, blockTo int64) ethereum.FilterQuery {
	contractAddress := common.HexToAddress(address)

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(blockFrom),
		ToBlock:   big.NewInt(blockTo),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	return query
}

func unpackTokenLogs(contractAbi abi.ABI, eventName string, vLog types.Log) TokenLog {
	var tokenEvent TokenLog

	err := contractAbi.UnpackIntoInterface(&tokenEvent, eventName, vLog.Data)
	if err != nil {
		log.Fatalln(err)
	}

	tokenEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
	tokenEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

	fmt.Printf("From: %s\n", tokenEvent.From.Hex())
	fmt.Printf("To: %s\n", tokenEvent.To.Hex())
	fmt.Printf("Value: %s\n\n", tokenEvent.Value.String())

	return tokenEvent
}
