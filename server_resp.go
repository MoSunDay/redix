// Copyright 2018 The Redix Authors. All rights reserved.
// Use of this source code is governed by a Apache 2.0
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/MoSunDay/go-color"
	"github.com/MoSunDay/redcon"
	"github.com/MoSunDay/redix/hash"
)

func initRespServer() error {
	return redcon.ListenAndServe(
		*flagRESPListenAddr,
		func(conn redcon.Conn, cmd redcon.Command) {
			defer (func() {
				if err := recover(); err != nil {
					conn.WriteError(fmt.Sprintf("fatal error: %s", (err.(error)).Error()))
				}
			})()

			todo := strings.TrimSpace(strings.ToLower(string(cmd.Args[0])))
			args := []string{}
			for _, v := range cmd.Args[1:] {
				v := strings.TrimSpace(string(v))
				args = append(args, v)
			}

			if *flagVerbose {
				log.Println(color.YellowString(todo), color.CyanString(strings.Join(args, " ")))
			}

			db, err := selectDB("0")
			if err != nil {
				conn.WriteError(fmt.Sprintf("db error: %s", err.Error()))
				return
			}

			// our internal change log
			if changelog.Subscribers(defaultPubSubAllTopic) > 0 {
				changelog.Broadcast(Change{
					Command:   todo,
					Arguments: args,
				}, defaultPubSubAllTopic)
			}

			if todo == "ping" {
				conn.WriteString("PONG")
				return
			}

			if todo == "quit" {
				conn.WriteString("OK")
				conn.Close()
				return
			}

			for _, commands := range h_commands {
				fn := commands[todo]
				if nil == fn {
					continue
				}

				slot := hash.GetSlotNumber(args[0])
				address := SlotCache.CM.Get("0")
				if address != *flagRESPListenAddr {
					conn.WriteError(fmt.Sprintf("MOVED %d %s", slot, address))
				}
				args[0] = strconv.Itoa(int(slot)) + "/" + args[0]
				fn(Context{
					Conn:   conn,
					action: todo,
					args:   args,
					db:     db,
				})
				return
			}

			fn := p_commands[todo]
			if fn != nil {
				fn(Context{
					Conn:   conn,
					action: todo,
					args:   args,
					db:     db,
				})
				return
			}
			conn.WriteError(fmt.Sprintf("unknown commands [%s]", todo))
		},
		func(conn redcon.Conn) bool {
			conn.SetContext(map[string]interface{}{})
			return true
		},
		nil,
	)
}
