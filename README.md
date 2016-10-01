[![Build Status](https://travis-ci.org/utamaro/neje-ui.svg?branch=master)](https://travis-ci.org/utamaro/neje-ui)
[![GoDoc](https://godoc.org/github.com/utamaro/neje-ui?status.svg)](https://godoc.org/github.com/utamaro/neje-ui)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/utamaro/neje-ui/master/LICENSE)


# neje-ui

Not Embed, Just Execute chrome browser for UI in golang.

For now just a PoC(proof of concept).  Don't believe so much :)

## Overview

This library is an UI alternative for golang by using installed chrome browser(or something else) that is already installed.

This communicates with browser by JSON-RPC on websocket using [gopherjs](https://github.com/gopherjs/gopherjs).

You can call funcs on browser from server (and vice versa) in [golang RPC-style](https://golang.org/pkg/net/rpc/) 
without considering about websocket and javascript.

You can write server-side program and client side program in golang.

## Requirements

This requires

* git
* go 1.6+
* gopherjs
```
go get -u github.com/gopherjs/gopherjs
```

## Installation

    $ go get -u github.com/utamaro/neje-uis


## Example
(This example omits error handlings for simplicity.)

## browser side

[ex.go](https://github.com/utamaro/wsrpc/blob/master/example/browser/ex.go)

```go

//GUI is struct to bel called from remote by rpc.
type GUI struct{}

//Write writes a response from the server.
func (g *GUI) Write(msg *string, reply *string) error {
	//show welcome message:
	jQuery("#from_server").SetText(msg)
	return nil
}

func main() {
	b,_ := browser.New(new(GUI))
	jQuery("button").On(jquery.CLICK, func(e jquery.Event) {
		go func() {
			m := jQuery("#to_server").Val()
			response := ""
			b.Call("Msg.Message", &m, &response)
			//show welcome message:
			jQuery("#response").SetText(response)
		}()
	})

}


```

Then compile it by gopherjs to create ex.js:

```
go get  
gopherjs build ex.go
```

## webserver side

[ex.go](https://github.com/utamaro/wsrpc/blob/master/example/webserver/ex.go)

```go

//Msg  is struct to bel called from remote by rpc.
type Msg struct{}

//Message writes a message to the browser.
func (t *Msg) Message(m *string, response *string) error {
	*response = "OK, I heard that you said\"" + *m + "\""
	return nil
}

func main() {
	ws,_ := webserver.New("", "ex.html", new(Msg))

	for {
		select {
		case <-ws.Finished:
			log.Println("browser was closed. Exiting...")
			return
		case <-time.After(10 * time.Second):
			msg := "Now " + time.Now().String() + " at server!"
			reply := ""
			ws.Call("GUI.Write", &msg, &reply)
		}
	}
}

```

Then copy ex.html and ex.js to the webserver directory,
```
go run ex.go
```

Then your chrome browser (or something else if chrome is not installed) automatically is opened and
display the demo.

## What is It Doing?

1. Start a webserver including websocket at a free port at localhost.
1. Search chrome browser. 
	1. If windows, read registry to get the path to chrome. 
	2. If macos, just run "open -a google chrome ".
    3. if linux, run "google-chrome", "chrome", or something else.
1. If chrome is found, run chrome with the options "--disable-extension --app=<url>"
1. if chrome is not found, 
	1. If windows, just run "start <url>". 
	2. If macos, just run "open <url>  ".
    3. if linux, just run "xdg-open <url>"
1. communicate between webserver and browser using websocket.

## Why not embed chrome lib?

1. chrome lib is very big(about 100MB?) for single apps.
2. chrome lib APIs are always changing.
3. Not want to loose eco system(easy to cross compile etc) of golang.
4. Chrome lib is too difficult to understand for me :( .
5. Chrome browser has convinent options for application (--app etc).

### Pros
 Can make golang progs that can be cross compiled easily with small size.

### Cons
 Cann't control browser precisely, must control them by javascript manually. (window size, menu etc.)


# Contribution
Improvements to the codebase and pull requests are encouraged.


