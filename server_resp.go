// Copyright 2018 The Redix Authors. All rights reserved.
// Use of this source code is governed by a Apache 2.0
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/MoSunDay/go-color"
	"github.com/MoSunDay/redcon"
)

func initRespServer() error {
	return redcon.ListenAndServe(
		*flagRESPListenAddr,
		func(conn redcon.Conn, cmd redcon.Command) {
			// handles any panic
			defer (func() {
				if err := recover(); err != nil {
					conn.WriteError(fmt.Sprintf("fatal error: %s", (err.(error)).Error()))
				}
			})()

			ctx := (conn.Context()).(map[string]interface{})
			todo := strings.TrimSpace(strings.ToLower(string(cmd.Args[0])))
			args := []string{}
			for _, v := range cmd.Args[1:] {
				v := strings.TrimSpace(string(v))
				args = append(args, v)
			}

			if *flagVerbose {
				log.Println(color.YellowString(todo), color.CyanString(strings.Join(args, " ")))
			}

			if ctx["db"] == nil || ctx["db"].(string) == "" {
				ctx["db"] = "0"
			}

			db, err := selectDB(ctx["db"].(string))
			if err != nil {
				conn.WriteError(fmt.Sprintf("db error: %s", err.Error()))
				return
			}

			// our internal change log
			if changelog.Subscribers(defaultPubSubAllTopic) > 0 {
				changelog.Broadcast(Change{
					Namespace: ctx["db"].(string),
					Command:   todo,
					Arguments: args,
				}, defaultPubSubAllTopic)
			}

			// internal ping-pong
			if todo == "ping" {
				conn.WriteString("PONG")
				return
			}

			// close the connection
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

				slot := crc16sum(args[0]) % 16384

				if *flagRaftNode == "node1" {
					if slot <= 5461 {
					} else if slot >= 10923 {
						conn.WriteError(fmt.Sprintf("MOVED %d 0.0.0.0:6668", slot))
					} else {
						conn.WriteError(fmt.Sprintf("MOVED %d 0.0.0.0:6667", slot))
					}
				}

				if *flagRaftNode == "node2" {
					if slot >= 5462 || slot <= 10922 {
					} else if slot >= 10923 {
						conn.WriteError(fmt.Sprintf("MOVED %d 0.0.0.0:6668", slot))
					} else {
						conn.WriteError(fmt.Sprintf("MOVED %d 0.0.0.0:6666", slot))
					}
				}

				if *flagRaftNode == "node3" {
					if slot >= 10923 || slot <= 16383 {
					} else if slot <= 5461 {
						conn.WriteError(fmt.Sprintf("MOVED %d 0.0.0.0:6666", slot))
					} else {
						conn.WriteError(fmt.Sprintf("MOVED %d 0.0.0.0:6667", slot))
					}
				}

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
