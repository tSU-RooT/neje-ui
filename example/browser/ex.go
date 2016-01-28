package main

import (
	"log"
	"strconv"

	"github.com/gopherjs/jquery"
	"github.com/utamaro/wsrpc/browser"
)

//convenience:
var jQuery = jquery.NewJQuery

//aa
const (
	INPUT   = "input#name"
	OUTPUT  = "span#output"
	OUTPUT2 = "span#output2"
)

//Args is
type Args struct {
	A int
	B int
	C string
}

//GUI is
type GUI struct{}

//Write is
func (g *GUI) Write(args *Args, reply *int) error {
	//show welcome message:
	jQuery(OUTPUT).SetText("from server:" + args.C)
	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	b, err := browser.New("localhost:7000", new(GUI))
	if err != nil {
		log.Fatal(err)
	}
	//	defer b.Close()
	jQuery(INPUT).On(jquery.KEYUP, func(e jquery.Event) {
		go func() {
			args := Args{A: 17, B: 8}
			var reply int
			err = b.Call("Arith.Multiply", args, &reply)
			if err != nil {
				log.Fatal("arith error:", err)
			}
			//show welcome message:
			jQuery(OUTPUT2).SetText("result of mult" + strconv.Itoa(reply))
		}()
	})

}
