# bb

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
curl -X PUT --data-binary @myfile.jpg localhost:8080/myfile
curl localhost:8080/myfile > myfile.jpg
```
