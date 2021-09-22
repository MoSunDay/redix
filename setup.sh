#!/bin/sh
rm -fr raft-data
./redix -node node1 -bootstrap true &
sleep 5
./redix -node node2 -join 127.0.0.1:6000 -raft-http 127.0.0.1:6001 -resp-addr :6667 -http-addr :7091 -raft 127.0.0.1:7001 &
./redix -node node3 -join 127.0.0.1:6000 -raft-http 127.0.0.1:6002 -resp-addr :6668 -http-addr :7092 -raft 127.0.0.1:7002 &
