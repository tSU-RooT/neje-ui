/*
 * Copyright (c) 2016, Shinya Yagyu
 * All rights reserved.
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 * 3. Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived from this
 *    software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package webserver

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"golang.org/x/net/websocket"
)

//WebServer represents web server side.
type WebServer struct {
	client   *rpc.Client
	ch       chan struct{}
	Finished chan struct{}
}

//New starts web server and  browser, returns WebServer obj.
func New(bpath, firstPage string, strs ...interface{}) (*WebServer, error) {
	for _, str := range strs {
		if err := rpc.Register(str); err != nil {
			return nil, err
		}
	}
	w := &WebServer{
		ch: make(chan struct{}),
	}
	var err error
	addr := w.start()
	w.Finished, err = tryBrowser(bpath, addr.String()+"/"+firstPage)
	<-w.ch
	return w, err
}

//Close closes client RPC connection.
func (w *WebServer) Close() {
	w.ch <- struct{}{}
}

//Call calls calls RPC.
func (w *WebServer) Call(m string, args interface{}, reply interface{}) error {
	return w.client.Call(m, args, reply)
}

//regHandlers registers handlers to http.
func regHandlers(w *WebServer) {
	http.HandleFunc("/ws-server",
		func(w http.ResponseWriter, req *http.Request) {
			log.Println("connected to ws-server")
			s := websocket.Server{
				Handler: websocket.Handler(func(ws *websocket.Conn) {
					jsonrpc.ServeConn(ws)
				}),
			}
			s.ServeHTTP(w, req)
		})
	http.HandleFunc("/ws-client",
		func(rw http.ResponseWriter, req *http.Request) {
			log.Println("connected to ws-client")
			s := websocket.Server{
				Handler: websocket.Handler(func(ws *websocket.Conn) {
					w.client = jsonrpc.NewClient(ws)
					w.ch <- struct{}{}
					<-w.ch
				}),
			}
			s.ServeHTTP(rw, req)
		})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
}

//start starts webserver at localhost:0(i.e. free port) and returns listen address.
func (w *WebServer) start() net.Addr {
	regHandlers(w)
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := http.Serve(l, nil); err != nil {
			log.Fatal(err)
		}
	}()
	return l.Addr()
}
