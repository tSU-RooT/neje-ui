[![Build Status](https://travis-ci.org/utamaro/wsrpc.svg?branch=master)](https://travis-ci.org/utamaro/wsrpc)
[![GoDoc](https://godoc.org/github.com/utamaro/wsrpc?status.svg)](https://godoc.org/github.com/utamaro/wsrpc)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/utamaro/wsrpc/master/LICENSE)
[![Coverage Status](https://coveralls.io/repos/utamaro/wsrpc/badge.svg?branch=master)](https://coveralls.io/r/utamaro/wsrpc?branch=master)


# wsrpc

## Overview

This is a library to communicate with browser by JSON-RPC on websocket using
[gopherjs](https://github.com/gopherjs/gopherjs).

You can call funcs on browser from server (and vice versa) in [golang RPC-style](https://golang.org/pkg/net/rpc/) without considering about websocket and javascript.

I made this library to make GUI by html5 on browser easily.


## Requirements

This requires

* git
* go 1.4+


## Installation

    $ go get -u github.com/utamaro/wsrpc


## Example
(This example omits error handlings for simplicity.)

## browser side

ex.go

```go
type Args struct {
	A int
	B int
	C string
}

type GUI struct{}

//func to be called from web server
func (g *GUI) Write(args *Args, reply *int) error {
	//show welcome message:
	jQuery(OUTPUT).SetText(args.C)
	return nil
}

func main() {
	b, _ := browser.New("7000", new(GUI))
	jQuery(INPUT).On(jquery.KEYUP, func(e jquery.Event) {
		go func() {
			args := Args{A: 17, B: 8}
			var reply int
//call func in web server from browser 
			b.Call("Arith.Multiply", args, &reply)
			jQuery("#output").SetText(strconv.Itoa(reply))
		}()
	})
}
```

Then compile it by gopherjs to create ex.js:

ex.html
```
gopherjs build ex.go
```


```html
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Welcome to GopherJS with jQuery</title>
    <script src="https://code.jquery.com/jquery-2.2.0.min.js"></script>
</head>
<body>
    <span id="output"></span>
    <script src="ex.js"></script>
</body>
</html>
```

## webserver side

```go
type Args struct {
	A int
	B int
	C string
}

type Arith struct{}

//func to be called from browser
func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	ws, _ := webserver.New("localhost:7000", new(Arith))
	var reply int
//call func in browswer from webserver 
	ws.Call("GUI.Write", &Args{C: "test"}, &reply)
}
```


# Contribution
Improvements to the codebase and pull requests are encouraged.


