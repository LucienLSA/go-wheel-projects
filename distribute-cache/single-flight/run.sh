#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8001 &
./server -port=8002 &
./server -port=8003 -api=1 &

sleep 2
echo ">>> start test"

# 默认并发数 10，可通过第一个参数指定
CONCURRENCY=${1:-10}
for i in $(seq 1 $CONCURRENCY); do
    curl "http://localhost:9999/api?key=Tom" &
done

wait