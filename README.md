# Gonyan
Gonyan is a simple stream based logging library for Go.

The basic idea behind the package is to allow the usage of different logging targets (called _Streams_) and use them associated with the different log level messages.

### Example

```go
package main

import (
  "os"
  "github.com/FedericoMaggi/gonyan"
)

func main() {
  log := gonyan.NewLogger("GH-Example", nil, true)
  log.RegisterStream(gonyan.Debug, os.Stdout)
  
  log.Debug("Hello, World!")
}
```

This will send the message `Hello, World!` to all the registered `Debug` streams, in the example only stdout is registered causing the log to be printed in the standard output.

Each log level can have multiple stream registered and custom streams can be defined and used leaving you the power to decide what log, where and when.

### Formatting

As of now log messages are streamed as JSON strings, better support for other/custom formats has to be defined.

Example: 
```
{"tag":"GH-Example","timestamp":1515161633123,"message":"Hello, World!"}
```
