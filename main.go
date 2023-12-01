package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: merlin output")
		os.Exit(1)
	}

	f, err := os.OpenFile(os.Args[1], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	csvw := csv.NewWriter(f)

	if fi.Size() < 1 {
		if err := csvw.Write(CSVHeader); err != nil {
			panic(err)
		}

		csvw.Flush()
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			log.Printf("Client attempted %s", req.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if req.Header.Get("Content-type") != "application/json" {
			log.Println("Client attempted bad content type")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		dec := json.NewDecoder(req.Body)
		var data Data
		if err := dec.Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Client gave invalid decoded data: %s", err)
			return
		}

		if err := data.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Client gave bad data: %s", err)
			return
		}

		data.Sanitize()

		log.Println("Recieved data, writing to CSV output file")

		defer csvw.Flush()
		if err := csvw.Write(data.CSV()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("unable to write to csv: %s", err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})

	log.Println("Serving")

	err = http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		panic(err)
	}
}
