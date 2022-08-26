package main

import (
	json2 "encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"url/url"
)

var (
	logEnabled *bool
	port       *int
	urlBase    string
)

func init() {
	domain := flag.String("d", "localhost", "domain")
	port = flag.Int("p", 8888, "port")
	logEnabled = flag.Bool("l", true, "log enabled/disabled")
	urlBase = fmt.Sprintf("http://localhost:%d", port)
	flag.Parse()
	urlBase = fmt.Sprintf("http://%s:%d", *domain, *port)
}

type Headers map[string]string
type Redirector struct {
	stats chan string
}

func (r *Redirector) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	findUrlAndExecute(writer, request, func(url *url.Url) {
		http.Redirect(writer, request, url.Final, http.StatusMovedPermanently)
		r.stats <- url.Id
	})
}

func Tiny(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		responseWith(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})
		return
	}

	url, newUrl, err := url.FindOrCreateNewUrl(extractUrl(r))

	if err != nil {
		responseWith(w, http.StatusBadRequest, nil)
		return
	}

	var status int
	if newUrl {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	tinyUrl := fmt.Sprintf("%s/r/%s", urlBase, url.Id)

	responseWith(w, status, Headers{
		"Location": tinyUrl,
		"Link":     fmt.Sprintf("<%s/api/stats/%s>; rel=\"stats\"", urlBase, url.Id),
	})

	logg("URL %s shorted with success to %s", url.Final, tinyUrl)
}

func Visualizer(w http.ResponseWriter, r *http.Request) {
	findUrlAndExecute(w, r, func(url *url.Url) {
		json, err := json2.Marshal(url.Stats)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseWithJSON(w, string(json))
	})

}

func findUrlAndExecute(w http.ResponseWriter, r *http.Request, executor func(url2 *url.Url)) {
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]
	if url := url.Find(id); url != nil {
		executor(url)
	} else {
		http.NotFound(w, r)
	}
}

func responseWith(w http.ResponseWriter, status int, headers Headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
}
func responseWithJSON(w http.ResponseWriter, response string) {
	responseWith(w, http.StatusOK, Headers{"Content-Type": "application/json"})
	fmt.Fprintf(w, response)
}
func extractUrl(r *http.Request) string {
	url := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(url)
	return string(url)
}

func buildStatistics(stats <-chan string) {
	for id := range stats {
		url.RegisterClick(id)
		logg("Click registered with success to %s.", id)
	}
}

func logg(format string, values ...interface{}) {
	if *logEnabled {
		log.Printf(fmt.Sprintf("%s\n", format), values...)
	}
}

func main() {
	url.ConfigRepository(url.NewMemoryRepository())

	stats := make(chan string)
	defer close(stats)
	go buildStatistics(stats)

	http.Handle("/r/", &Redirector{stats})
	http.HandleFunc("/api/tiny", Tiny)
	http.HandleFunc("/api/stats/", Visualizer)

	logg("Starting server on port %d...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))

}
