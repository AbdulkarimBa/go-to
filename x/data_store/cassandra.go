package data_store

import (
	"github.com/gocql/gocql"
	log "github.com/shyyawn/go-to/x/logging"
	"github.com/shyyawn/go-to/x/source"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type Cassandra struct {
	cluster           *gocql.ClusterConfig
	session           *gocql.Session
	lock              *sync.RWMutex
	Keyspace          string   `mapstructure:"keyspace"`
	Hosts             []string `mapstructure:"hosts"`
	Timeout           int      `mapstructure:"timeout"`
	ReconnectInterval int      `mapstructure:"reconnect_interval"`
	Port              int      `mapstructure:"port"`
}

func (ds *Cassandra) LoadFromConfig(key string, config *viper.Viper) error {
	ds.lock = &sync.RWMutex{}
	return source.LoadFromConfig(key, config, ds)
}

func (ds *Cassandra) Init() error {
	ds.lock = &sync.RWMutex{}
	_ = ds.Cluster()
	return nil
}

func (ds *Cassandra) Session() *gocql.Session {
	if ds.session == nil || ds.session.Closed() {
		if ds.cluster == nil {
			_ = ds.Cluster()
		}
		ds.lock.Lock()
		defer ds.lock.Unlock()
		var err error
		ds.session, err = ds.cluster.CreateSession()
		if err != nil {
			log.Error("Unable to connect to cassandra", err)
		}
	}
	return ds.session
}

func (ds *Cassandra) Cluster() *gocql.ClusterConfig {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	ds.cluster = gocql.NewCluster(ds.Hosts...)

	ds.cluster.Port = ds.Port
	ds.cluster.Keyspace = ds.Keyspace
	ds.cluster.Timeout = time.Second * time.Duration(ds.Timeout)
	ds.cluster.ReconnectInterval = time.Second * time.Duration(ds.ReconnectInterval)
	ds.cluster.Consistency = gocql.LocalQuorum
	return ds.cluster
}