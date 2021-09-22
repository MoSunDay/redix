package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	rcache "github.com/MoSunDay/redix/rcache"
)

func initRaftServer() {
	opts := rcache.NewOptions()
	opts.DataDir = *flageRaftDataDir + "/" + *flagRaftNode
	opts.HttpAddress = *flagHttpAddress
	opts.Bootstrap = *flagRaftBootstrap
	opts.RaftTCPAddress = *flagRaftTCPAddress
	opts.JoinAddress = *flagRaftJoinAddress

	cache := &rcache.Cached{
		Opts: opts,
		Log:  log.New(os.Stderr, "Cached: ", log.Ldate|log.Ltime),
		CM:   rcache.NewCacheManager(),
	}
	ctx := &rcache.CachedContext{cache}

	var l net.Listener
	var err error
	l, err = net.Listen("tcp", cache.Opts.HttpAddress)
	if err != nil {
		cache.Log.Fatal(fmt.Sprintf("listen %s failed: %s", cache.Opts.HttpAddress, err))
	}
	cache.Log.Printf("http server listen:%s", l.Addr())

	logger := log.New(os.Stderr, "httpserver: ", log.Ldate|log.Ltime)
	httpServer := rcache.NewHttpServer(ctx, logger)
	cache.HttpServer = httpServer

	go func() {
		http.Serve(l, httpServer.Mux)
	}()

	raft, err := rcache.NewRaftNode(cache.Opts, ctx)
	if err != nil {
		cache.Log.Fatal(fmt.Sprintf("new raft node failed:%v", err))
	}
	cache.Raft = raft

	if cache.Opts.JoinAddress != "" {
		err = rcache.JoinRaftCluster(cache.Opts)
		if err != nil {
			cache.Log.Fatal(fmt.Sprintf("join raft cluster failed:%v", err))
		}
	}

	for {
		select {
		case leader := <-cache.Raft.LeaderNotifyCh:
			if leader {
				cache.Log.Println("become leader, enable write api")
				cache.HttpServer.SetWriteFlag(true)
			} else {
				cache.Log.Println("become follower, close write api")
				cache.HttpServer.SetWriteFlag(false)
			}
		}
	}
}
