package main

import (
	"fmt"
	"testing"
)

func TestGetAudioPathImpl(t *testing.T) {
	audioDir := "/tmp/jen-418706227/audio"

	cases := []struct {
		videoPath string
		audioPath string
	}{
		{"/tmp/jen-418706227/video/951134894", "/tmp/jen-418706227/audio/951134894"},
	}

	printer := func(audioPath, videoPath, audioDir string) string {
		return fmt.Sprintf("%q := getAudioPathImpl(%q, %q)",
			audioPath, videoPath, audioDir)
	}

	for _, c := range cases {
		audioPath := getAudioPathImpl(c.videoPath, audioDir)
		if audioPath != c.audioPath {
			t.Errorf("wrong audioPath: %s, want %q",
				printer(audioPath, c.videoPath, audioDir),
				c.audioPath)
			continue
		}
	}
}
