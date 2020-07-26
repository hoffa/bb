# bb

[![Go Report Card](https://goreportcard.com/badge/github.com/hoffa/bb)](https://goreportcard.com/report/github.com/hoffa/bb)

Easily store small amounts of data behind an HTTP server.

Built so I can have something remotely useful to play with on my Raspberry Pi Zero W.

## Install

```shell
go get github.com/hoffa/bb
```

## Usage

```shell
bb
```

## Example

```shell
curl --data-binary @file.txt localhost:8080/file
curl localhost:8080/file
```
