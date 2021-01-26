package main

import (
	"fmt"
	"os"
	"github.com/akamensky/argparse"	
	"github.com/aandersonl/gofast/pkg/fastHttp"
)



func main() {
	parser := argparse.NewParser("gofast", "A HTTP downloader accelerator using concurrency")
	var numTasks *int = parser.Int("n", "num-tasks", &argparse.Options{Required: false,
									  Help: "Number of concurrent connections, default: Num cores * 2",
									  Default: 0})

	var outputFilename *string = parser.String("o", "output", &argparse.Options{Required: false,
									           Help: "Output filename",
										   Default: ""})
	var url *string = parser.String("u", "url", &argparse.Options{Required: true,
								     Help: "Remote file url"})
	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}


	var response fastHttp.FastResponse = fastHttp.GetResponse(*url)

	if (*outputFilename != "") {
		response.Filename = *outputFilename
	}

	if response.SupportConcurrent {
		fmt.Printf("Starting concurrent download of %s\n", response.Filename)
		fastHttp.ConcurrentDownload(&response, *numTasks)
	} else {
		fmt.Printf("Remote doens't support multiple connections, downloading %s as normal file\n", response.Filename)
		fastHttp.NormalDownload(&response)
	}
}
