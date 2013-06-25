package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"code.google.com/p/goconf/conf" // new bsd
)

func init() {
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusForbidden)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "Only / allowed as path", http.StatusForbidden)
		return
	}
	hostname := r.Header.Get(HTTP_HOST_HEADER)
	if hostname == "" {
		http.Error(w, "No hostname header given", http.StatusForbidden)
		return
	}

	tmpfile, err := ioutil.TempFile(SPOOL_DIR, PREFIX)
	if err != nil {
		log.Printf("Error writing spool file: ", err)
		http.Error(w, "Error writing spool file", http.StatusInternalServerError)
		return
	}
	tmpfile_name := tmpfile.Name()
	defer os.Remove(tmpfile_name)

	// sanitize target path
	target_filename_unsafe := filepath.Clean(PREFIX+hostname)
	escape_check, target_filename := filepath.Split(target_filename_unsafe)
	if escape_check != "" {
		log.Printf("Error: invalid spool filename, escapes spool directory - %q", target_filename_unsafe)
		http.Error(w, "Error invalid spool filename", http.StatusInternalServerError)
	}
	target_name := filepath.Join(SPOOL_DIR, target_filename)

	body_len, err := io.Copy(tmpfile, r.Body)
	if DEBUG {
		log.Printf("Read %d bytes for %q -> %s\n", body_len, hostname, target_name)
	}
	tmpfile.Close()

	if err := os.Chmod(tmpfile_name, 0660); err != nil {
		log.Printf("Error changing permission: ", err)
	}
	if err := os.Rename(tmpfile_name, target_name); err != nil {
		log.Printf("Error renaming file: ", err)
	}
}

var (
	HTTP_HOST_HEADER string
	LISTEN           string
	SPOOL_DIR        string
	PREFIX           string
	PREFIX_TMP       string
	DEBUG            bool
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config_file := flag.String("conf", "default.conf", "Config file to use")
	flag.BoolVar(&DEBUG, "debug", false, "Enable debug output")
	flag.Parse()

	c, err := conf.ReadConfigFile(*config_file)
	if err != nil {
		log.Fatal("Error parsing config file: ", err)
	}

	LISTEN = getString(c, "", "listen")
	HTTP_HOST_HEADER = getString(c, "", "header")
	SPOOL_DIR = getString(c, "", "spool_dir")
	PREFIX = getString(c, "", "file_prefix")
	PREFIX_TMP = getString(c, "", "tmpfile_prefix")

	if !isDir(SPOOL_DIR) {
		log.Fatalf("Spool directory %s does not exist or is not a directory", SPOOL_DIR)
	}

	// routing configuration
	http.HandleFunc("/", Handler)

	log.Print("Start listening on ", LISTEN, " spool=", SPOOL_DIR)
	log.Fatal(http.ListenAndServe(LISTEN, nil))
}

func isDir(path string) bool {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	filemode := fileinfo.Mode()
	return filemode.IsDir()
}
