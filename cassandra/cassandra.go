package cassandra

import (
	"github.com/gocql/gocql"
	"log"
)

const DefaultURL = "127.0.0.1"

type Cassandra struct {
	cluster    *gocql.ClusterConfig
	Connection *gocql.Session
}

func Connect(url string, keyspace string) (*Cassandra, error) {
	var err error
	con := &Cassandra{}
	con.cluster = gocql.NewCluster(url)
	con.cluster.Keyspace = keyspace
	if err != nil {
		return nil, err
	}

	con.Connection, err = con.cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	log.Println("db initialized")

	return con, nil
}

func (c *Cassandra) CreateTable(table string) error {
	if err := c.Connection.Query(table).RetryPolicy(nil).Exec(); err != nil {
		log.Printf("error creating table table=%q err=%v\n", table, err)
		return err
	}

	return nil
}

func (c *Cassandra) Close() {
	log.Println("db connection closed")
	c.Connection.Close()
}