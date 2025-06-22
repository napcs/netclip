package netclip_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"netclip"

	"github.com/stretchr/testify/assert"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(netclip.IndexHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Saved clips")
}

func TestSaveHandler(t *testing.T) {
	// Set up the test request
	formData := url.Values{}
	formData.Set("text", "testing123")

	req, err := http.NewRequest("POST", "/save", strings.NewReader(formData.Encode()))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(netclip.SaveHandler)
	handler.ServeHTTP(rr, req)

	// Check the response
	assert.Equal(t, http.StatusSeeOther, rr.Code)
	assert.Equal(t, "/", rr.Header().Get("Location"))

	// Follow the redirect and check that the record was saved
	req, err = http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(netclip.IndexHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "testing123")
}

// TestDeleteHandler needs to run a full series of requests because
// the datastore is in-memory and not exposed outside of the server, so
// there's no way to check that records saved without looking at the body
// and there's no way to get the key of the last saved record because we
// redirect to the index page rather than redirecting to the page with the key.
//
// Thus the sequence is:
// * Post a value
// * Follow the redirect to /
// * Read the HTML response to extract the key
// * Use that key to make a delete request.
func TestDeleteHandler(t *testing.T) {
	// Set up the test request
	formData := url.Values{}
	formData.Set("text", "testing123")

	req, err := http.NewRequest("POST", "/save", strings.NewReader(formData.Encode()))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	// Call the save handler to add a record to the datastore
	handler := http.HandlerFunc(netclip.SaveHandler)
	handler.ServeHTTP(rr, req)

	// Check the response
	assert.Equal(t, http.StatusSeeOther, rr.Code)
	assert.Equal(t, "/", rr.Header().Get("Location"))

	// Follow the redirect and get the key of the saved record
	req, err = http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(netclip.IndexHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Extract the key from the index page
	body := rr.Body.String()
	re := regexp.MustCompile(`<input type="hidden" value="(.*)" name="key">`)
	matches := re.FindStringSubmatch(body)
	key := matches[1]

	// Set up the delete request
	formData = url.Values{}
	formData.Set("key", key)

	req, err = http.NewRequest("POST", "/delete", strings.NewReader(formData.Encode()))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	// Call the delete handler to remove the record from the datastore
	handler = http.HandlerFunc(netclip.DeleteHandler)
	handler.ServeHTTP(rr, req)

	// Check the response
	assert.Equal(t, http.StatusSeeOther, rr.Code)
	assert.Equal(t, "/", rr.Header().Get("Location"))

	// Follow the redirect and check that the record was deleted
	req, err = http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(netclip.IndexHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotContains(t, rr.Body.String(), key)
}
