package netclip

import (
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"text/template"
	"time"

	"tailscale.com/tsnet"
)

var dataStore = NewDataStore()

//go:embed static
var staticFiles embed.FS

// AppVersion holds the application version
var AppVersion = "0.6.1"

// Server interface for different server types
type Server interface {
	Listen() (net.Listener, error)
	Serve(ln net.Listener) error
}

// CreateServer creates the appropriate server type based on configuration
func CreateServer(config Config, authKey string) Server {
	if config.Tailscale.Enabled {
		return &TSNetServer{
			Hostname: config.Tailscale.Hostname,
			AuthKey:  authKey,
			UseTLS:   config.Tailscale.UseTLS,
		}
	}
	return &HTTPServer{
		Port:     config.Port,
		CertFile: config.CertFile,
		KeyFile:  config.KeyFile,
	}
}

// setupHandlers registers all HTTP handlers
func setupHandlers() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/save", SaveHandler)
	http.HandleFunc("/delete", DeleteHandler)
	http.HandleFunc("/static/", StaticFileHandler)
}

// HTTPServer implements Server interface for regular HTTP/HTTPS
type HTTPServer struct {
	Port     string
	CertFile string
	KeyFile  string
}

func (s *HTTPServer) Listen() (net.Listener, error) {
	return net.Listen("tcp", ":"+s.Port)
}

func (s *HTTPServer) Serve(ln net.Listener) error {
	if s.CertFile == "" && s.KeyFile == "" {
		log.Println("starting http on port", s.Port)
		return http.Serve(ln, nil)
	} else {
		log.Println("starting https on port", s.Port)
		tlsConfig := &tls.Config{}
		tlsListener := tls.NewListener(ln, tlsConfig)
		return http.ServeTLS(tlsListener, nil, s.CertFile, s.KeyFile)
	}
}

// TSNetServer implements Server interface for Tailscale networking
type TSNetServer struct {
	Hostname string
	AuthKey  string
	UseTLS   bool
}

func (s *TSNetServer) Listen() (net.Listener, error) {
	srv := &tsnet.Server{
		Hostname: s.Hostname,
	}

	// Only set AuthKey if provided, otherwise TSNet will prompt for manual auth
	if s.AuthKey != "" {
		srv.AuthKey = s.AuthKey
	}

	addr := ":80"
	if s.UseTLS {
		addr = ":443"
	}

	ln, err := srv.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	if s.UseTLS {
		lc, err := srv.LocalClient()
		if err != nil {
			return nil, err
		}
		ln = tls.NewListener(ln, &tls.Config{
			GetCertificate: lc.GetCertificate,
		})
	}

	return ln, nil
}

func (s *TSNetServer) Serve(ln net.Listener) error {
	if s.UseTLS {
		log.Printf("starting TSNet HTTPS server as %s", s.Hostname)
	} else {
		log.Printf("starting TSNet HTTP server as %s", s.Hostname)
	}
	return http.Serve(ln, nil)
}

// Run starts the server using the provided Server implementation
func Run(server Server) {
	setupHandlers()

	ln, err := server.Listen()
	if err != nil {
		log.Fatal("Could not create listener: ", err)
	}

	err = server.Serve(ln)
	if err != nil {
		log.Fatal("Could not start server: ", err)
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
