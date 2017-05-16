# TODO Service
Simple service to manage a TODO list

## Requirements
- docker

## Build the service

To build the docker image of the service, run the following command:  
`make build`

## Start the service

### Start Dependencies

Start NATS:      
`bin/run-nats-streaming.sh`

Start the scylla database:  
`bin/run-scylla.sh`

#### Prepare the database (create keyspace)

```bash  
docker exec -it scylla sh;
cqlsh;
CREATE KEYSPACE todos WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
```

### Start the service

`bin/run-todo-service.sh`

The API should then be available at:  
`127.0.0.1:8011/api/v1/todos`

## API
The service exposes the following API resources:

Get a list of TODOs:  
`GET /api/v1/todos`

Create a TODO:  
`POST /api/v1/todos`

Get a single TODO:  
`GET /api/v1/todos/{id}`

Remove a TODO:  
`DELETE /api/v1/todos/{id}`

