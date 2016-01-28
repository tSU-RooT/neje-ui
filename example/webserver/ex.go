package main

import (
	"log"
	"strconv"
	"time"

	"github.com/utamaro/wsrpc/webserver"
)

//Args is
type Args struct {
	A int
	B int
	C string
}

//Arith is
type Arith struct{}

//Multiply is
func (t *Arith) Multiply(args *Args, reply *int) error {
	log.Println("mult")
	*reply = args.A * args.B
	return nil
}
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	ws, err := webserver.New("localhost:7000", new(Arith))
	if err != nil {
		log.Fatal(err)
	}
	var reply int
	for i := 0; i < 10; i++ {
		log.Println("writing", i)
		if err := ws.Call("GUI.Write", &Args{C: "test" + strconv.Itoa(i)}, &reply); err != nil {
			log.Fatal(err)
		}
		time.Sleep(10 * time.Second)
	}
}
