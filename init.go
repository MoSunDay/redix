// Copyright 2018 The Redix Authors. All rights reserved.
// Use of this source code is governed by a Apache 2.0
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/MoSunDay/go-color"
	"github.com/MoSunDay/go-pubsub"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
)

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(*flagWorkers)

	*flagStorageDir = *flagStorageDir + "/" + *flagRaftNode

	RaftAddrSlice := strings.Split(*flagRaftTCPAddress, ":")
	port, err := strconv.Atoi(RaftAddrSlice[1])
	if err != nil {
		log.Fatal(err)
	}

	*flagRESPListenAddr = RaftAddrSlice[0] + ":" + strconv.Itoa(port+200)
	*flagHTTPListenAddr = RaftAddrSlice[0] + ":" + strconv.Itoa(port+400)
	*flagHttpAddress = RaftAddrSlice[0] + ":" + strconv.Itoa(port+600)

	if !*flagVerbose {
		logger := logrus.New()
		logger.SetOutput(ioutil.Discard)
	}

	if !supportedEngines[*flagEngine] {
		fmt.Println(color.RedString("Invalid strorage engine specified"))
		os.Exit(0)
		return
	}

	databases = new(sync.Map)
	changelog = pubsub.NewBroker()
	webhooks = new(sync.Map)
	websockets = new(sync.Map)
	engineOptions = (func() url.Values {
		opts, _ := url.ParseQuery(*flagEngineOptions)
		return opts
	})()

	snowflakenode, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(color.RedString(err.Error()))
		os.Exit(0)
		return
	}

	snowflakeGenerator = snowflakenode
}
