// +build ignore

package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/MJKWoolnough/scs"
)

var (
	port  = flag.Int("port", 11212, "The port the server listens on")
	limit = flag.Int("items", 65535, "Total number of items the server can store")
)

func myLogFunc(err error) {
	fmt.Println(err)
}

func main() {
	flag.Parse()
	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: *port})
	if err != nil {
		fmt.Println(err)
		return
	}
	s := scs.NewStore(scs.Limit(uint64(*limit)), scs.LogFunc(myLogFunc))
	s.Serve(l)
}
