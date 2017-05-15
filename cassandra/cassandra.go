package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/Sirupsen/logrus"
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

	return con, nil
}

func (c *Cassandra) CreateTable(table string) error {
	if err := c.Connection.Query(table).RetryPolicy(nil).Exec(); err != nil {
		logrus.Errorf("error creating table table=%q err=%v\n", table, err)
		return err
	}

	return nil
}

func (c *Cassandra) Close() {
	logrus.Info("db connection closed")
	c.Connection.Close()
}