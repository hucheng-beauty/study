package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// to current-limiting
var rateLimiter = time.Tick(10 * time.Millisecond)

// Fetch is according to the url get html page.
func Fetch(url string) ([]byte, error) {
	<-rateLimiter
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Deal with status which the StatusCode is not 200.
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: status code is %d.", resp.StatusCode)
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}

	// Style others to UTF-8.
	bodyReader := bufio.NewReader(resp.Body)
	e := determinerEncoding(bodyReader)
	utfReader := transform.NewReader(resp.Body, e.NewDecoder())

	return ioutil.ReadAll(utfReader)
}

func determinerEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
