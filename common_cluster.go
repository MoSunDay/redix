package main

import (
	"fmt"
	"strconv"
	"strings"
)

func clusterSlotsInit(c Context) {
	if SlotCache.CM.Get("0") == "" {
		for i := 0; i <= 16383; i++ {
			addr := SlotCache.Opts.RaftTCPAddress
			addrSlice := strings.Split(addr, ":")
			portStr := addrSlice[1]
			port, err := strconv.Atoi(portStr)
			if err != nil {
				SlotCache.Log.Fatal("cluster slots verification failed")
				c.WriteError("init slots failed err: " + err.Error())
				return
			}
			portStr = strconv.Itoa(port + 200)
			addr = addrSlice[0] + ":" + portStr
			SlotCache.CM.Set(strconv.Itoa(i), addr)
		}
		SlotCache.Log.Println("slots init done")
	} else {
		c.WriteError("slots have been allocated")
		return
	}
	c.WriteString("OK")
}

func clusterInfo(c Context) {
	raftStats := SlotCache.Raft.Raft.Stats()

	future := SlotCache.Raft.Raft.GetConfiguration()
	if err := future.Error(); err != nil {
		c.WriteError("could not get configuration for stats")
		return
	}
	configuration := future.Configuration()

	epoch := raftStats["term"]
	size := len(configuration.Servers)
	clusterStatus := "ok"
	if SlotCache.CM.Get("0") == "" {
		clusterStatus = "down"
	}

	c.WriteBulkString(fmt.Sprintf(""+
		"cluster_state:%s\n"+
		"cluster_slots_assigned:16384\n"+
		"cluster_slots_ok:16384\n"+
		"cluster_slots_pfail:0\n"+
		"cluster_slots_fail:0\n"+
		"cluster_known_nodes:%d\n"+
		"cluster_size:%d\n"+
		"cluster_current_epoch:%s\n"+
		"cluster_my_epoch:%s\n"+
		"cluster_stats_messages_sent:0\n"+
		"cluster_stats_messages_received:0\n",
		clusterStatus, size, size, epoch, epoch,
	))
}

func clusterHelp(c Context) {
	c.WriteString("CLUSTER [ help | NODES | SLOTS | INIT ]")
}

func clusterNodes(c Context) {
	future := SlotCache.Raft.Raft.GetConfiguration()
	if err := future.Error(); err != nil {
		c.WriteError("could not get configuration for stats")
		return
	}
	configuration := future.Configuration()

	response := make([]string, len(configuration.Servers))
	nodeSlots := getNodeSlots()
	for _, server := range configuration.Servers {

		addr := string(server.Address)
		addrSlice := strings.Split(addr, ":")
		portStr := addrSlice[1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			c.WriteError("cluster slots verification failed")
			return
		}
		portStr = strconv.Itoa(port + 200)
		addr = addrSlice[0] + ":" + portStr
		uuid := make([]byte, 40)

		for i := 0; i < 40; i++ {
			if i < len(addr) {
				if addr[i] != '.' && addr[i] != ':' {
					uuid[i] = addr[i]
				} else {
					uuid[i] = 'a'
				}
			} else {
				uuid[i] = 'b'
			}
		}

		nodeInfo := fmt.Sprintf("%s %s@%s myself,master - 0 0 connected %s\r\n", string(uuid), addr, portStr, nodeSlots[addr])
		response = append(response, nodeInfo)
	}
	c.WriteBulkString(strings.Join(response, ""))
}

func getNodeSlots() map[string]string {
	info := make(map[string][]string)
	nodeSlots := make(map[string]string)

	headIndex := 0
	head := SlotCache.CM.Get(strconv.Itoa(headIndex))
	cur := head
	for i := 1; i <= 16383; i++ {
		curIndex := strconv.Itoa(i)
		cur = SlotCache.CM.Get(curIndex)
		if head == cur {
			continue
		} else {
			headInexToString := strconv.Itoa(headIndex)
			indexToString := strconv.Itoa(i - 1)
			if headInexToString == indexToString {
				info[head] = append(info[head], indexToString)
			} else {
				info[head] = append(info[head], headInexToString+"-"+indexToString)
			}
			head = cur
			headIndex = i
		}
	}

	if head == cur {
		info[head] = append(info[head], strconv.Itoa(headIndex)+"-16383")
	}

	for k, v := range info {
		nodeSlots[k] = strings.Join(v, " ")
	}
	return nodeSlots
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
		"init":  clusterSlotsInit,
		"test":  clusterTest,
	}
	if fn, ok := subCommand[c.args[0]]; ok {
		fn(c)
	} else {
		clusterHelp(c)
	}
}
