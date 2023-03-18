package netclip

import (
	"fmt"
	"net/http"
	"sort"
	"time"
)

var dataStore = DataStore{
	data: make(map[string]string),
}

// AppVersion holds the application version
var AppVersion = "0.0.3"

func Banner() {
	fmt.Println("netclip v" + AppVersion)
}

// Run starts the server using the given port (string)
func Run(port string) {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/delete", deleteHandler)
	_ = http.ListenAndServe(":"+port, nil)
}

// Index page shows the subscription form.
func indexHandler(w http.ResponseWriter, _ *http.Request) {

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

	w.Header().Set("Content-Type", "text/html")

	_, _ = fmt.Fprint(w, head())
	_, _ = fmt.Fprint(w, "<form method='post' action='save'><textarea required name='text'></textarea><br><input type='submit' value='Save'></form>")
	_, _ = fmt.Fprint(w, "<div class='items'>")
	_, _ = fmt.Fprint(w, "<h1>Saved clips</h1>")

	// Display sorted data
	for _, key := range keys {
		value, _ := dataStore.data[key]
		_, _ = fmt.Fprintf(w, "<div class='item'><div class='snippet'><pre>%s</pre></div><form method='post' action='/delete'><input type='hidden' value='%s' name='key'><input type='submit' value='Delete this clip'></form></div>", value, key)
	}
	_, _ = fmt.Fprint(w, "</div>") // closing outer div
	_, _ = fmt.Fprint(w, script())
	_, _ = fmt.Fprint(w, footer())
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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
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

func head() string {

	str := `<!DOCTYPE html5>
<html>
  <head>
    <meta charset="utf-8">
    <title>netclip</title>
    <meta name=viewport content="width=device-width,initial-scale=1">
	<style>`

	str += css()

	str += `
	</style>
  </head>
  <body>
    <div class="container">
      <main>
        <h1>netclip</h1> `

	return str
}

func footer() string {
	return "</main><footer><small>netclip v" + AppVersion + " &copy; 2023 Brian Hogan</small></footer></body></html>"
}

func script() string {
	return `
<script>
function addButtons() {
  var snippets = document.querySelectorAll('.snippet pre');
  var numberOfSnippets = snippets.length;


  for (var i = 0; i < numberOfSnippets; i++) {
    var p = snippets[i].parentElement;
    var b = document.createElement("button");
    b.classList.add('btn-copy')
    b.innerText="Copy";

    b.addEventListener("click", function () {
      this.innerText = 'Copying..';
      code = this.nextSibling.innerText;
      console.log(this.nextSibling);
      navigator.clipboard.writeText(code);
      this.innerText = 'Copied!';
      var that = this;
      setTimeout(function () {
        that.innerText = 'Copy';
      }, 1000)
    });
    p.prepend(b)
  }
}

addButtons();


document.querySelectorAll("pre").forEach(el => el.innerText = el.innerHTML);
</script>
`
}

func css() string {
	return `
* ,*:before, *:after {
  box-sizing:border-box;
}

body {
  font-family: geneva, verdana, 'lucida sans', 'lucida grande', 'lucida sans unicode',sans-serif;
  color:#345;
}

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
.item {
  border: 1px solid #ddd;
  padding: 1em;

}
.snippet pre {
  overflow-x: auto;
  white-space: pre-wrap;
  white-space: -moz-pre-wrap;
  white-space: -pre-wrap;
  white-space: -o-pre-wrap;
  word-wrap: break-word;
}

.snippet {
  position: relative;
}

@media screen and (min-width: 760px) {

  .snippet .btn-copy {
    top: 0;
    right: 0;
    position: absolute;
  }
}

@media screen and (min-width: 1440px) {

  .container {
    width: 60%;
  }

}
	`

}
