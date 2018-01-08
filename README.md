# Gonyan
Gonyan is a simple stream based logging library for Go.

[![Go Report Card](https://goreportcard.com/badge/github.com/FedericoMaggi/gonyan)](https://goreportcard.com/report/github.com/FedericoMaggi/gonyan)&nbps;
[![Build Status](https://travis-ci.org/FedericoMaggi/gonyan.svg?branch=master)](https://travis-ci.org/FedericoMaggi/gonyan)&nbsp;
[![GoDoc](https://godoc.org/github.com/FedericoMaggi/gonyan?status.svg)](https://godoc.org/github.com/FedericoMaggi/gonyan)&nbsp;
[![codecov](https://codecov.io/gh/FedericoMaggi/gonyan/branch/master/graph/badge.svg)](https://codecov.io/gh/FedericoMaggi/gonyan)&nbsp;
[![GitHub issues](https://img.shields.io/github/issues/FedericoMaggi/gonyan.svg "GitHub issues")](https://github.com/FedericoMaggi/gonyan)



The idea behind the package is to allow the creation of a logging utility capable of sending machine readable logs to different targets (called _Streams_) based on the desired logging level.

### Example


```go
package main

import (
  "os"
  "github.com/FedericoMaggi/gonyan"
)

func main() {
  log := gonyan.NewLogger("GH-Example", nil, true)

  log.RegisterStream(gonyan.Fatal, os.Stderr)
  if verboseMode {
    log.RegisterStream(gonyan.Debug, os.Stdout)
  }
  
  log.Debug("Hello, World!")

  if err := doSomething(); err != nil {
    log.Errorf("doSomething() failed: %s", err.Error())
    return
  }

  log.Debug("doSomething() worked fine")
}
```
In this example two streams are registered for two different logging levels: `stdout` for `Debug` level and `stderr` for `Error` level. Log messages sent for the `Debug` level will be streamed only to `stdout` stream while those sent with `Error` level will be streamed only to `stderr` stream.

Please note that streams are completely optional, you can call all logging functions even without registering any stream, your program won't crash but won't even log anything (of course).

Also, each logging level can have multiple stream registered and custom streams can be defined and used leaving you the power to decide what log, how, where and when.

## Stream

> Yeah, but what is a stream?

A stream is whatever `struct` that implements the function `Write([]byte) (int, error)` this choice allows Gonyan to natively support many `I/O` structures (e.g. `File`, `bytes.Buffer`, `bufio.Writer`, etc..) and being agnostic regarding where the log will be actually used. 

With time, many streams will be provided out-of-the-box but everyone can create its own custom stream object and transparently provide it to the Gonyan logger.
 
### Formatting

As of now log messages are streamed as JSON strings, better support for other/custom formats has to be defined.

Example: 
```
{"tag":"GH-Example","timestamp":1515161633123,"message":"Hello, World!"}
```
