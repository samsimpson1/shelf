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
			app := NewApp(tt.mediaList, tmpl, "/test/media", "")

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

	app := NewApp(mediaList, tmpl, "/test/media", "")

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

	app := NewApp(mediaList, tmpl, "/test/media", "")

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

func TestDetailHandler(t *testing.T) {
	// Create a realistic detail template
	tmpl := template.Must(template.New("detail.html").Parse(`
<!DOCTYPE html>
<html>
<body>
<h1>{{.Media.DisplayTitle}}</h1>
<p>Type: {{.Media.Type}}</p>
<p>Disks: {{.Media.DiskCount}}</p>
{{if .Media.Disks}}
<table>
<tr><th>Name</th><th>Format</th><th>Size</th></tr>
{{range .Media.Disks}}
<tr><td>{{.Name}}</td><td>{{.Format}}</td><td>{{printf "%.1f GB" .SizeGB}}</td></tr>
{{end}}
</table>
{{end}}
<p>Description: {{.Description}}</p>
{{range .Genres}}<span>{{.}}</span>{{end}}
</body>
</html>
`))

	tests := []struct {
		name           string
		mediaList      []Media
		requestPath    string
		expectedStatus int
		expectedInBody []string
		notInBody      []string
	}{
		{
			name: "Film with disks",
			mediaList: []Media{
				{
					Title:     "The Thing",
					Type:      Film,
					Year:      1982,
					DiskCount: 2,
					Disks: []Disk{
						{Name: "Disk 1", Format: "Blu-Ray", SizeGB: 45.2},
						{Name: "Disk 2", Format: "DVD", SizeGB: 4.7},
					},
					TMDBID: "1091",
					Path:   "/test/the-thing",
				},
			},
			requestPath:    "/media/the-thing-1982",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				"The Thing (1982)",
				"Film",
				"Disk 1",
				"Disk 2",
				"Blu-Ray",
				"DVD",
				"45.2 GB",
				"4.7 GB",
			},
		},
		{
			name: "TV show with disks",
			mediaList: []Media{
				{
					Title:     "Better Call Saul",
					Type:      TV,
					Year:      0,
					DiskCount: 2,
					Disks: []Disk{
						{Name: "Series 1 Disk 1", Format: "Blu-Ray", SizeGB: 23.5},
						{Name: "Series 1 Disk 2", Format: "Blu-Ray UHD", SizeGB: 66.8},
					},
					TMDBID: "60059",
					Path:   "/test/better-call-saul",
				},
			},
			requestPath:    "/media/better-call-saul",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				"Better Call Saul",
				"TV",
				"Series 1 Disk 1",
				"Series 1 Disk 2",
				"Blu-Ray UHD",
				"23.5 GB",
				"66.8 GB",
			},
		},
		{
			name: "Film with no disks",
			mediaList: []Media{
				{
					Title:     "Empty Film",
					Type:      Film,
					Year:      2020,
					DiskCount: 0,
					Disks:     []Disk{},
					Path:      "/test/empty-film",
				},
			},
			requestPath:    "/media/empty-film-2020",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				"Empty Film (2020)",
				"Film",
			},
			notInBody: []string{
				"<table>",
			},
		},
		{
			name:           "Invalid slug - not found",
			mediaList:      []Media{},
			requestPath:    "/media/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Empty slug",
			mediaList:      []Media{},
			requestPath:    "/media/",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(tt.mediaList, tmpl, "/test/media", "")

			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			w := httptest.NewRecorder()

			app.DetailHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("DetailHandler() status = %v, want %v", res.StatusCode, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				body := w.Body.String()
				for _, expected := range tt.expectedInBody {
					if !strings.Contains(body, expected) {
						t.Errorf("DetailHandler() body does not contain %q", expected)
					}
				}
				for _, notExpected := range tt.notInBody {
					if strings.Contains(body, notExpected) {
						t.Errorf("DetailHandler() body should not contain %q", notExpected)
					}
				}
			}
		})
	}
}

func TestFindMediaBySlug(t *testing.T) {
	mediaList := []Media{
		{Title: "The Thing", Type: Film, Year: 1982, Path: "/test/thing"},
		{Title: "Better Call Saul", Type: TV, Year: 0, Path: "/test/bcs"},
	}

	tmpl := template.Must(template.New("test").Parse("test"))
	app := NewApp(mediaList, tmpl, "/test/media", "")

	tests := []struct {
		name      string
		slug      string
		wantTitle string
		wantNil   bool
	}{
		{
			name:      "Find film by slug",
			slug:      "the-thing-1982",
			wantTitle: "The Thing",
			wantNil:   false,
		},
		{
			name:      "Find TV show by slug",
			slug:      "better-call-saul",
			wantTitle: "Better Call Saul",
			wantNil:   false,
		},
		{
			name:    "Nonexistent slug",
			slug:    "nonexistent",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			media := app.findMediaBySlug(tt.slug)

			if tt.wantNil {
				if media != nil {
					t.Errorf("findMediaBySlug(%q) = %v, want nil", tt.slug, media)
				}
			} else {
				if media == nil {
					t.Errorf("findMediaBySlug(%q) = nil, want media", tt.slug)
				} else if media.Title != tt.wantTitle {
					t.Errorf("findMediaBySlug(%q).Title = %q, want %q", tt.slug, media.Title, tt.wantTitle)
				}
			}
		})
	}
}
