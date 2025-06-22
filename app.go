package netclip

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

var dataStore = NewDataStore()

//go:embed static
var staticFiles embed.FS

// AppVersion holds the application version
var AppVersion = "0.0.5"

// Run starts the server using the given port (string)
func Run(port string, certFile string, keyFile string) {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/save", SaveHandler)
	http.HandleFunc("/delete", DeleteHandler)
	http.HandleFunc("/static/", StaticFileHandler)

	if certFile == "" && keyFile == "" {
		log.Println("starting http")
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal("Could not start server: ", err)
		}
	} else {
		log.Println("starting https")
		err := http.ListenAndServeTLS(":"+port, certFile, keyFile, nil)
		if err != nil {
			log.Fatal("Could not start server: ", err)
		}
	}
}

// IndexHandler shows the page that displays the form and the results
func IndexHandler(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	templateData := struct {
		AppVersion string
		DataStore  *DataStore
		Year       int
	}{
		AppVersion: AppVersion,
		DataStore:  &dataStore,
		Year:       time.Now().Year(),
	}

	indexTemplate, err := template.ParseFS(staticFiles, "static/index.html")

	if err != nil {
		log.Printf("Error parsing index template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = indexTemplate.Execute(w, templateData)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// StaticFileHandler serves static files from the embedded file system
func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	data, err := staticFiles.ReadFile(filePath)
	if err != nil {
		// Handle the error, e.g., send a 404 status code
		//w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, err)

		return
	}

	// Set the appropriate Content-Type header
	if strings.HasSuffix(filePath, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(filePath, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	} else {
		// You can add more file types here if needed
		w.Header().Set("Content-Type", "text/plain")
	}

	_, _ = w.Write(data)
}

// SaveHandler saves records to the DataStore
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := r.ParseForm()
	if err != nil {
		// in case of any error
		_, _ = fmt.Fprint(w, "<h1>Error processing form</h1>")
		return
	}

	textToSave := r.PostForm.Get("text")

	if textToSave == "" {
		// in case of any error
		_, _ = fmt.Fprint(w, "<h1>Text is blank</h1>")
		return
		// Use a unique key for each saved text
	}

	key := fmt.Sprintf("%d", time.Now().UnixNano())
	dataStore.Store(key, textToSave)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteHandler deletes records from the DataStore
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := r.ParseForm()
	if err != nil {
		_, _ = fmt.Fprint(w, "<h1>Error processing form</h1>")
		return
	}

	keyToDelete := r.PostForm.Get("key")

	if keyToDelete == "" {
		_, _ = fmt.Fprint(w, "<h1>Key is blank</h1>")
		return
	}

	dataStore.Delete(keyToDelete)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
