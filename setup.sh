#!/bin/sh
pkill redix
sleep 2
rm -fr raft-data
go build .
./redix -node node1 -bootstrap true &
sleep 5
./redix -node node2 -join 127.0.0.1:7600 -raft 127.0.0.1:7001 &
./redix -node node3 -join 127.0.0.1:7600 -raft 127.0.0.1:7002 &
