package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	// Create a simple template for testing
	tmpl := template.Must(template.New("index.html").Parse(`
<!DOCTYPE html>
<html>
<body>
<h1>Media Backup Manager</h1>
<table>
{{range .MediaList}}
<tr>
<td>{{.DisplayTitle}}</td>
<td>{{.Type}}</td>
<td>{{.DiskCount}}</td>
<td>{{.TMDBID}}</td>
</tr>
{{end}}
</table>
</body>
</html>
`))

	tests := []struct {
		name           string
		mediaList      []Media
		expectedStatus int
		expectedInBody []string
	}{
		{
			name: "Multiple media items",
			mediaList: []Media{
				{
					Title:     "War of the Worlds",
					Type:      Film,
					Year:      2025,
					DiskCount: 1,
					TMDBID:    "755898",
					Path:      "/test/path1",
				},
				{
					Title:     "Better Call Saul",
					Type:      TV,
					Year:      0,
					DiskCount: 5,
					TMDBID:    "60059",
					Path:      "/test/path2",
				},
			},
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				"Media Backup Manager",
				"War of the Worlds (2025)",
				"Better Call Saul",
				"Film",
				"TV",
				"755898",
				"60059",
			},
		},
		{
			name:           "Empty media list",
			mediaList:      []Media{},
			expectedStatus: http.StatusOK,
			expectedInBody: []string{"Media Backup Manager"},
		},
		{
			name: "Film without TMDB ID",
			mediaList: []Media{
				{
					Title:     "Unknown Film",
					Type:      Film,
					Year:      2020,
					DiskCount: 1,
					TMDBID:    "",
					Path:      "/test/path",
				},
			},
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				"Unknown Film (2020)",
				"Film",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(tt.mediaList, tmpl, "/test/media")

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			app.IndexHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("IndexHandler() status = %v, want %v", res.StatusCode, tt.expectedStatus)
			}

			body := w.Body.String()
			for _, expected := range tt.expectedInBody {
				if !strings.Contains(body, expected) {
					t.Errorf("IndexHandler() body does not contain %q", expected)
				}
			}
		})
	}
}

func TestIndexHandlerSorting(t *testing.T) {
	// Create a simple template for testing
	tmpl := template.Must(template.New("index.html").Parse(`{{range .MediaList}}{{.Title}},{{end}}`))

	mediaList := []Media{
		{Title: "Zebra Show", Type: TV, DiskCount: 1, Path: "/test/z"},
		{Title: "Alpha Film", Type: Film, Year: 2020, DiskCount: 1, Path: "/test/a"},
		{Title: "Beta Show", Type: TV, DiskCount: 1, Path: "/test/b"},
		{Title: "Gamma Film", Type: Film, Year: 2021, DiskCount: 1, Path: "/test/g"},
	}

	app := NewApp(mediaList, tmpl, "/test/media")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	app.IndexHandler(w, req)

	body := w.Body.String()

	// Expected order: Films first (alphabetically), then TV shows (alphabetically)
	// Alpha Film, Gamma Film, Beta Show, Zebra Show
	expected := "Alpha Film,Gamma Film,Beta Show,Zebra Show,"
	if body != expected {
		t.Errorf("IndexHandler() sorting = %q, want %q", body, expected)
	}
}

func TestNewApp(t *testing.T) {
	mediaList := []Media{
		{Title: "Test", Type: Film, Year: 2020, DiskCount: 1, Path: "/test"},
	}
	tmpl := template.Must(template.New("test").Parse("test"))

	app := NewApp(mediaList, tmpl, "/test/media")

	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
	if len(app.mediaList) != 1 {
		t.Errorf("NewApp() mediaList length = %v, want 1", len(app.mediaList))
	}
	if app.templates == nil {
		t.Error("NewApp() templates is nil")
	}
}
