package fastHttp


import (
	"net/http"
	"log"
	"fmt"
	"strconv"
	"io"
	"os"
	"github.com/aandersonl/gofast/pkg/utils"
	"github.com/schollz/progressbar/v3"
)

// Define the type of the served file in response header
type ContentType []string


const (
	ACCEPT_RANGE string = "Accept-Ranges"
	CONTENT_TYPE string = "Content-Type"
	CONTENT_LENGTH string = "Content-Length"
	CONTENT_DISPOSITION string = "Content-Disposition"
)


// Based on this enum, we can now how work with the given file
type FastResponse struct {
	Remote string
	Header *http.Header
	SupportConcurrent bool
	Filename string
	contentLength uint64
}

// Makes a Get request to a given url and prepare the FastResponse struct to future downloads
func GetResponse(url string) (FastResponse) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(fmt.Sprintf("Error on get url %s: %q\n", url, err))
	}
	
	acceptRanges, hasRange := res.Header[ACCEPT_RANGE]
	contentLength, hasLength := res.Header[CONTENT_LENGTH]
	// TODO: Check above for filename
	//contentDispoistion, hasDispositorion := res.Header[CONTENT_DISPOSITION]

	if hasLength && len(contentLength) > 0{
		contentLengthParsed, err := strconv.ParseUint(contentLength[0], 10, 32)
			
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to extract Content-Length from header: %q\n", err)
			contentLengthParsed = 0
		}
		
		fastResponse := FastResponse{
			url,
			&res.Header,
			hasRange && utils.Any(acceptRanges, "bytes") && contentLengthParsed != 0,
			"nofilename",
			contentLengthParsed}

		return fastResponse
	} else {
		panic("Unable to download, content-length not found")
	}
}

func NormalDownload(fResponse* FastResponse) {
	res, err := http.Get(fResponse.Remote)

	if err != nil {
		panic(fmt.Sprintf("Unable to download: %s\n", fResponse.Remote))
	}

	defer res.Body.Close()

	out, _ := os.OpenFile("test", os.O_CREATE | os.O_WRONLY, 0644)
	defer out.Close()

	bar := progressbar.DefaultBytes(
		int64(fResponse.contentLength),
		"Downloading",
	)

	io.Copy(io.MultiWriter(out, bar), res.Body)

}