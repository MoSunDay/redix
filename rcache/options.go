package rcache

import (
	"flag"
)

type options struct {
	DataDir        string // data directory
	HttpAddress    string // http server address
	RaftTCPAddress string // construct Raft Address
	Bootstrap      bool   // start as master or not
	JoinAddress    string // peer address to join
}

func NewOptions() *options {
	opts := &options{}

	var HttpAddress = flag.String("http", "127.0.0.1:6000", "Http address")
	var RaftTCPAddress = flag.String("raft", "127.0.0.1:7000", "raft tcp address")
	var node = flag.String("node", "node1", "raft node name")
	var Bootstrap = flag.Bool("bootstrap", false, "start as raft cluster")
	var JoinAddress = flag.String("join", "", "join address for raft cluster")
	flag.Parse()

	opts.DataDir = "./" + *node
	opts.HttpAddress = *HttpAddress
	opts.Bootstrap = *Bootstrap
	opts.RaftTCPAddress = *RaftTCPAddress
	opts.JoinAddress = *JoinAddress
	return opts
}
