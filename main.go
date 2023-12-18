package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/csv"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/didip/tollbooth/v7"
	"errors"
)

func main() {
	cf := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *cf == "" {
		log.Fatal("cannot run without configuration file")
	}

	cfg, err := LoadConfig(*cf)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(serve(&cfg))
}

func serve(cfg *Config) error {
	// CSV Setup
	f, err := os.OpenFile(cfg.OutputFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	csvw := csv.NewWriter(f)

	if fi.Size() < 1 {
		if err := csvw.Write(CSVHeader); err != nil {
			return err
		}

		csvw.Flush()
	} else {
		log.Println("Warning: will not write CSV header to existing output file")
	}

	// Limiter setup
	limiter := tollbooth.NewLimiter(cfg.RateLimit, nil)
	limiter.SetMethods([]string{"POST"})
	if cfg.CFMode {
		limiter.SetIPLookups([]string{"CF-Connecting-IP", "X-Forwarded-For", "RemoteAddr", "X-Real-IP"})
	}

	tlsConfig := &tls.Config{}

	/* mTLS optional setup
	   Note: This is NOT the cert/key for the server.t
	   In the case of cloudflare, it will come from:
	   https://developers.cloudflare.com/ssl/static/authenticated_origin_pull_ca.pem
	*/
	if cfg.MTLSFile != "" {
		mtlsCert, err := ioutil.ReadFile(cfg.MTLSFile)
		if err != nil {
			return err
		}
		certPool := x509.NewCertPool()
		if success := certPool.AppendCertsFromPEM(mtlsCert); !success {
			return errors.New("Failed to add cert to pool")
		}

		tlsConfig = &tls.Config{
			ClientCAs:  certPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
	}

	httpServer := &http.Server{
		Addr:      ":" + cfg.Port,
		TLSConfig: tlsConfig,
	}

	// Handler setup
	http.Handle("/", tollbooth.LimitFuncHandler(limiter, func(w http.ResponseWriter, req *http.Request) {
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
			log.Printf("Failed to decode data: %s", err)
			return
		}

		if err := data.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Client gave bad data")
			return
		}

		data.Sanitize()

		log.Println("Recieved data, writing to CSV output file")

		defer csvw.Flush()
		if err := csvw.Write(data.CSV()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Unable to write to csv: %s", err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}))

	log.Println("Serving")

	// Start HTTP
	return httpServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
}
