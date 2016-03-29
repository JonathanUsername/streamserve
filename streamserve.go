package main

// This package simply takes a streaming input and serves it on a port

import (
  "bufio"
  "os"
  "io"
  "log"
  "flag"
  "fmt"
  "path/filepath"
  "strconv"
  "net/http"
)

var portStr string
var pathStr string

func check (err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func handleReq (w http.ResponseWriter, r *http.Request) {
  fmt.Println("Request from", r.RemoteAddr)
  f, err := os.Open(pathStr)
  check(err)

  reader := bufio.NewReader(f)
  for {
    buf := make([]byte, 0, 4*1024)
    v, err := reader.Read(buf[:cap(buf)])
    buf = buf[:v]
    if v == 0 {
      if err == nil {
        continue
      }
      if err == io.EOF {
          break
      }
      log.Fatal(err)
    }
    w.Write(buf)
  }
}

func listen () {
  startMsg := fmt.Sprintf("Serving %s at 0.0.0.0%s", pathStr, portStr)
  fmt.Println(startMsg)
  http.HandleFunc("/", handleReq)
  http.ListenAndServe(portStr, nil)
}

func main () {

  flag.Parse()
  pathName := flag.Arg(0)
  portName := flag.Arg(1)

  if pathName == "" {
    fmt.Fprintf(os.Stderr, "No path supplied")
    return
  } else {
    pathStr = filepath.FromSlash(pathName)
  }

  if portName == "" {
    portStr = ":7777"
  } else {
    _, err := strconv.Atoi(portName)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error parsing port number")
      return
    } else {
      portStr = fmt.Sprintf(":%s", portName)
    }
  }

  listen()
}