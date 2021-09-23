package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

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
