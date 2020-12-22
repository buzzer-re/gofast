package main

import (
	"fmt"
	"os"
	"github.com/aandersonl/gofast/pkg/fastHttp"
)



func main() {
	
	var url = os.Args[1]

	var res fastHttp.FastResponse = fastHttp.GetResponse(url)

	fmt.Printf("Starting downloading %s\n", url)
	fmt.Println(res)

}