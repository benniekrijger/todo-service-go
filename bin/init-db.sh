#!/usr/bin/env bash

docker exec -it scylla sh << EOF
  cqlsh;
  CREATE KEYSPACE todos WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
  use todos;
  create table todos (
    id UUID,
    title text,
    completed boolean,
    PRIMARY KEY(id)
  );
EOF