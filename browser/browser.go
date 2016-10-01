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

package browser

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
)

//Browser represents RPC of browser side
type Browser struct {
	Client *rpc.Client
	s      net.Conn
	c      net.Conn
}

//New connects websocket and returns Browser obj.
func New(strs ...interface{}) (*Browser, error) {
	var err error
	b := &Browser{}
	for _, str := range strs {
		if errr := rpc.Register(str); errr != nil {
			return nil, errr
		}
	}
	port := js.Global.Get("window").Get("location").Get("port").String()
	b.s, err = websocket.Dial("ws://localhost:" + port + "/ws-client") // Blocks until connection is established
	if err != nil {
		return nil, err
	}
	log.Println("connected to ws-client")
	go jsonrpc.ServeConn(b.s)

	b.c, err = websocket.Dial("ws://localhost:" + port + "/ws-server") // Blocks until connection is established
	if err != nil {
		return nil, err
	}
	log.Println("connected to ws-server")
	b.Client = jsonrpc.NewClient(b.c)
	return b, nil
}

//Call Calls calls RPC.
func (b *Browser) Call(m string, args interface{}, reply interface{}) error {
	return b.Client.Call(m, args, reply)
}

//Close closes RPC client and server connections.
func (b *Browser) Close() {
	if err := b.c.Close(); err != nil {
		log.Println(err)
	}
	if err := b.s.Close(); err != nil {
		log.Println(err)
	}
}
