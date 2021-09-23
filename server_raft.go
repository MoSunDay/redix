package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	rcache "github.com/MoSunDay/redix/rcache"
)

func runRatfServer() {
	opts := rcache.NewOptions()
	opts.DataDir = *flageRaftDataDir + "/" + *flagRaftNode
	opts.HttpAddress = *flagHttpAddress
	opts.Bootstrap = *flagRaftBootstrap
	opts.RaftTCPAddress = *flagRaftTCPAddress
	opts.JoinAddress = *flagRaftJoinAddress

	SlotCache = &rcache.Cached{
		Opts: opts,
		Log:  log.New(os.Stderr, "Cached: ", log.Ldate|log.Ltime),
		CM:   rcache.NewCacheManager(),
	}
	ctx := &rcache.CachedContext{SlotCache}

	raft, err := rcache.NewRaftNode(SlotCache.Opts, ctx)
	if err != nil {
		SlotCache.Log.Fatal(fmt.Sprintf("new raft node failed:%v", err))
	}
	SlotCache.Raft = raft

	if SlotCache.Opts.JoinAddress != "" {
		err = rcache.JoinRaftCluster(SlotCache.Opts)
		if err != nil {
			SlotCache.Log.Fatal(fmt.Sprintf("join raft cluster failed:%v", err))
		}
	}

	logger := log.New(os.Stderr, "httpserver: ", log.Ldate|log.Ltime)
	httpServer := rcache.NewHttpServer(ctx, logger)
	SlotCache.HttpServer = httpServer

	go func() {
		for {
			select {
			case leader := <-SlotCache.Raft.LeaderNotifyCh:
				if leader {
					SlotCache.Log.Println("become leader, enable write api")
					SlotCache.HttpServer.SetWriteFlag(true)
					if SlotCache.CM.Get("0") == "" {
						for i := 0; i <= 16383; i++ {
							addr := SlotCache.Opts.RaftTCPAddress
							addrSlice := strings.Split(addr, ":")
							portStr := addrSlice[1]
							port, err := strconv.Atoi(portStr)
							if err != nil {
								SlotCache.Log.Fatal("cluster slots verification failed")
							}
							portStr = strconv.Itoa(port + 200)
							addr = addrSlice[0] + ":" + portStr
							SlotCache.CM.Set(strconv.Itoa(i), addr)
						}
						SlotCache.Log.Println("slots init done")
					}
				} else {
					SlotCache.Log.Println("become follower, close write api")
					SlotCache.HttpServer.SetWriteFlag(false)
				}
			}
		}
	}()

	go func() {
		var l net.Listener
		var err error
		l, err = net.Listen("tcp", SlotCache.Opts.HttpAddress)
		if err != nil {
			logger.Fatal(fmt.Sprintf("listen %s failed: %s", SlotCache.Opts.HttpAddress, err))
		}
		logger.Printf("http server listen:%s", l.Addr())
		http.Serve(l, httpServer.Mux)
	}()
}
