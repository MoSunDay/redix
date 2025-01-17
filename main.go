// Copyright 2018 The Redix Authors. All rights reserved.
// Use of this source code is governed by a Apache 2.0
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"strconv"

	"github.com/MoSunDay/go-color"
	_ "github.com/MoSunDay/redix/rcache"
)

func main() {
	fmt.Println(color.MagentaString(redixBrand))
	fmt.Printf("⇨ redix server version: %s \n", color.GreenString(redixVersion))
	fmt.Printf("⇨ redix selected engine: %s \n", color.GreenString(*flagEngine))
	fmt.Printf("⇨ redix workers count: %s \n", color.GreenString(strconv.Itoa(*flagWorkers)))
	fmt.Printf("⇨ redix resp server available at: %s \n", color.GreenString(*flagRESPListenAddr))
	fmt.Printf("⇨ redix http server available at: %s \n", color.GreenString(*flagHTTPListenAddr))
	fmt.Printf("⇨ raft http server available at: %s \n", color.GreenString(*flagHttpAddress))
	fmt.Printf("⇨ raft tcp available at: %s \n", color.GreenString(*flagRaftTCPAddress))

	err := make(chan error)

	go (func() {
		err <- initRespServer()
	})()

	go (func() {
		err <- initHTTPServer()
	})()

	go (func() {
		runRatfServer()
	})()

	if err := <-err; err != nil {
		color.Red(err.Error())
	}
}
