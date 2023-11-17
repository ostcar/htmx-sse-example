package main

import (
	"embed"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"nhooyr.io/websocket"
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
	http.HandleFunc("/websocket", handleWebsocket)
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
			<script src="/static/websocket.js"></script>
		</head>
		<body>
			<div hx-get="/poll" hx-trigger="every 1s"></div>
			<div hx-ext="sse" sse-connect="/sse" sse-swap="message"></div>
			<div hx-ext="ws" ws-connect="/websocket">
				<div id="websocket"></div>
			</div>
		</body>
	</html>
`

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(layout))
}

func handlePoll(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Poll Works: %d.", rand.Intn(1000))
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/event-stream")

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			fmt.Fprintf(w, "data: SEE Works: %d.\n\n", rand.Intn(1000))
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("Accept websocket: %v", err)
		return
	}
	defer c.CloseNow()

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if err := c.Write(r.Context(), websocket.MessageText, []byte(fmt.Sprintf(`<div id="websocket">Websocket works: %d.</div>`, rand.Intn(1000)))); err != nil {
				log.Printf("Send websocket message: %v", err)
				return
			}

		case <-r.Context().Done():
			return
		}
	}
}
