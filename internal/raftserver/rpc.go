package raftserver

import (
	"context"
	"net"

	"github.com/btmorr/leifdb/internal/node"
	"github.com/btmorr/leifdb/internal/raft"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type server struct {
	raft.UnimplementedRaftServer
	Node *node.Node
}

// RequestVote handles RPC vote requests from other nodes
func (s *server) RequestVote(ctx context.Context, v *raft.VoteRequest) (*raft.VoteReply, error) {
	log.Debug().Msgf("Received vote request: %v", v)
	return s.Node.HandleVote(v), nil
}

// AppendLogs handles RPC log-append requests from other nodes
func (s *server) AppendLogs(ctx context.Context, a *raft.AppendRequest) (*raft.AppendReply, error) {
	log.Debug().Msgf("Received append request: %v", a)
	return s.Node.HandleAppend(a), nil
}

// StartRaftServer constructs and starts a gRPC server for Raft protocol routes
// Note: `port` must be in the form ":12345"
func StartRaftServer(lis net.Listener, n *node.Node) *grpc.Server {
	s := grpc.NewServer()
	raft.RegisterRaftServer(s, &server{Node: n})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC failed to serve")
		}
	}()
	return s
}
