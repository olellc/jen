package main

import (
	"fmt"
	"testing"
)

func TestSwitchRootImpl(t *testing.T) {
	videoRoot := "/video"
	outRoot := "/sound"

	cases := []struct {
		path    string
		newPath string
	}{
		{"/video", "/sound"},
		{"/video/aaa", "/sound/aaa"},
		{"/video/aaa/bbb", "/sound/aaa/bbb"},
	}

	printer := func(newPath string, err error, videoRoot, outRoot, path string) string {
		return fmt.Sprintf("%q, <%v> := switchRootImpl(%q, %q, %q)",
			newPath, err, videoRoot, outRoot, path)
	}

	for _, c := range cases {
		newPath, err := switchRootImpl(videoRoot, outRoot, c.path)
		if err != nil {
			t.Error("unexpected error: " + printer(newPath, err, videoRoot, outRoot, c.path))
			continue
		}
		if newPath != c.newPath {
			t.Errorf("wrong newPath: %s, want %q",
				printer(newPath, err, videoRoot, outRoot, c.path),
				c.newPath)
			continue
		}
	}
}

func TestGetAudioPathImpl(t *testing.T) {
	videoRoot := "/video"
	outRoot := "/sound"

	cases := []struct {
		videoPath string
		audioExt  string
		audioPath string
	}{
		{"/video/aaa.flv", "m4a", "/sound/aaa.m4a"},
	}

	printer := func(audioPath string, err error, videoRoot, outRoot, videoPath, audioExt string) string {
		return fmt.Sprintf("%q, <%v> := getAudioPathImpl(%q, %q, %q, %q)",
			audioPath, err, videoRoot, outRoot, videoPath, audioExt)
	}

	for _, c := range cases {
		audioPath, err := getAudioPathImpl(videoRoot, outRoot, c.videoPath, c.audioExt)
		if err != nil {
			t.Error("unexpected error: " + printer(audioPath, err,
				videoRoot, outRoot, c.videoPath, c.audioExt))
			continue
		}
		if audioPath != c.audioPath {
			t.Errorf("wrong audioPath: %s, want %q",
				printer(audioPath, err, videoRoot, outRoot, c.videoPath, c.audioExt),
				c.audioPath)
			continue
		}
	}
}
