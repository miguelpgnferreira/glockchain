package types

import (
	"crypto/sha256"

	"github.com/miguelpgnferreira/glockchain/crypto"

	pb "github.com/golang/protobuf/proto"
	"github.com/miguelpgnferreira/glockchain/proto"
)

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	return pk.Sign(HashBlock(b))
}

func HashBlock(block *proto.Block) []byte {
	b, err := pb.Marshal(block)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)

	return hash[:]
}
