package node

import (
	"context"
	"fmt"
	"net"

	"github.com/miguelpgnferreira/glockchain/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	peers map[net.Addr]*grpc.ClientConn
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.None, error) {
	peer, _ := peer.FromContext(ctx)
	fmt.Println("Received tx from: ", peer)
	return nil, nil
}
