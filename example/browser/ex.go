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

package main

import (
	"log"

	"github.com/gopherjs/jquery"

	"github.com/utamaro/neje-ui/browser"
)

//Jquery is for convenience.
var jQuery = jquery.NewJQuery

//GUI is struct to bel called from remote by rpc.
type GUI struct{}

//Write writes a response from the server.
func (g *GUI) Write(msg *string, reply *string) error {
	//show welcome message:
	jQuery("#from_server").SetText(msg)
	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	b, err := browser.New(new(GUI))
	if err != nil {
		log.Fatal(err)
	}
	//defer b.Close()
	jQuery("button").On(jquery.CLICK, func(e jquery.Event) {
		go func() {
			m := jQuery("#to_server").Val()
			response := ""
			err = b.Call("Msg.Message", &m, &response)
			if err != nil {
				log.Fatal(err)
			}
			//show welcome message:
			jQuery("#response").SetText(response)
		}()
	})

}
