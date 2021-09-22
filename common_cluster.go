package main

import (
	"fmt"
)

func clusterInfo(c Context) {
	size := 3
	epoch := 1

	c.WriteBulkString(fmt.Sprintf(""+
		"cluster_state:ok\n"+
		"cluster_slots_assigned:16384\n"+
		"cluster_slots_ok:16384\n"+
		"cluster_slots_pfail:0\n"+
		"cluster_slots_fail:0\n"+
		"cluster_known_nodes:%d\n"+
		"cluster_size:%d\n"+
		"cluster_current_epoch:%d\n"+
		"cluster_my_epoch:%d\n"+
		"cluster_stats_messages_sent:0\n"+
		"cluster_stats_messages_received:0\n",
		size, size, epoch, epoch,
	))
}

func clusterHelp(c Context) {
	c.WriteString("CLUSTER [ help | NODES | SLOTS ]")
}

func clusterNodes(c Context) {
	c.WriteBulkString("356a192b7913b04c54574d18c28d46e6395428ab 0.0.0.0:6666@6666 myself,master - 0 0 connected 0-5461\r\nda4b9237bacccdf19c0760cab7aec4a8359010b0 0.0.0.0:6667@6667 myself,master - 0 0 connected 5462-10922\r\n77de68daecd823babbb58edb1c8e14d7106e83bb 0.0.0.0:6668@6668 myself,master - 0 0 connected 10923-16383\r\n")
}

func clusterSlots(c Context) {
	c.WriteArray(1)
	c.WriteArray(3)
	c.WriteInt64(0)
	c.WriteInt64(16383)
	c.WriteArray(3)
	c.WriteBulkString("0.0.0.0")
	c.WriteInt64(6666)
	c.WriteBulkString("356a192b7913b04c54574d18c28d46e6395428ab")
}

func clusterTest(c Context) {
	c.WriteString("-MOVED 0 0.0.0.0:6666")
}

func clusterCommand(c Context) {
	subCommand := map[string]func(c Context){
		"help":  clusterHelp,
		"info":  clusterInfo,
		"nodes": clusterNodes,
		"slots": clusterSlots,
		"test":  clusterTest,
	}
	if fn, ok := subCommand[c.args[0]]; ok {
		fn(c)
	} else {
		clusterHelp(c)
	}
}
