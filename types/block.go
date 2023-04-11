package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/cbergoon/merkletree"
	"github.com/miguelpgnferreira/glockchain/crypto"

	pb "github.com/golang/protobuf/proto"
	"github.com/miguelpgnferreira/glockchain/proto"
)

type TxHash struct {
	hash []byte
}

func NewTxhash(hash []byte) TxHash {
	return TxHash{
		hash: hash,
	}
}

func (h TxHash) CalculateHash() ([]byte, error) {
	return h.hash, nil
}

func (h TxHash) Equals(other merkletree.Content) (bool, error) {
	equals := bytes.Equal(h.hash, other.(TxHash).hash)
	return equals, nil
}

func VerifyBlock(b *proto.Block) bool {
	if len(b.Transactions) > 0 {
		if !VerifyRootHash(b) {
			fmt.Println("INVALID root hash")
			return false
		}
	}

	if len(b.PublicKey) != crypto.PubKeyLen {
		fmt.Println("INVALID public key length")
		return false
	}
	if len(b.Signature) != crypto.SignatureLen {
		fmt.Println("INVALID signature length")
		return false
	}
	sig := crypto.SignatureFromBytes(b.Signature)
	pubKey := crypto.PublicKeyFromBytes(b.PublicKey)
	hash := HashBlock(b)
	if !sig.Verify(pubKey, hash) {
		fmt.Println("INVALID signature")
		return false
	}
	return true
}

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	if len(b.Transactions) > 0 {
		tree, err := GetMerkleTree(b)
		if err != nil {
			panic(err)
		}
		b.Header.RootHash = tree.MerkleRoot()
	}
	hash := HashBlock(b)
	sig := pk.Sign(hash)
	b.PublicKey = pk.Public().Bytes()
	b.Signature = sig.Bytes()

	return sig
}

func VerifyRootHash(b *proto.Block) bool {
	tree, err := GetMerkleTree(b)
	if err != nil {
		return false
	}
	valid, err := tree.VerifyTree()
	if err != nil {
		return false
	}
	if len(b.Header.RootHash) == 0 {
		b.Header.RootHash = tree.MerkleRoot()
	}
	if !valid {
		return false
	}

	return bytes.Equal(b.Header.RootHash, tree.MerkleRoot())
}

func GetMerkleTree(b *proto.Block) (*merkletree.MerkleTree, error) {
	list := make([]merkletree.Content, len(b.Transactions))
	for i := 0; i < len(b.Transactions); i++ {
		list[i] = NewTxhash(HashTransaction(b.Transactions[i]))
	}

	//Create a new Merkle Tree from the list of Content
	t, err := merkletree.NewTree(list)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)
}

func HashHeader(header *proto.Header) []byte {
	h, err := pb.Marshal(header)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(h)

	return hash[:]
}
