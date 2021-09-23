// Copyright 2018 The Redix Authors. All rights reserved.
// Use of this source code is governed by a Apache 2.0
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"net/url"
	"runtime"
	"sync"

	"github.com/MoSunDay/go-pubsub"
	rcache "github.com/MoSunDay/redix/rcache"
	"github.com/bwmarrin/snowflake"
)

var (
	flagRESPListenAddr  = flag.String("resp-addr", "", "the address of resp server, raft server port + 200")
	flagHTTPListenAddr  = flag.String("http-addr", "", "the address of the http server, raft server port + 400")
	flagStorageDir      = flag.String("data-storage", "./redix-data", "the data storage directory")
	flagEngine          = flag.String("engine", "leveldb", "the storage engine to be used")
	flagEngineOptions   = flag.String("engine-options", "", "options related to used engine in the url query format, i.e (opt1=val2&opt2=val2)")
	flagWorkers         = flag.Int("workers", runtime.NumCPU(), "the default workers number")
	flagVerbose         = flag.Bool("verbose", false, "verbose or not")
	flagACK             = flag.Bool("ack", true, "acknowledge write or return immediately")
	flagHttpAddress     = flag.String("raft-http", "", "Http address, raft server port + 600")
	flagRaftTCPAddress  = flag.String("raft", "127.0.0.1:7000", "raft tcp address")
	flagRaftNode        = flag.String("node", "node1", "raft node name")
	flagRaftBootstrap   = flag.Bool("bootstrap", false, "start as raft cluster")
	flagRaftJoinAddress = flag.String("join", "", "join address for raft cluster")
	flageRaftDataDir    = flag.String("raft-storage", "./raft-data", "the raft storage directory")
)

var (
	databases          *sync.Map
	changelog          *pubsub.Broker
	webhooks           *sync.Map
	websockets         *sync.Map
	snowflakeGenerator *snowflake.Node
	kvjobs             chan func()
	SlotCache          *rcache.Cached
)

var (
	r_commands = map[string]CommandHandler{
		// strings
		"get":  getCommand,
		"mget": mgetCommand,

		// list
		"lrange": lrangeCommand,
		"lcount": lcountCommand,

		// sets (list alias)
		"smembers": lrangeCommand,
		"srem":     lremCommand,
		"scard":    lcountCommand,
		"sscan":    lrangeCommand,

		// hashes
		"hget":    hgetCommand,
		"hgetall": hgetallCommand,
		"hkeys":   hkeysCommand,
		"hexists": hexistsCommand,
		"hlen":    hlenCommand,
	}

	w_commands = map[string]CommandHandler{
		// strings
		"set":  setCommand,
		"mset": msetCommand,

		"del":    delCommand,
		"exists": existsCommand,
		"incr":   incrCommand,
		"ttl":    ttlCommand,

		// lists
		"lpush":      lpushCommand,
		"lpushu":     lpushuCommand,
		"lrem":       lremCommand,
		"lcard":      lcountCommand,
		"lsum":       lsumCommand,
		"lavg":       lavgCommand,
		"lmin":       lminCommand,
		"lmax":       lmaxCommand,
		"lsrch":      lsearchCommand,
		"lsrchcount": lsearchcountCommand,

		// sets (list alias)
		"sadd": lpushuCommand,

		// hashes
		"hset":  hsetCommand,
		"hdel":  hdelCommand,
		"hmset": hmsetCommand,
		"hincr": hincrCommand,
		"httl":  httlCommand,

		// pubsub
		"publish":        publishCommand,
		"subscribe":      subscribeCommand,
		"webhookset":     webhooksetCommand,
		"webhookdel":     webhookdelCommand,
		"websocketopen":  websocketopenCommand,
		"websocketclose": websocketcloseCommand,
	}

	p_commands = map[string]CommandHandler{
		// utils
		"keys": keysCommand,

		"encode":   encodeCommand,
		"uuidv4":   uuid4Command,
		"uniqid":   uniqidCommand,
		"randstr":  randstrCommand,
		"randint":  randintCommand,
		"time":     timeCommand,
		"dbsize":   dbsizeCommand,
		"gc":       gcCommand,
		"info":     infoCommand,
		"echo":     echoCommand,
		"flushdb":  flushdbCommand,
		"flushall": flushallCommand,
		"cluster":  clusterCommand,

		// ratelimit
		"ratelimitset":  ratelimitsetCommand,
		"ratelimittake": ratelimittakeCommand,
		"ratelimitget":  ratelimitgetCommand,
	}

	h_commands = []map[string]CommandHandler{r_commands, w_commands}
)

var (
	supportedEngines = map[string]bool{
		"boltdb":  true,
		"leveldb": true,
		"null":    true,
		"sqlite":  true,
	}
	engineOptions         = url.Values{}
	defaultPubSubAllTopic = "*"
)

const (
	redixVersion = "2.00-dev"
	redixBrand   = `

		 _______  _______  ______  _________         
		(  ____ )(  ____ \(  __  \ \__   __/|\     /|
		| (    )|| (    \/| (  \  )   ) (   ( \   / )
		| (____)|| (__    | |   ) |   | |    \ (_) / 
		|     __)|  __)   | |   | |   | |     ) _ (  
		| (\ (   | (      | |   ) |   | |    / ( ) \ 
		| ) \ \__| (____/\| (__/  )___) (___( /   \ )
		|/   \__/(_______/(______/ \_______/|/     \|

A high-concurrency standalone NoSQL datastore with the support for redis protocol 
and multiple backends/engines, also there is a native support for
real-time apps via webhook & websockets besides the basic redis channels.

	`
)
