package main

import (
	"fmt"
	"os"

	"github.com/aandersonl/gofast/pkg/fastHttp"
)

func main() {

	var url = os.Args[1]

	var response fastHttp.FastResponse = fastHttp.GetResponse(url)

	if response.SupportConcurrent {
		fmt.Printf("Starting concurrent download of %s\n", response.Filename)
		fastHttp.ConcurrentDownload(&response)
	} else {
		fmt.Printf("Remote doens't support multiple connections, downloading %s as normal file\n", response.Filename)
		fastHttp.NormalDownload(&response)
	}
}
