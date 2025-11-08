package main

import "testing"

func TestMediaTypeString(t *testing.T) {
	tests := []struct {
		name     string
		mt       MediaType
		expected string
	}{
		{"Film type", Film, "Film"},
		{"TV type", TV, "TV"},
		{"Unknown type", MediaType(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.mt.String()
			if result != tt.expected {
				t.Errorf("MediaType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMediaDisplayTitle(t *testing.T) {
	tests := []struct {
		name     string
		media    Media
		expected string
	}{
		{
			name: "Film with year",
			media: Media{
				Title: "War of the Worlds",
				Type:  Film,
				Year:  2025,
			},
			expected: "War of the Worlds (2025)",
		},
		{
			name: "Film without year",
			media: Media{
				Title: "Some Film",
				Type:  Film,
				Year:  0,
			},
			expected: "Some Film",
		},
		{
			name: "TV show",
			media: Media{
				Title: "Better Call Saul",
				Type:  TV,
				Year:  0,
			},
			expected: "Better Call Saul",
		},
		{
			name: "TV show with year (should ignore)",
			media: Media{
				Title: "Breaking Bad",
				Type:  TV,
				Year:  2008,
			},
			expected: "Breaking Bad",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.media.DisplayTitle()
			if result != tt.expected {
				t.Errorf("Media.DisplayTitle() = %v, want %v", result, tt.expected)
			}
		})
	}
}
