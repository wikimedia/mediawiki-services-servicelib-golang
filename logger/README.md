# logger

## Usage

```
$ go get github.com/eevans/servicelib-golang/logger
```

```
package main

import (
    "os"
    "github.com/eevans/servicelib-golang/logger"
)

func main() {
    var log logger.Logger = logger.NewLogger(os.Stdout, "sessionstore", "kask", logger.INFO)

    log.DEBUG("Debugging yer bugs")
    log.INFO("The current time is %s", time.Now().Format(time.RFC3339))
    log.TraceID("0a762a9c-b8d6-11eb-87bc-4f82287279b0").Log(logger.INFO, "request received by %s", os.Hostname())

}
```
