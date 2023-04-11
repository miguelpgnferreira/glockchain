package node

import (
	"fmt"
	"testing"

	"github.com/miguelpgnferreira/glockchain/crypto"
	"github.com/miguelpgnferreira/glockchain/proto"
	"github.com/miguelpgnferreira/glockchain/types"
	"github.com/miguelpgnferreira/glockchain/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.GeneratePrivateKey()
	b := util.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	require.Nil(t, err)
	b.Header.PrevHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, b)
	return b
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	assert.Equal(t, 0, chain.Height())

	_, err := chain.GetBlockByHeight(0)
	assert.Nil(t, err)
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	for i := 0; i < 100; i++ {
		b := randomBlock(t, chain)
		require.Nil(t, chain.AddBlock(b))
		require.Equal(t, chain.Height(), i+1)
	}
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())

	for i := 0; i < 100; i++ {
		block := randomBlock(t, chain)

		blockHash := types.HashBlock(block)

		require.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlock)

		fetchedBlockByheight, err := chain.GetBlockByHeight(i + 1)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlockByheight)
	}
}

func TestAddBlockWithTxInsufficientFunds(t *testing.T) {
	var (
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
		block     = randomBlock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedString(godSeed)
		recipient = crypto.GeneratePrivateKey().Public().Address().Bytes()
	)

	prevTx, err := chain.txStore.Get("29215d9313649e394524b53fe0cd85e3cafb678655339bb1ea214b15c1dff902")
	assert.Nil(t, err)
	fmt.Println(prevTx)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  1001,
			Address: recipient,
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	require.NotNil(t, chain.AddBlock(block))
}

func TestAddBlockWithTx(t *testing.T) {
	var (
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
		block     = randomBlock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedString(godSeed)
		recipient = crypto.GeneratePrivateKey().Public().Address().Bytes()
	)

	prevTx, err := chain.txStore.Get("29215d9313649e394524b53fe0cd85e3cafb678655339bb1ea214b15c1dff902")
	assert.Nil(t, err)
	fmt.Println(prevTx)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  100,
			Address: recipient,
		},
		{
			Amount:  900,
			Address: privKey.Public().Address().Bytes(),
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privKey, block)
	require.Nil(t, chain.AddBlock(block))
}
