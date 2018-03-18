# icii

[![Go Report Card](https://goreportcard.com/badge/broadcastle.co/code/icii)](https://goreportcard.com/report/broadcastle.co/code/icii)
[![GoDoc](https://godoc.org/broadcastle.co/code/icii?status.svg)](https://godoc.org/broadcastle.co/code/icii)

Stream MP3 files to icecast with go.

## Installation

_MOSTLY EMPTY_

```bash
go get broadcastle.co/code/icii
```

## Usage

```go
stream, err := icii.CreateStream("Station", "We play audio!", "10.10.10.42:8080", "/stream.mp3")
if err != nil {
    log.Fatal(err)
}

// Stream a using a filepath.
if err := stream.File("/path/to/audio.mp3"); err != nil {
    log.Error(err)
}

// io.Reader can also be used.

data, err := os.Open("/path/to/music.mp3")
if err != nil {
    log.Error(err)
}

if err := stream.Reader(data); err != nil {
    log.Error(err)
}

```
