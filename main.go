package main // import "github.com/dd-caleb/dd-trace-mock-agent"

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/DataDog/datadog-trace-agent/model"
	"github.com/tinylib/msgp/msgp"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer w.WriteHeader(http.StatusOK)

		log.Println(r.Method, r.RequestURI)
		traces, _ := getTraces(w, r)

		for _, trace := range traces {
			bs, _ := json.MarshalIndent(trace.APITrace(), "", "  ")
			log.Println(string(bs))
		}
	})

	log.Println("listening for traces on 127.0.0.1:8126")
	err := http.ListenAndServe("127.0.0.1:8126", nil)
	if err != nil {
		panic(err)
	}
}

func getTraces(w http.ResponseWriter, r *http.Request) (model.Traces, bool) {
	var traces model.Traces
	contentType := r.Header.Get("Content-Type")

	if err := decodeReceiverPayload(r.Body, &traces, contentType); err != nil {
		log.Printf("cannot decode traces payload: %v", err)
		return nil, false
	}

	return traces, true
}

func decodeReceiverPayload(r io.Reader, dest msgp.Decodable, contentType string) error {
	switch contentType {
	case "application/msgpack":
		return msgp.Decode(r, dest)

	case "application/json":
		fallthrough
	case "text/json":
		fallthrough
	case "":
		return json.NewDecoder(r).Decode(dest)

	default:
		panic(fmt.Sprintf("unhandled content type %q", contentType))
	}
}
