package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadDescription(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "description.txt"), []byte("  hello world  \n"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := NewMetadataLoader().LoadDescription(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello world" {
		t.Errorf("got %q, want %q", got, "hello world")
	}
}

func TestLoadDescriptionMissing(t *testing.T) {
	_, err := NewMetadataLoader().LoadDescription(t.TempDir())
	if !errors.Is(err, ErrMetadataNotFound) {
		t.Errorf("want ErrMetadataNotFound, got %v", err)
	}
}

func TestLoadDescriptionPropagatesIOError(t *testing.T) {
	if runtime.GOOS == "windows" || os.Getuid() == 0 {
		t.Skip("permission test requires non-root POSIX")
	}
	dir := t.TempDir()
	desc := filepath.Join(dir, "description.txt")
	if err := os.WriteFile(desc, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(desc, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(desc, 0644) })

	_, err := NewMetadataLoader().LoadDescription(dir)
	if err == nil {
		t.Fatal("expected error from unreadable file, got nil")
	}
	if errors.Is(err, ErrMetadataNotFound) {
		t.Errorf("permission error misclassified as ErrMetadataNotFound: %v", err)
	}
}

func TestLoadGenres(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "genre.txt"), []byte("Action, Drama, Crime\n"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := NewMetadataLoader().LoadGenres(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"Action", "Drama", "Crime"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestLoadGenresEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "genre.txt"), []byte("   \n"), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := NewMetadataLoader().LoadGenres(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}

func TestLoadGenresMissing(t *testing.T) {
	_, err := NewMetadataLoader().LoadGenres(t.TempDir())
	if !errors.Is(err, ErrMetadataNotFound) {
		t.Errorf("want ErrMetadataNotFound, got %v", err)
	}
}

func TestLoadTrackList(t *testing.T) {
	dir := t.TempDir()
	tl := &TrackList{Media: []Medium{{Position: 1, Format: "CD", Tracks: []Track{{Position: "1", Title: "T", Duration: 180000}}}}}
	if err := SaveTrackList(tl, dir); err != nil {
		t.Fatal(err)
	}

	got, err := NewMetadataLoader().LoadTrackList(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || len(got.Media) != 1 || got.Media[0].Tracks[0].Title != "T" {
		t.Errorf("unexpected track list: %+v", got)
	}
}

func TestLoadTrackListMissing(t *testing.T) {
	_, err := NewMetadataLoader().LoadTrackList(t.TempDir())
	if !errors.Is(err, ErrMetadataNotFound) {
		t.Errorf("want ErrMetadataNotFound, got %v", err)
	}
}

func TestLoadTrackListMalformed(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "tracks.json"), []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := NewMetadataLoader().LoadTrackList(dir)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if errors.Is(err, ErrMetadataNotFound) {
		t.Errorf("parse error misclassified as ErrMetadataNotFound: %v", err)
	}
}

func TestFindPoster(t *testing.T) {
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".webp"} {
		t.Run(ext, func(t *testing.T) {
			dir := t.TempDir()
			p := filepath.Join(dir, "poster"+ext)
			if err := os.WriteFile(p, []byte("x"), 0644); err != nil {
				t.Fatal(err)
			}
			got, found, err := NewMetadataLoader().FindPoster(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !found {
				t.Fatal("expected poster to be found")
			}
			if got != p {
				t.Errorf("got %q, want %q", got, p)
			}
		})
	}
}

func TestFindPosterMissing(t *testing.T) {
	_, found, err := NewMetadataLoader().FindPoster(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Error("expected no poster")
	}
}

func TestFindPosterIgnoresDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "poster.jpg"), 0755); err != nil {
		t.Fatal(err)
	}
	_, found, err := NewMetadataLoader().FindPoster(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Error("directory named poster.jpg should not be reported as poster")
	}
}
