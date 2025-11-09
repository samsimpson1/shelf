package main

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to parse a template file
func parseTemplate(t *testing.T, name string) *template.Template {
	t.Helper()
	tmplPath := filepath.Join("templates", name)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		t.Fatalf("Failed to parse template %s: %v", name, err)
	}
	return tmpl
}

// Helper function to execute a template and return the output
func executeTemplate(t *testing.T, tmpl *template.Template, data interface{}) string {
	t.Helper()
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}
	return buf.String()
}

// TestIndexTemplateParses tests that index.html template parses without syntax errors
func TestIndexTemplateParses(t *testing.T) {
	parseTemplate(t, "index.html")
}

// TestIndexTemplateWithMediaList tests index.html with a valid media list
func TestIndexTemplateWithMediaList(t *testing.T) {
	tmpl := parseTemplate(t, "index.html")

	data := struct {
		MediaList []Media
	}{
		MediaList: []Media{
			{
				Title:     "The Thing",
				Type:      Film,
				Year:      1982,
				DiskCount: 1,
				TMDBID:    "1091",
				Path:      "/test/the-thing",
			},
			{
				Title:     "Better Call Saul",
				Type:      TV,
				Year:      0,
				DiskCount: 12,
				TMDBID:    "60059",
				Path:      "/test/better-call-saul",
			},
		},
	}

	output := executeTemplate(t, tmpl, data)

	// Verify expected content is present
	expectedStrings := []string{
		"The Thing (1982)",
		"Better Call Saul",
		"Film",
		"TV",
		"the-thing-1982",
		"better-call-saul",
		"2 items",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

// TestIndexTemplateWithEmptyList tests index.html with no media items
func TestIndexTemplateWithEmptyList(t *testing.T) {
	tmpl := parseTemplate(t, "index.html")

	data := struct {
		MediaList []Media
	}{
		MediaList: []Media{},
	}

	output := executeTemplate(t, tmpl, data)

	// Verify empty state is shown
	expectedStrings := []string{
		"No Media Found",
		"No media items were found",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

// TestIndexTemplateWithNilMediaList tests index.html with nil media list
func TestIndexTemplateWithNilMediaList(t *testing.T) {
	tmpl := parseTemplate(t, "index.html")

	data := struct {
		MediaList []Media
	}{
		MediaList: nil,
	}

	output := executeTemplate(t, tmpl, data)

	// Verify empty state is shown
	if !strings.Contains(output, "No Media Found") {
		t.Error("Expected output to show empty state for nil media list")
	}
}

// TestIndexTemplateWithSingleItem tests singular "item" vs plural "items"
func TestIndexTemplateWithSingleItem(t *testing.T) {
	tmpl := parseTemplate(t, "index.html")

	data := struct {
		MediaList []Media
	}{
		MediaList: []Media{
			{
				Title:     "Solo Movie",
				Type:      Film,
				Year:      2020,
				DiskCount: 1,
				Path:      "/test/solo",
			},
		},
	}

	output := executeTemplate(t, tmpl, data)

	// Verify singular form
	if !strings.Contains(output, "1 item") {
		t.Error("Expected output to contain '1 item' (singular)")
	}
	if strings.Contains(output, "1 items") {
		t.Error("Expected output NOT to contain '1 items' (incorrect plural)")
	}
}

// TestIndexTemplateDiskCountPluralization tests disk count pluralization
func TestIndexTemplateDiskCountPluralization(t *testing.T) {
	tmpl := parseTemplate(t, "index.html")

	tests := []struct {
		name      string
		diskCount int
		expected  string
	}{
		{"Single disk", 1, "1 disk"},
		{"Multiple disks", 3, "3 disks"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := struct {
				MediaList []Media
			}{
				MediaList: []Media{
					{
						Title:     "Test",
						Type:      Film,
						Year:      2020,
						DiskCount: tt.diskCount,
						Path:      "/test",
					},
				},
			}

			output := executeTemplate(t, tmpl, data)
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain %q", tt.expected)
			}
		})
	}
}

// TestDetailTemplateParses tests that detail.html template parses without syntax errors
func TestDetailTemplateParses(t *testing.T) {
	parseTemplate(t, "detail.html")
}

// TestDetailTemplateWithFullData tests detail.html with complete data
func TestDetailTemplateWithFullData(t *testing.T) {
	tmpl := parseTemplate(t, "detail.html")

	media := &Media{
		Title:     "The Matrix",
		Type:      Film,
		Year:      1999,
		DiskCount: 2,
		TMDBID:    "603",
		Path:      "/test/the-matrix",
	}

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		PlayURLPrefix string
	}{
		Media:         media,
		Description:   "A computer hacker learns about the true nature of reality.",
		Genres:        []string{"Action", "Science Fiction"},
		HasPoster:     true,
		PlayURLPrefix: "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify expected content
	expectedStrings := []string{
		"The Matrix (1999)",
		"Film",
		"1999",
		"2",
		"603",
		"A computer hacker learns about the true nature of reality.",
		"Action",
		"Science Fiction",
		"the-matrix-1999",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

// TestDetailTemplateWithNoPoster tests detail.html without a poster
func TestDetailTemplateWithNoPoster(t *testing.T) {
	tmpl := parseTemplate(t, "detail.html")

	media := &Media{
		Title:     "Test Film",
		Type:      Film,
		Year:      2020,
		DiskCount: 1,
		Path:      "/test/test-film",
	}

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		PlayURLPrefix string
	}{
		Media:         media,
		Description: "Test description",
		Genres:        []string{"Drama"},
		HasPoster:     false,
		PlayURLPrefix: "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify placeholder emoji is shown
	if !strings.Contains(output, "🎬") {
		t.Error("Expected output to contain film emoji placeholder")
	}
}

// TestDetailTemplateWithNoDescription tests detail.html with missing description
func TestDetailTemplateWithNoDescription(t *testing.T) {
	tmpl := parseTemplate(t, "detail.html")

	media := &Media{
		Title:     "Test Film",
		Type:      Film,
		Year:      2020,
		DiskCount: 1,
		Path:      "/test/test-film",
	}

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		PlayURLPrefix string
	}{
		Media:         media,
		Description: "",
		Genres:        []string{},
		HasPoster:     false,
		PlayURLPrefix: "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify "No description available" message
	if !strings.Contains(output, "No description available") {
		t.Error("Expected output to show 'No description available'")
	}
}

// TestDetailTemplateWithEmptyGenres tests detail.html with no genres
func TestDetailTemplateWithEmptyGenres(t *testing.T) {
	tmpl := parseTemplate(t, "detail.html")

	media := &Media{
		Title:     "Test Film",
		Type:      Film,
		Year:      2020,
		DiskCount: 1,
		Path:      "/test/test-film",
	}

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		PlayURLPrefix string
	}{
		Media:         media,
		Description: "Test",
		Genres:        []string{},
		HasPoster:     false,
		PlayURLPrefix: "",
	}

	output := executeTemplate(t, tmpl, data)

	// Should not error, just not show genre section
	if strings.Contains(output, "class=\"genre\"") {
		t.Error("Expected no genre elements when genres list is empty")
	}
}

// TestDetailTemplateWithNoTMDBID tests detail.html without TMDB ID
func TestDetailTemplateWithNoTMDBID(t *testing.T) {
	tmpl := parseTemplate(t, "detail.html")

	media := &Media{
		Title:     "Test Film",
		Type:      Film,
		Year:      2020,
		DiskCount: 1,
		TMDBID:    "",
		Path:      "/test/test-film",
	}

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		PlayURLPrefix string
	}{
		Media:         media,
		Description: "Test",
		Genres:        []string{},
		HasPoster:     false,
		PlayURLPrefix: "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify TMDB section is not shown
	if strings.Contains(output, "<strong>TMDB:</strong>") {
		t.Error("Expected TMDB metadata not to be shown when TMDBID is empty")
	}

	// Verify "Search for TMDB ID" button is shown instead
	if !strings.Contains(output, "Search for TMDB ID") {
		t.Error("Expected 'Search for TMDB ID' button to be shown")
	}
}

// TestDetailTemplateWithTVShow tests detail.html with TV show (no year)
func TestDetailTemplateWithTVShow(t *testing.T) {
	tmpl := parseTemplate(t, "detail.html")

	media := &Media{
		Title:     "Breaking Bad",
		Type:      TV,
		Year:      0,
		DiskCount: 5,
		TMDBID:    "1396",
		Path:      "/test/breaking-bad",
	}

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		PlayURLPrefix string
	}{
		Media:         media,
		Description: "A chemistry teacher turned meth cook.",
		Genres:        []string{"Drama", "Crime"},
		HasPoster:     true,
		PlayURLPrefix: "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify TV show displays correctly (no year in title)
	if strings.Contains(output, "Breaking Bad (") {
		t.Error("Expected TV show title not to include year")
	}

	// Verify TMDB link points to TV endpoint
	if !strings.Contains(output, "themoviedb.org/tv/1396") {
		t.Error("Expected TMDB link to point to TV endpoint")
	}
}

// TestDetailTemplateWithDisks tests detail.html with disk listing
func TestDetailTemplateWithDisks(t *testing.T) {
	tmpl := parseTemplate(t, "detail.html")

	media := &Media{
		Title:     "The Thing",
		Type:      Film,
		Year:      1982,
		DiskCount: 2,
		Disks: []Disk{
			{Name: "Disk 1", Format: "Blu-Ray", SizeGB: 45.2, Path: "/media/the-thing/Disk [Blu-Ray]"},
			{Name: "Disk 2", Format: "DVD", SizeGB: 4.7, Path: "/media/the-thing/Disk 2 [DVD]"},
		},
		TMDBID: "1091",
		Path:   "/test/the-thing",
	}

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		PlayURLPrefix string
	}{
		Media:         media,
		Description:   "Scientists in Antarctica discover an alien.",
		Genres:        []string{"Horror", "Science Fiction"},
		HasPoster:     true,
		PlayURLPrefix: "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify disk listing is shown
	expectedStrings := []string{
		"Disks",
		"Disk 1",
		"Disk 2",
		"Blu-Ray",
		"DVD",
		"45.2 GB",
		"4.7 GB",
		"Name",
		"Format",
		"Size",
		"Action",
		"Copy VLC Command",
		"Copy MPV Command",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}

	// Verify empty state is NOT shown
	if strings.Contains(output, "No disks found") {
		t.Error("Expected 'No disks found' not to be shown when disks are present")
	}
}

// TestDetailTemplateWithNoDisks tests detail.html with no disks
func TestDetailTemplateWithNoDisks(t *testing.T) {
	tmpl := parseTemplate(t, "detail.html")

	media := &Media{
		Title:     "Test Film",
		Type:      Film,
		Year:      2020,
		DiskCount: 0,
		Disks:     []Disk{},
		Path:      "/test/test-film",
	}

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		PlayURLPrefix string
	}{
		Media:         media,
		Description: "Test",
		Genres:        []string{},
		HasPoster:     false,
		PlayURLPrefix: "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify empty state is shown
	if !strings.Contains(output, "No disks found") {
		t.Error("Expected 'No disks found' to be shown when no disks are present")
	}

	// Verify table is NOT shown
	if strings.Contains(output, "<table") {
		t.Error("Expected disk table not to be shown when no disks are present")
	}
}

// TestDetailTemplateWithTVDisks tests detail.html with TV show disks
func TestDetailTemplateWithTVDisks(t *testing.T) {
	tmpl := parseTemplate(t, "detail.html")

	media := &Media{
		Title:     "Better Call Saul",
		Type:      TV,
		Year:      0,
		DiskCount: 2,
		Disks: []Disk{
			{Name: "Series 1 Disk 1", Format: "Blu-Ray", SizeGB: 23.5, Path: "/media/better-call-saul/Series 1 Disk 1 [Blu-Ray]"},
			{Name: "Series 1 Disk 2", Format: "Blu-Ray UHD", SizeGB: 66.8, Path: "/media/better-call-saul/Series 1 Disk 2 [Blu-Ray UHD]"},
		},
		TMDBID: "60059",
		Path:   "/test/better-call-saul",
	}

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		PlayURLPrefix string
	}{
		Media:         media,
		Description: "The trials and tribulations of criminal lawyer Jimmy McGill.",
		Genres:        []string{"Drama", "Crime"},
		HasPoster:     true,
		PlayURLPrefix: "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify TV disk naming is correct
	expectedStrings := []string{
		"Series 1 Disk 1",
		"Series 1 Disk 2",
		"Blu-Ray UHD",
		"23.5 GB",
		"66.8 GB",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

// TestSearchTemplateParses tests that search.html template parses without syntax errors
func TestSearchTemplateParses(t *testing.T) {
	parseTemplate(t, "search.html")
}

// TestSearchTemplateWithMovieResults tests search.html with movie search results
func TestSearchTemplateWithMovieResults(t *testing.T) {
	tmpl := parseTemplate(t, "search.html")

	media := &Media{
		Title:     "The Thing",
		Type:      Film,
		Year:      1982,
		DiskCount: 1,
		Path:      "/test/the-thing",
	}

	results := []MovieSearchResult{
		{
			ID:          1091,
			Title:       "The Thing",
			ReleaseDate: "1982-06-25",
			Overview:    "A research team in Antarctica is hunted by a shape-shifting alien.",
			PosterPath:  "/tzGY49kseSE9QAKk47uuDGwnSCu.jpg",
			Popularity:  45.2,
		},
		{
			ID:          2254,
			Title:       "The Thing from Another World",
			ReleaseDate: "1951-04-27",
			Overview:    "Scientists discover an alien spaceship in the Arctic.",
			PosterPath:  "/abc123.jpg",
			Popularity:  12.5,
		},
	}

	data := struct {
		Media   *Media
		Query   string
		Year    int
		Results interface{}
		Error   string
	}{
		Media:   media,
		Query:   "The Thing",
		Year:    1982,
		Results: results,
		Error:   "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify search results are displayed
	expectedStrings := []string{
		"The Thing",
		"1982-06-25",
		"A research team in Antarctica",
		"The Thing from Another World",
		"1951-04-27",
		"Search Results (2)",
		"Popularity: 45.2",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

// TestSearchTemplateWithTVResults tests search.html with TV search results
func TestSearchTemplateWithTVResults(t *testing.T) {
	tmpl := parseTemplate(t, "search.html")

	media := &Media{
		Title:     "Breaking Bad",
		Type:      TV,
		Year:      0,
		DiskCount: 5,
		Path:      "/test/breaking-bad",
	}

	results := []TVSearchResult{
		{
			ID:           1396,
			Name:         "Breaking Bad",
			FirstAirDate: "2008-01-20",
			Overview:     "A chemistry teacher turned meth cook.",
			PosterPath:   "/zzz.jpg",
			Popularity:   100.5,
		},
	}

	data := struct {
		Media   *Media
		Query   string
		Year    int
		Results interface{}
		Error   string
	}{
		Media:   media,
		Query:   "Breaking Bad",
		Year:    0,
		Results: results,
		Error:   "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify TV results are displayed
	expectedStrings := []string{
		"Breaking Bad",
		"2008-01-20",
		"A chemistry teacher turned meth cook.",
		"Search Results (1)",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

// TestSearchTemplateWithNoResults tests search.html with empty results
func TestSearchTemplateWithNoResults(t *testing.T) {
	tmpl := parseTemplate(t, "search.html")

	media := &Media{
		Title:     "Nonexistent Movie",
		Type:      Film,
		Year:      2099,
		DiskCount: 1,
		Path:      "/test/nonexistent",
	}

	data := struct {
		Media   *Media
		Query   string
		Year    int
		Results interface{}
		Error   string
	}{
		Media:   media,
		Query:   "Nonexistent Movie",
		Year:    2099,
		Results: []MovieSearchResult{},
		Error:   "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify empty state is shown
	expectedStrings := []string{
		"No Results Found",
		"Try searching with a different title or year",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

// TestSearchTemplateWithError tests search.html with an error message
func TestSearchTemplateWithError(t *testing.T) {
	tmpl := parseTemplate(t, "search.html")

	media := &Media{
		Title:     "Test Movie",
		Type:      Film,
		Year:      2020,
		DiskCount: 1,
		Path:      "/test/test-movie",
	}

	data := struct {
		Media   *Media
		Query   string
		Year    int
		Results interface{}
		Error   string
	}{
		Media:   media,
		Query:   "",
		Year:    0,
		Results: nil,
		Error:   "TMDB API error: invalid API key",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify error message is displayed
	if !strings.Contains(output, "TMDB API error: invalid API key") {
		t.Error("Expected output to contain error message")
	}
}

// TestSearchTemplateWithNoPoster tests search.html with results missing poster
func TestSearchTemplateWithNoPoster(t *testing.T) {
	tmpl := parseTemplate(t, "search.html")

	media := &Media{
		Title:     "Test Movie",
		Type:      Film,
		Year:      2020,
		DiskCount: 1,
		Path:      "/test/test-movie",
	}

	results := []MovieSearchResult{
		{
			ID:          123,
			Title:       "Movie Without Poster",
			ReleaseDate: "2020-01-01",
			Overview:    "Test overview",
			PosterPath:  "",
			Popularity:  5.0,
		},
	}

	data := struct {
		Media   *Media
		Query   string
		Year    int
		Results interface{}
		Error   string
	}{
		Media:   media,
		Query:   "test",
		Year:    0,
		Results: results,
		Error:   "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify placeholder is shown
	if !strings.Contains(output, "result-placeholder") {
		t.Error("Expected output to contain placeholder for missing poster")
	}
}

// TestConfirmTemplateParses tests that confirm.html template parses without syntax errors
func TestConfirmTemplateParses(t *testing.T) {
	parseTemplate(t, "confirm.html")
}

// TestConfirmTemplateWithMovieMatch tests confirm.html with a movie match
func TestConfirmTemplateWithMovieMatch(t *testing.T) {
	tmpl := parseTemplate(t, "confirm.html")

	media := &Media{
		Title:     "Fight Club",
		Type:      Film,
		Year:      1999,
		DiskCount: 1,
		TMDBID:    "",
		Path:      "/test/fight-club",
	}

	tmdbMatch := MovieSearchResult{
		ID:          550,
		Title:       "Fight Club",
		ReleaseDate: "1999-10-15",
		Overview:    "An insomniac office worker and a devil-may-care soap maker form an underground fight club.",
		PosterPath:  "/pB8BM7pdSp6B6Ih7QZ4DrQ3PmJK.jpg",
		Popularity:  85.3,
	}

	data := struct {
		Media       *Media
		TMDBID      string
		TMDBMatch   interface{}
		Query       string
		Description string
		HasPoster   bool
		Error       string
	}{
		Media:       media,
		TMDBID:      "550",
		TMDBMatch:   tmdbMatch,
		Query:       "Fight Club",
		Description: "",
		HasPoster:   false,
		Error:       "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify comparison is shown
	expectedStrings := []string{
		"Fight Club",
		"1999",
		"550",
		"1999-10-15",
		"An insomniac office worker",
		"Confirm and Save TMDB ID",
		"Download metadata now",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

// TestConfirmTemplateWithTVMatch tests confirm.html with a TV show match
func TestConfirmTemplateWithTVMatch(t *testing.T) {
	tmpl := parseTemplate(t, "confirm.html")

	media := &Media{
		Title:     "The Wire",
		Type:      TV,
		Year:      0,
		DiskCount: 5,
		TMDBID:    "",
		Path:      "/test/the-wire",
	}

	tmdbMatch := TVSearchResult{
		ID:           1438,
		Name:         "The Wire",
		FirstAirDate: "2002-06-02",
		Overview:     "A look at the drug scene in Baltimore through the eyes of law enforcement.",
		PosterPath:   "/the-wire.jpg",
		Popularity:   50.2,
	}

	data := struct {
		Media       *Media
		TMDBID      string
		TMDBMatch   interface{}
		Query       string
		Description string
		HasPoster   bool
		Error       string
	}{
		Media:       media,
		TMDBID:      "1438",
		TMDBMatch:   tmdbMatch,
		Query:       "The Wire",
		Description: "",
		HasPoster:   false,
		Error:       "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify TV-specific fields are shown
	expectedStrings := []string{
		"The Wire",
		"1438",
		"2002-06-02",
		"First Air Date",
		"A look at the drug scene in Baltimore",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

// TestConfirmTemplateWithExistingTMDBID tests confirm.html when media already has TMDB ID
func TestConfirmTemplateWithExistingTMDBID(t *testing.T) {
	tmpl := parseTemplate(t, "confirm.html")

	media := &Media{
		Title:     "The Matrix",
		Type:      Film,
		Year:      1999,
		DiskCount: 1,
		TMDBID:    "603",
		Path:      "/test/the-matrix",
	}

	tmdbMatch := MovieSearchResult{
		ID:          604,
		Title:       "The Matrix Reloaded",
		ReleaseDate: "2003-05-15",
		Overview:    "Six months after the first movie.",
		PosterPath:  "/reloaded.jpg",
		Popularity:  40.0,
	}

	data := struct {
		Media       *Media
		TMDBID      string
		TMDBMatch   interface{}
		Query       string
		Description string
		HasPoster   bool
		Error       string
	}{
		Media:       media,
		TMDBID:      "604",
		TMDBMatch:   tmdbMatch,
		Query:       "Matrix",
		Description: "Existing description",
		HasPoster:   true,
		Error:       "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify warning is shown about replacing existing ID
	expectedStrings := []string{
		"Warning",
		"already has a TMDB ID",
		"603",
		"replace the existing metadata",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

// TestConfirmTemplateWithError tests confirm.html with an error
func TestConfirmTemplateWithError(t *testing.T) {
	tmpl := parseTemplate(t, "confirm.html")

	media := &Media{
		Title:     "Test Movie",
		Type:      Film,
		Year:      2020,
		DiskCount: 1,
		Path:      "/test/test-movie",
	}

	data := struct {
		Media       *Media
		TMDBID      string
		TMDBMatch   interface{}
		Query       string
		Description string
		HasPoster   bool
		Error       string
	}{
		Media:       media,
		TMDBID:      "999999",
		TMDBMatch:   nil,
		Query:       "",
		Description: "",
		HasPoster:   false,
		Error:       "Failed to fetch TMDB details",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify error is displayed
	if !strings.Contains(output, "Failed to fetch TMDB details") {
		t.Error("Expected output to contain error message")
	}
}

// TestConfirmTemplateWithNoPosterOnBothSides tests confirm.html with no posters
func TestConfirmTemplateWithNoPosterOnBothSides(t *testing.T) {
	tmpl := parseTemplate(t, "confirm.html")

	media := &Media{
		Title:     "Test Film",
		Type:      Film,
		Year:      2020,
		DiskCount: 1,
		Path:      "/test/test-film",
	}

	tmdbMatch := MovieSearchResult{
		ID:          123,
		Title:       "Test Film",
		ReleaseDate: "2020-01-01",
		Overview:    "Test overview",
		PosterPath:  "",
		Popularity:  1.0,
	}

	data := struct {
		Media       *Media
		TMDBID      string
		TMDBMatch   interface{}
		Query       string
		Description string
		HasPoster   bool
		Error       string
	}{
		Media:       media,
		TMDBID:      "123",
		TMDBMatch:   tmdbMatch,
		Query:       "test",
		Description: "",
		HasPoster:   false,
		Error:       "",
	}

	output := executeTemplate(t, tmpl, data)

	// Verify placeholders are shown for both sides
	placeholderCount := strings.Count(output, "media-placeholder")
	if placeholderCount < 2 {
		t.Errorf("Expected at least 2 placeholders, got %d", placeholderCount)
	}
}

// TestAllTemplatesWithTypeComparison tests that type comparisons work for Film vs TV
func TestAllTemplatesWithTypeComparison(t *testing.T) {
	templates := []string{"index.html", "detail.html", "search.html", "confirm.html"}

	for _, tmplName := range templates {
		t.Run(tmplName, func(t *testing.T) {
			tmpl := parseTemplate(t, tmplName)

			// Test with Film type
			filmMedia := &Media{
				Title:     "Test Film",
				Type:      Film,
				Year:      2020,
				DiskCount: 1,
				Path:      "/test/film",
			}

			// Test with TV type
			tvMedia := &Media{
				Title:     "Test TV",
				Type:      TV,
				Year:      0,
				DiskCount: 5,
				Path:      "/test/tv",
			}

			// Create appropriate data structures for each template
			switch tmplName {
			case "index.html":
				filmData := struct{ MediaList []Media }{
					MediaList: []Media{*filmMedia},
				}
				tvData := struct{ MediaList []Media }{
					MediaList: []Media{*tvMedia},
				}

				filmOutput := executeTemplate(t, tmpl, filmData)
				tvOutput := executeTemplate(t, tmpl, tvData)

				// Both should render without error
				if len(filmOutput) == 0 || len(tvOutput) == 0 {
					t.Error("Template execution produced empty output")
				}

			case "detail.html":
				filmData := struct {
					Media         *Media
					Description   string
					Genres        []string
					HasPoster     bool
					PlayURLPrefix string
				}{Media: filmMedia, Description: "Test", Genres: []string{}, HasPoster: false, PlayURLPrefix: ""}

				tvData := struct {
					Media         *Media
					Description   string
					Genres        []string
					HasPoster     bool
					PlayURLPrefix string
				}{Media: tvMedia, Description: "Test", Genres: []string{}, HasPoster: false, PlayURLPrefix: ""}

				filmOutput := executeTemplate(t, tmpl, filmData)
				tvOutput := executeTemplate(t, tmpl, tvData)

				if len(filmOutput) == 0 || len(tvOutput) == 0 {
					t.Error("Template execution produced empty output")
				}

			case "search.html":
				filmData := struct {
					Media   *Media
					Query   string
					Year    int
					Results interface{}
					Error   string
				}{Media: filmMedia, Query: "", Year: 0, Results: nil, Error: ""}

				tvData := struct {
					Media   *Media
					Query   string
					Year    int
					Results interface{}
					Error   string
				}{Media: tvMedia, Query: "", Year: 0, Results: nil, Error: ""}

				filmOutput := executeTemplate(t, tmpl, filmData)
				tvOutput := executeTemplate(t, tmpl, tvData)

				if len(filmOutput) == 0 || len(tvOutput) == 0 {
					t.Error("Template execution produced empty output")
				}

			case "confirm.html":
				filmData := struct {
					Media       *Media
					TMDBID      string
					TMDBMatch   interface{}
					Query       string
					Description string
					HasPoster   bool
					Error       string
				}{Media: filmMedia, TMDBID: "123", TMDBMatch: nil, Query: "", Description: "", HasPoster: false, Error: ""}

				tvData := struct {
					Media       *Media
					TMDBID      string
					TMDBMatch   interface{}
					Query       string
					Description string
					HasPoster   bool
					Error       string
				}{Media: tvMedia, TMDBID: "456", TMDBMatch: nil, Query: "", Description: "", HasPoster: false, Error: ""}

				filmOutput := executeTemplate(t, tmpl, filmData)
				tvOutput := executeTemplate(t, tmpl, tvData)

				if len(filmOutput) == 0 || len(tvOutput) == 0 {
					t.Error("Template execution produced empty output")
				}
			}
		})
	}
}
