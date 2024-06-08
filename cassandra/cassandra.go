package cassandra

import (
	log "gerrit.wikimedia.org/r/mediawiki/services/servicelib-golang/logger"

	"github.com/gocql/gocql"
)

// LoggingConnectObserver is a [ConnectObserver] implementation for [GoCQL].  It uses
// [servicelib-golang/logger] to log new connection events (success logged at DEBUG
// level, failures at ERROR).
//
// Example usage:
//
//	logger, _ = log.NewLogger(os.Stdout, "service", "debug")
//	cluster := gocql.NewCluster("192.168.1.1", "192.168.1.2", "192.168.1.3")
//	...
//	cluster.ConnectObserver = &LoggingConnectObserver{Logger: logger}
//	session := cluster.CreateSession()
//
// [ConnectObserver]: https://pkg.go.dev/github.com/gocql/gocql?utm_source=godoc#ConnectObserver
// [GoCQL]: https://github.com/gocql/gocql
// [servicelib-golang/logger]: https://gerrit.wikimedia.org/g/mediawiki/services/servicelib-golang
type LoggingConnectObserver struct {
	Logger *log.Logger
}

func (o *LoggingConnectObserver) ObserveConnect(conn gocql.ObservedConnect) {
	if conn.Err != nil {
		o.Logger.Error("Cassandra: Problem connecting to %s, (%s)", conn.Host.ConnectAddress(), conn.Err)
		return
	}

	o.Logger.Debug("Cassandra: Opened new connection to %s", conn.Host.ConnectAddress())
}
