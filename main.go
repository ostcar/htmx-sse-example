package main

import (
	"embed"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//go:embed static
var static embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	http.Handle("/static/", http.FileServer(http.FS(static)))
	http.HandleFunc("/sse", handleSSE)
	http.HandleFunc("/poll", handlePoll)
	http.HandleFunc("/", handleHome)
	log.Println("listen on :8080")
	return http.ListenAndServe(":8080", nil)
}

const layout = `
<!DOCTYPE html>
	<html lang="de">
		<head>
			<meta charset="utf-8"/>
			<title>test</title>

			<script src="/static/htmx.js"></script>
			<script src="/static/sse.js"></script>

			<!--
			<script src="https://unpkg.com/htmx.org@1.9.6"></script>
			<script src="https://unpkg.com/htmx.org/dist/ext/sse.js"></script>
			-->
		</head>
		<body>
			<div hx-ext="sse" sse-connect="/sse" sse-swap="message"></div>
			<div hx-get="/poll" hx-trigger="every 1s"></div>
		</body>
	</html>
`

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(layout))
}

func handlePoll(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "poll %d. <strong>this is a </strong> <a href=\"https://github.com/bigskysoftware/htmx/pull/1794\">Here</a> is a fantastic pr.", rand.Intn(1000))
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/event-stream")

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			fmt.Fprintf(w, "data: sse %d. <a href=\"https://github.com/bigskysoftware/htmx/pull/1794\">Here</a> is a fantastic pr.\n\n", rand.Intn(1000))
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}
