package main

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWriteErrorDoesNotLeakErr verifies that the writeError helper writes only
// the user-facing message into the response body, never the underlying error.
func TestWriteErrorDoesNotLeakErr(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("internal: /etc/secret-path: permission denied")

	writeError(w, http.StatusInternalServerError, "Something went wrong", err)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", res.StatusCode, http.StatusInternalServerError)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Something went wrong") {
		t.Errorf("body missing user message: %q", body)
	}
	for _, leaked := range []string{"/etc/secret-path", "permission denied", "internal:"} {
		if strings.Contains(body, leaked) {
			t.Errorf("body leaked %q from underlying error: %q", leaked, body)
		}
	}
}

// TestImportExecuteHandlerSanitizesError verifies that ImportExecuteHandler
// returns a sanitised body when the underlying ExecuteImport call fails with
// an error containing internal filesystem paths. Previously the handler
// embedded the full error chain in the response with `fmt.Sprintf("Import
// failed: %v", err)`, which leaked destination paths to the user.
func TestImportExecuteHandlerSanitizesError(t *testing.T) {
	tmpDir := t.TempDir()
	importDir := filepath.Join(tmpDir, "import")
	mediaDir := filepath.Join(tmpDir, "media")
	sourceDir := filepath.Join(importDir, "source-disk")
	leakedPath := filepath.Join(mediaDir, "Leaky Film (2020) [Film]", "Disk [Blu-Ray]")

	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Pre-create the destination so ExecuteImport returns
	// "destination already exists: <leakedPath>".
	if err := os.MkdirAll(leakedPath, 0o755); err != nil {
		t.Fatal(err)
	}

	app := NewApp(nil, nil, template.Must(template.New("ignored").Parse("")), mediaDir, importDir)
	app.importScanner = NewImportScanner(importDir)
	session := &ImportSession{
		SourceDir:     &ImportDirectory{Name: "source-disk", Path: sourceDir},
		MediaKind:     Film,
		Title:         "Leaky Film",
		Year:          2020,
		DiskType:      DiskTypeBluRay,
		AddToExisting: false,
	}
	sessionID := app.importSessions.Create(session)

	form := url.Values{"session": {sessionID}}
	req := httptest.NewRequest(http.MethodPost, "/import/execute", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	app.ImportExecuteHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", res.StatusCode, http.StatusInternalServerError)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Import failed") {
		t.Errorf("body missing user message: %q", body)
	}
	if strings.Contains(body, leakedPath) {
		t.Errorf("body leaked destination path %q: %q", leakedPath, body)
	}
	if strings.Contains(body, "destination already exists") {
		t.Errorf("body leaked internal error chain: %q", body)
	}
}

// TestSetTMDBHandlerNoIDDoesNotLeakErr asserts the empty-ID validation path
// returns only the sanitised message.
func TestSetTMDBHandlerNoIDDoesNotLeakErr(t *testing.T) {
	tmpDir := t.TempDir()
	media := Media{Title: "Test", Type: Film, Year: 2020, Path: tmpDir, DiskCount: 1}

	app := NewApp([]Media{media}, nil, template.Must(template.New("ignored").Parse("")), "/media", "")
	app.SetTMDBClient(NewTMDBClient("test-key"))

	form := url.Values{"tmdb_id": {""}}
	req := httptest.NewRequest(http.MethodPost, "/media/test-2020/set-tmdb", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("slug", "test-2020")
	w := httptest.NewRecorder()

	app.SaveTMDBHandler(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "TMDB ID is required") {
		t.Errorf("body missing user message: %q", body)
	}
	// Should not leak any wrapped error details (none in this path, but make
	// sure %v formatting wasn't reintroduced).
	if strings.Contains(body, "%!") || strings.Contains(body, "<nil>") {
		t.Errorf("body contains formatting artefacts: %q", body)
	}
}
