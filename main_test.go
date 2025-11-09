package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	var buf bytes.Buffer
	printHelp(&buf)

	output := buf.String()

	// Check that all required elements are present
	tests := []struct {
		name     string
		contains string
	}{
		{"Title", "Media Backup Manager"},
		{"Usage section", "Usage:"},
		{"MEDIA_DIR env var", "MEDIA_DIR"},
		{"MEDIA_DIR description", "Path to media backup directory"},
		{"MEDIA_DIR default", "/home/sam/Scratch/media/backup"},
		{"PORT env var", "PORT"},
		{"PORT description", "HTTP server port"},
		{"PORT default", "8080"},
		{"TMDB_API_KEY env var", "TMDB_API_KEY"},
		{"TMDB_API_KEY description", "TMDB API key"},
		{"Optional indicator", "optional"},
		{"DEV_MODE env var", "DEV_MODE"},
		{"DEV_MODE description", "Development mode"},
		{"PLAY_URL_PREFIX env var", "PLAY_URL_PREFIX"},
		{"PLAY_URL_PREFIX description", "URL prefix for VLC play commands"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Help output missing %q\nGot:\n%s", tt.contains, output)
			}
		})
	}
}

func TestShouldShowHelp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "No arguments",
			args:     []string{"shelf"},
			expected: false,
		},
		{
			name:     "Help with single dash",
			args:     []string{"shelf", "-help"},
			expected: true,
		},
		{
			name:     "Help with double dash",
			args:     []string{"shelf", "--help"},
			expected: true,
		},
		{
			name:     "Help with -h",
			args:     []string{"shelf", "-h"},
			expected: true,
		},
		{
			name:     "Other arguments",
			args:     []string{"shelf", "something"},
			expected: false,
		},
		{
			name:     "Multiple arguments with help",
			args:     []string{"shelf", "-help", "other"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldShowHelp(tt.args)
			if got != tt.expected {
				t.Errorf("shouldShowHelp(%v) = %v, want %v", tt.args, got, tt.expected)
			}
		})
	}
}
