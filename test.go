// +build ignore

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"sync"
)

var port = flag.Int("port", 11212, "The port the server listens on")

const (
	input = "set sushi\r\n" +
		"delicious\r\n" +
		"set topcoder\r\n" +
		"fun\r\n" +
		"get sushi topcoder\r\n" +
		"get sushi\r\n" +
		"delete sushi\r\n" +
		"get sushi\r\n" +
		"get topcoder\r\n" +
		"get topcoder sushi\r\n" +
		"delete sushi\r\n" +
		"stats\r\n" +
		"quit\r\n"
	output = "STORED\r\n" +
		"STORED\r\n" +
		"VALUE sushi\r\n" +
		"delicious\r\n" +
		"VALUE topcoder\r\n" +
		"fun\r\n" +
		"END\r\n" +
		"VALUE sushi\r\n" +
		"delicious\r\n" +
		"END\r\n" +
		"DELETED\r\n" +
		"END\r\n" +
		"VALUE topcoder\r\n" +
		"fun\r\n" +
		"END\r\n" +
		"VALUE topcoder\r\n" +
		"fun\r\n" +
		"END\r\n" +
		"NOT_FOUND\r\n" +
		"cmd_get 7\r\n" +
		"cmd_set 2\r\n" +
		"get_hits 5\r\n" +
		"get_misses 2\r\n" +
		"delete_hits 1\r\n" +
		"delete_misses 1\r\n" +
		"curr_items 1\r\n" +
		"limit_items 65535\r\n" +
		"END\r\n"
)

func main() {
	flag.Parse()
	c, err := net.Dial("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		fmt.Println(err)
		return
	}

	var w sync.WaitGroup
	w.Add(2)
	go func() {
		_, err = c.Write([]byte(input))
		if err != nil {
			fmt.Println(err)
		}
		w.Done()
	}()
	var result string
	go func() {
		d, err := ioutil.ReadAll(c)
		if err != nil {
			fmt.Println(err)
		}
		result = string(d)
		w.Done()
	}()
	w.Wait()
	if result == output {
		fmt.Println("\x1b[32mPASSED\x1b[39m")
	} else {
		fmt.Printf("\x1b[31mFAILED\x1b[39m\nexpecting: %s\n\ngot: %s\n", output, result)
	}
}
