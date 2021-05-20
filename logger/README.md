# logger

## Usage

```
$ go get github.com/eevans/servicelib-golang/logger
```

```
package main

import (
    "os"
    "time"

    "github.com/eevans/servicelib-golang/logger"
)

func main() {
    log, _ := logger.NewLogger(os.Stdout, "sessionstore", "kask", logger.INFO)

    // The basics...
    log.Debug("Debugging yer bugs")
    log.Info("The current time is %s", time.Now().Format(time.RFC3339))

    // Using a request-scoped logger...
    hostname, _ := os.Hostname()
    log.Request().
        Trace("0a762a9c-b8d6-11eb-87bc-4f82287279b0").
        ClientIP("127.0.0.1").
        ClientPort(9000).
        ClientBytes(1500).
        Log(logger.INFO, "request received by %s", hostname)

}
```
