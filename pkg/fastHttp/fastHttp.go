package fastHttp


import (
	"net/http"
	"log"
	"fmt"
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

	DEFAULT_PERMISSION os.FileMode = 0644
)


// Based on this enum, we can now how work with the given file
type FastResponse struct {
	Res *http.Response
	Remote string
	Header *http.Header
	SupportConcurrent bool
	Filename string
	contentLength int64
}

// Makes a Get request to a given url and prepare the FastResponse struct to future downloads
func GetResponse(url string) (FastResponse) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(fmt.Sprintf("Error on get url %s: %q\n", url, err))
	}


	acceptRanges, hasRange := res.Header[ACCEPT_RANGE]
	contentLength := res.ContentLength
	contentDisposition, _ := res.Header[CONTENT_DISPOSITION]

	filename := utils.ExtractFilename(contentDisposition, url)
	hasRange, _ = utils.Any(acceptRanges, "bytes")

	fastResponse := FastResponse{
		res,
		url,
		&res.Header,
		hasRange && contentLength != 0,
		filename,
		contentLength}
	
	return fastResponse
}

// If remote does not accept ranges as bytes, download as a normal file
func NormalDownload(fResponse* FastResponse) {
	defer fResponse.Res.Body.Close()

	out, _ := os.OpenFile(fResponse.Filename, os.O_CREATE | os.O_WRONLY, DEFAULT_PERMISSION)
	defer out.Close()

	bar := progressbar.DefaultBytes(
		int64(fResponse.contentLength),
		"Downloading",
	)

	io.Copy(io.MultiWriter(out, bar), fResponse.Res.Body)

}