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

    log.Debug("Debugging yer bugs")
    log.Info("The current time is %s", time.Now().Format(time.RFC3339))

    hostname, _ := os.Hostname()
    log.TraceID("0a762a9c-b8d6-11eb-87bc-4f82287279b0").Log(logger.INFO, "request received by %s", hostname)

}
```
