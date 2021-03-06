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
  "regexp"
)

var portStr string
var pathStr string
var endpoint string

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
  startMsg := fmt.Sprintf("Serving %s at 0.0.0.0%s%s", pathStr, portStr, endpoint)
  fmt.Println(startMsg)
  http.HandleFunc(endpoint, handleReq)
  http.ListenAndServe(portStr, nil)
}

func main () {
  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
    fmt.Fprintf(os.Stderr, " %s [OPTIONS] FILEPATH PORTNUMBER\n\nOptions:\n", os.Args[0])
    flag.PrintDefaults()
  }
  flag.StringVar(&endpoint, "e", "/", "Specify a particular endpoint.")

  flag.Parse()

  pathName := flag.Arg(0)
  portName := flag.Arg(1)

  if pathName == "" {
    fmt.Fprintf(os.Stderr, "No path supplied.")
    return
  } else {
    pathStr = filepath.FromSlash(pathName)
  }

  if portName == "" {
    portStr = ":7777"
  } else {
    _, err := strconv.Atoi(portName)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error parsing port number.")
      return
    } else {
      portStr = fmt.Sprintf(":%s", portName)
    }
  }

  match, _ := regexp.MatchString("^/", endpoint)

  if !match {
    fmt.Fprintf(os.Stderr, "Endpoint '%s' invalid. Must begin with slash.", endpoint)
    return
  }

  listen()
}