package real

import (
	"net/http"
	"net/http/httputil"
	"time"
)

type Retriever struct {
	UserAgent string
	TimeOut   time.Duration
}

func (r *Retriever) Get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	return string(out)

	// out, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	//	panic(err)
	// }
	// return string(out)
}
