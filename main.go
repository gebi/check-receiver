package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only posts allowed", http.StatusForbidden)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "Only / allowed as path", http.StatusForbidden)
		return
	}
	hostname := r.Header.Get(*HTTP_HOST_HEADER)
	if hostname == "" {
		http.Error(w, "No hostname header given", http.StatusForbidden)
		return
	}

	tmpfile, err := ioutil.TempFile(*SPOOL_DIR, *PREFIX)
	defer os.Remove(tmpfile.Name())
	if err != nil {
		log.Printf("Error writing spool file: ", err)
		http.Error(w, "Error writing spool file", http.StatusInternalServerError)
		return
	}

	body_len, err := io.Copy(tmpfile, r.Body)
	log.Printf("Read %d bytes for %s\n", body_len, hostname)
	tmpfile.Close()

	target_name := filepath.Join(*SPOOL_DIR, *PREFIX+hostname)
	os.Rename(tmpfile.Name(), target_name)
}

var (
	HTTP_HOST_HEADER *string
	LISTEN           *string
	SPOOL_DIR        *string
	PREFIX           *string
	PREFIX_TMP       *string
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	HTTP_HOST_HEADER = flag.String("header", "REMOTE_USER", "Header to get hostname from")
	LISTEN = flag.String("listen", "localhost:8443", "Socket to listen on")
	SPOOL_DIR = flag.String("spool", "", "Spool directory for uploaded files (required!)")
	PREFIX = flag.String("prefix", "nagios-receiver.", "Prefix name for spool files")
	tmp_prefix := (*PREFIX)[:len(*PREFIX)-1] + "-tmp."
	PREFIX_TMP = &tmp_prefix
	flag.Parse()

	if *SPOOL_DIR == "" {
		log.Fatal("No --spool <directory> configured")
	}

	// routing configuration
	http.HandleFunc("/", Handler)

	log.Print("Start listening on ", *LISTEN, " spool=", *SPOOL_DIR)
	log.Fatal(http.ListenAndServe(*LISTEN, nil))
}
