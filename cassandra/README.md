# cassandra

Cassandra/GoCQL utilities.

## LoggingConnectObserver

A [ConnectObserver][1] implementation for [GoCQL][2] that uses [servicelib-golang/logger][3] to
log new connections.  Successful connections are logged at `DEBUG`, and failures at `ERROR`.

```golang
package main

import (
    "os"

    "gerrit.wikimedia.org/r/mediawiki/services/servicelib-golang/cassandra"
    log "gerrit.wikimedia.org/r/mediawiki/services/servicelib-golang/logger"
    "github.com/gocql/gocql"
)

func main() {
    logger _ := log.NewLogger(os.Stdout, "sessionstore", "INFO")
    cluster := gocql.NewCluster("192.168.1.1", "192.168.1.2", "192.168.1.3")
    cluster.ConnectObserver = &LoggingConnectObserver{Logger: logger}

    session, err := cluster.CreateSession()
    //...
}
```

[1]: https://pkg.go.dev/github.com/gocql/gocql?utm_source=godoc#ConnectObserver
[2]: https://github.com/gocql/gocql
[3]: https://gerrit.wikimedia.org/g/mediawiki/services/servicelib-golang
