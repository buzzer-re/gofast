package fastHttp

import (
	"fmt"
	"io"
	"log"
	"time"
	"strings"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"sync"
	"github.com/aandersonl/gofast/pkg/utils"
	"github.com/schollz/progressbar/v3"
)

// Define the type of the served file in response header
type ContentType []string

const (
	ACCEPT_RANGE        string = "Accept-Ranges"
	CONTENT_TYPE        string = "Content-Type"
	CONTENT_LENGTH      string = "Content-Length"
	CONTENT_DISPOSITION string = "Content-Disposition"

	DEFAULT_PERMISSION os.FileMode = 0644
)

//Task struct that hold where the concurrent download should seek in the file and where must end
type Task struct {
	Offset int64
	End int64
	Size int64
}

//FastResponse Based on this enum, we can now how work with the given file
type FastResponse struct {
	Res               *http.Response
	Remote            string
	Header            *http.Header
	SupportConcurrent bool
	Filename          string
	contentLength     int64
}

//GetResponse makes a Get request to a given url and prepare the FastResponse struct to future downloads
func GetResponse(url string) FastResponse {
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

//NormalDownload If remote does not accept ranges as bytes, download as a normal file
func NormalDownload(fResponse *FastResponse) {
	defer fResponse.Res.Body.Close()

	out, _ := os.OpenFile(fResponse.Filename, os.O_CREATE|os.O_WRONLY, DEFAULT_PERMISSION)
	defer out.Close()

	bar := progressbar.DefaultBytes(
		int64(fResponse.contentLength),
		"Downloading",
	)

	io.Copy(io.MultiWriter(out, bar), fResponse.Res.Body)
}

func ConcurrentDownload(fResponse *FastResponse,numTasks int) {
	out, _ := os.OpenFile(fResponse.Filename, os.O_CREATE|os.O_WRONLY, DEFAULT_PERMISSION)

	syscall.Fallocate(int(out.Fd()), 0, 0, fResponse.contentLength)
	out.Close()

	var wg sync.WaitGroup
	var taskNum int
	if numTasks != 0 {
		taskNum = numTasks
	} else {
		taskNum = runtime.NumCPU()
		taskNum = taskNum * 2
	}

	tasks     := make([]Task, taskNum)
	splitSize := int64(fResponse.contentLength/int64(taskNum))
//	remain	  := fResponse.contentLength % int64(taskNum)

	start := time.Now()
	bar := progressbar.DefaultBytes(
		fResponse.contentLength,
		"Downloading")

	for i,task := range(tasks) {
		index := int64(i)
		task = Task{index*splitSize, (index*splitSize) + splitSize, splitSize}
		wg.Add(1)
		go downloadRange(fResponse, task, &wg, bar)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Downloaded in %s\n", elapsed)
}

//downloadRange task downloads a portion of the remote file based in the Range Requests(RFC-7233)
func downloadRange(fResponse *FastResponse, task Task,  wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	out, _ := os.OpenFile(fResponse.Filename, os.O_RDWR, DEFAULT_PERMISSION)
	defer out.Close()
	out.Seek(task.Offset, 0)

	req, err := http.NewRequest("GET", fResponse.Remote, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", task.Offset, task.End-1))
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if !strings.Contains(res.Status,"206") {
		panic("Remote doesn't responded with the expected code: 206")
	}

	defer res.Body.Close()

	io.Copy(io.MultiWriter(out, bar), res.Body)
}


