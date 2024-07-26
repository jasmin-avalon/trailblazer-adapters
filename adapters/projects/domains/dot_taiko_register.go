package domains

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/taikoxyz/trailblazer-adapters/adapters"
)

var (
	logMintedDomainSigHash   = crypto.Keccak256Hash([]byte("MintedDomain(string,uint256,address,uint256)"))
	logNameRegisteredSigHash = crypto.Keccak256Hash([]byte("NameRegistered(uint256,string,bytes32,address,uint256,uint256)"))
	logProfileCreatedSigHash = crypto.Keccak256Hash([]byte("ProfileCreated(uint256,address,address,string,string,address,bytes,string,uint256)"))
)

type DotTaikoIndexer struct {
	TargetAddresses []common.Address
}

func NewDotTaikoIndexer() *DotTaikoIndexer {
	return &DotTaikoIndexer{TargetAddresses: []common.Address{
		common.HexToAddress("0xD7b837A0E388B4c25200983bdAa3EF3A83ca86b7"),
		common.HexToAddress("0xFb2Cd41a8aeC89EFBb19575C6c48d872cE97A0A5"),
		common.HexToAddress("0x01412AAba531Cc6ef630CE5059120999f824CDAF"),
	}}
}

func (indexer *DotTaikoIndexer) Addresses() []common.Address {
	return indexer.TargetAddresses
}

func (indexer *DotTaikoIndexer) IndexLogs(ctx context.Context, chainID *big.Int, client *ethclient.Client, logs []types.Log) ([]adapters.Whitelist, error) {
	var result []adapters.Whitelist
	for _, vLog := range logs {
		if !indexer.isRelevantLog(vLog.Topics[0]) {
			continue
		}
		transferData, err := indexer.ProcessLog(ctx, chainID, client, vLog)
		if err != nil {
			return nil, err
		}
		result = append(result, *transferData)
	}
	return result, nil
}

func (indexer *DotTaikoIndexer) isRelevantLog(topic common.Hash) bool {
	return topic.Hex() == logNameRegisteredSigHash.Hex() || topic.Hex() == logMintedDomainSigHash.Hex() || topic.Hex() == logProfileCreatedSigHash.Hex()
}

func (indexer *DotTaikoIndexer) ProcessLog(ctx context.Context, chainID *big.Int, client *ethclient.Client, vLog types.Log) (*adapters.Whitelist, error) {
	ownerHex := common.BytesToAddress(vLog.Topics[2].Bytes()[12:])

	block, err := client.BlockByNumber(ctx, big.NewInt(int64(vLog.BlockNumber)))
	if err != nil {
		return nil, err
	}

	return &adapters.Whitelist{
		User:        ownerHex,
		Time:        block.Time(),
		BlockNumber: block.Number().Uint64(),
	}, nil
}