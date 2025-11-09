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

func TestDiskPlayCommand(t *testing.T) {
	tests := []struct {
		name     string
		disk     Disk
		prefix   string
		expected string
	}{
		{
			name: "Blu-Ray disk with no prefix",
			disk: Disk{
				Name:   "Disk 1",
				Format: "Blu-Ray",
				Path:   "/media/War of the Worlds (2025) [Film]/Disk [Blu-Ray]",
			},
			prefix:   "",
			expected: "vlc bluray:///media/War of the Worlds (2025) [Film]/Disk [Blu-Ray]",
		},
		{
			name: "Blu-Ray UHD disk with no prefix",
			disk: Disk{
				Name:   "Disk 1",
				Format: "Blu-Ray UHD",
				Path:   "/media/The Thing (1982) [Film]/Disk [Blu-Ray UHD]",
			},
			prefix:   "",
			expected: "vlc bluray:///media/The Thing (1982) [Film]/Disk [Blu-Ray UHD]",
		},
		{
			name: "DVD disk with no prefix",
			disk: Disk{
				Name:   "Disk 1",
				Format: "DVD",
				Path:   "/media/Some Movie (2020) [Film]/Disk [DVD]",
			},
			prefix:   "",
			expected: "vlc dvd:///media/Some Movie (2020) [Film]/Disk [DVD]",
		},
		{
			name: "Blu-Ray disk with prefix",
			disk: Disk{
				Name:   "Disk 1",
				Format: "Blu-Ray",
				Path:   "/media/War of the Worlds (2025) [Film]/Disk [Blu-Ray]",
			},
			prefix:   "/mnt/nas",
			expected: "vlc bluray:///mnt/nas/media/War of the Worlds (2025) [Film]/Disk [Blu-Ray]",
		},
		{
			name: "DVD disk with prefix",
			disk: Disk{
				Name:   "Series 1 Disk 1",
				Format: "DVD",
				Path:   "/media/Better Call Saul [TV]/Series 1 Disk 1 [DVD]",
			},
			prefix:   "/mnt/network",
			expected: "vlc dvd:///mnt/network/media/Better Call Saul [TV]/Series 1 Disk 1 [DVD]",
		},
		{
			name: "Unknown format defaults to file protocol",
			disk: Disk{
				Name:   "Disk 1",
				Format: "Unknown",
				Path:   "/media/Some Media/Disk [Unknown]",
			},
			prefix:   "",
			expected: "vlc file:///media/Some Media/Disk [Unknown]",
		},
		{
			name: "Case insensitive Blu-Ray matching",
			disk: Disk{
				Name:   "Disk 1",
				Format: "BLU-RAY",
				Path:   "/media/Movie/Disk [BLU-RAY]",
			},
			prefix:   "",
			expected: "vlc bluray:///media/Movie/Disk [BLU-RAY]",
		},
		{
			name: "Case insensitive DVD matching",
			disk: Disk{
				Name:   "Disk 1",
				Format: "dvd",
				Path:   "/media/Movie/Disk [dvd]",
			},
			prefix:   "",
			expected: "vlc dvd:///media/Movie/Disk [dvd]",
		},
		{
			name: "BluRay without hyphen",
			disk: Disk{
				Name:   "Disk 1",
				Format: "BluRay",
				Path:   "/media/Movie/Disk [BluRay]",
			},
			prefix:   "",
			expected: "vlc bluray:///media/Movie/Disk [BluRay]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.disk.PlayCommand(tt.prefix)
			if result != tt.expected {
				t.Errorf("Disk.PlayCommand() = %v, want %v", result, tt.expected)
			}
		})
	}
}
