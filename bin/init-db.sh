#!/usr/bin/env bash

docker exec -it scylla bash << EOF
  cqlsh;
  CREATE KEYSPACE todos WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
EOF