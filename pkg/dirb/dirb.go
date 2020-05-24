package dirb

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/leobeosab/sharingan/internal/helpers"
	"github.com/schollz/progressbar/v2"
)

// Response data container for response information
type Response struct {
	Path          string
	StatusCode    int
	ContentLength int64
}

// Dirb execute direcotry busting on specified domain
func Dirb(d string, wlp string, t int) []Response {
	wl, err := os.Open(wlp)
	if err != nil {
		panic("Specifed wordlsit cannot be found")
	}
	defer wl.Close()

	// TODO: implement proxy stuff
	// TODO: recursive stuff

	lines := helpers.GetNumberOfLinesInFile(wl)
	progress := progressbar.NewOptions(lines)

	// domain + path
	jobs := make(chan string, lines)
	responses := make(chan Response, lines)

	var wg sync.WaitGroup

	for ji := 0; ji < t; ji++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for job := range jobs {
				resp, _ := http.Get(job)
				progress.Add(1)

				if resp.StatusCode == 404 {
					continue
				}

				responses <- Response{
					Path:          job,
					StatusCode:    resp.StatusCode,
					ContentLength: resp.ContentLength,
				}
			}
		}()
	}

	// Add trailing slash if it doesn't exist
	if d[len(d)-1:] != "/" {
		d = d + "/"
	}

	// Add possible subdomains to jobs list
	stream := bufio.NewScanner(wl)
	for stream.Scan() {
		jobs <- d + stream.Text()
	}
	if err := stream.Err(); err != nil {
		log.Fatal(err)
	}

	close(jobs)
	wg.Wait()
	close(responses)
	progress.Finish()

	result := make([]Response, 0)
	for r := range responses {
		fmt.Printf("%s - %d - %d\n", r.Path, r.StatusCode, r.ContentLength) // REMOVE
		result = append(result, r)
	}

	return result
}
