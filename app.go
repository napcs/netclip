package netclip

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type DataStore struct {
	data map[string]string
	mu   sync.Mutex
}

func (ds *DataStore) Store(key, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[key] = value
}

func (ds *DataStore) Range(f func(key, value string) bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for key, value := range ds.data {
		if !f(key, value) {
			break
		}
	}
}

var dataStore = DataStore{
	data: make(map[string]string),
}

var AppVersion = "0.0.2"

func Run(port string) {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveHandler)
	_ = http.ListenAndServe(":"+port, nil)

}

// Index page shows the subscription form.
func indexHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	head := `<!DOCTYPE html5>
<html>
  <head>
    <meta charset="utf-8">
    <title>netclip</title>
    <meta name=viewport content="width=device-width,initial-scale=1">
	<style>

      .container {
        width: 80%;
        margin: 0 auto;
      }

      main form {

        display: flex;
        flex-direction:column;
      }

      main form textarea {
        height: 6rem;
      }

	</style>
  </head>
  <body>
    <div class="container">
      <main>
        <h1>netclip</h1> `

	footer := "</main><footer><small>netclip v" + AppVersion + " &copy; 2023 Brian Hogan</small></footer></body></html>"

	_, _ = fmt.Fprint(w, head)
	_, _ = fmt.Fprint(w, "<form method='post' action='save'><textarea required name='text'></textarea><br><input type='submit' value='Save'>")

	_, _ = fmt.Fprint(w, "<div class='items'>")

	// Collect keys and values from the datastore
	keys := make([]string, 0)
	dataStore.mu.Lock()
	for key := range dataStore.data {
		keys = append(keys, key)
	}
	dataStore.mu.Unlock()

	// Sort keys in reverse order
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	_, _ = fmt.Fprint(w, "<h1>All saved texts</h1>")

	// Display sorted data
	for _, key := range keys {
		value, _ := dataStore.data[key]
		_, _ = fmt.Fprintf(w, "<div class='item'><p>%s</p><pre>%s</pre></div>", key, value)
	}
	_, _ = fmt.Fprint(w, "</div>")
	_, _ = fmt.Fprint(w, footer)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

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
