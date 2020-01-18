package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"unicode"
)

// SafeWriter only writes printable charactars
type SafeWriter struct {
	w io.Writer
}

// Write removes unprintable characters from p and calls Write of underlaying Writer
func (s *SafeWriter) Write(p []byte) (n int, err error) {
	var buf = []byte{}
	for _, c := range p {
		if unicode.IsPrint(rune(c)) {
			buf = append(buf, c)
		}
	}
	return s.w.Write(buf)
}

func main() {
	var (
		forceJSON   bool
		enableTLS   bool
		showHeaders bool
		bindAddr    string
	)

	flag.BoolVar(&forceJSON, "json", false, "treats all bodies as json")
	flag.BoolVar(&enableTLS, "tls", false, "serve tls with tls.crt and tls.key")
	flag.BoolVar(&showHeaders, "header", true, "show headers")
	flag.StringVar(&bindAddr, "bind", ":8080", "bind address")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if showHeaders {
			header, err := httputil.DumpRequest(r, false)
			if err != nil {
				fmt.Fprintln(w, err)
				log.Println("could not parse header")
				return
			}
			fmt.Print(string(header))
		}

		payload := make(map[string]interface{})

		contentType := r.Header.Get("Content-Type")
		if contentType == "application/json" || forceJSON {
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				log.Println(err)
				fmt.Fprintln(w, err)
				io.Copy(&SafeWriter{os.Stdout}, r.Body)
				fmt.Println("")
				fmt.Println("")
				return
			} else {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				err = enc.Encode(payload)
				if err != nil {
					fmt.Fprintln(w, err)
					log.Println(err)
					return
				}
			}
		} else {
			io.Copy(&SafeWriter{os.Stdout}, r.Body)
			fmt.Println("")
			fmt.Println("")
		}
	})

	if enableTLS {
		log.Fatal(http.ListenAndServeTLS(bindAddr, "tls.crt", "tls.key", nil))
	} else {
		log.Fatal(http.ListenAndServe(bindAddr, nil))
	}
}
