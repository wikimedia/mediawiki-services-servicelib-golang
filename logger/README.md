# logger

A very simple golang logging library for k8s deployed microservices at the Wikimedia Foundation.

## Features

- JSON-formatted messages (one per line)
- Logs messages only in [Elasic Common Schema](https://doc.wikimedia.org/ecs/) format
- No external dependencies

## Usage

```
$ go get github.com/eevans/servicelib-golang/logger
```

```golang
package main

import (
    "math/rand"
    "os"
    "time"

    "github.com/eevans/servicelib-golang/logger"
)

func main() {
    log, _ := logger.NewLogger(os.Stdout, "sessionstore", logger.INFO)

    // The basics...
    log.Info("Random number %d is random", rand.Intn(100))

    // Using a request-scoped logger inside an http handler
    handler := func(w http.ResponseWriter, r *http.Request) {
        hostname, err := os.Hostname()
        if err != nil {
            log.Error("Oh no; Failed to get the hostname: %s", err)
            return
        }

        // Inside an http handler...
        log.Request(r).Log(logger.INFO, "request received by %s", hostname)
        io.WriteString(w, "Hello World!")
    }
}
```
