package fastHttp


import (
	"net/http"
	"log"
	"fmt"
	"strings"
)

// Define the type of the served file in response header
type ContentType int

// Based on this enum, we can now how work with the given file
const (
	HTML_CONTENT ContentType = iota
	TEXT_CONTENT
	BLOB_CONTENT
	MEDIA_CONTENT
)

const (
	HTTP_OK int = 200
	HTTP_FOUND  	   = 302
	HTTP_ERROR 		   = 500
)

type FastResponse struct {
	contentType ContentType
	remote string
	header *http.Header
}

func GetResponse(url string) (FastResponse) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(fmt.Sprintf("Error on get url %s: %q\n", url, err))
	}

	
	contentTypeString := strings.Join(res.Header["Content-Type"], ";")
	
	fmt.Println(contentTypeString)
	
	return FastResponse{}
}