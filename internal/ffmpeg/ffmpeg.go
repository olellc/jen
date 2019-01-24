// Package ffmpeg is a FFmpeg wrapper.
package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
)

// FFmpeg represents FFmpeg distribution
type FFmpeg struct {
	ffprobePath string
	ffmpegPath  string
}

func New(ffmpegDir string) *FFmpeg {
	var ff FFmpeg

	ff.ffprobePath = filepath.Join(ffmpegDir, "ffprobe")
	ff.ffmpegPath = filepath.Join(ffmpegDir, "ffmpeg")

	return &ff
}

// PathSwitcher returns a path to the file that would be used for
// audio extraction from the file at videoPath using audioFormat.
// The file at the audioPath must not exist.
type PathSwitcher func(videoPath string, audioFormat AudioFormat) (audioPath string, err error)

// Extract extracts audio from the file at videoPath without reencoding.
// It stores the extracted audio at the path obtained from the switcher.
func (ff *FFmpeg) Extract(videoPath string, switcher PathSwitcher) (
	audioPath string, audioFormat AudioFormat, err error) {

	audioFormat, err = ff.chooseFormat(videoPath)
	if err != nil {
		return "", AudioFormat{}, err
	}

	audioPath, err = switcher(videoPath, audioFormat)
	if err != nil {
		return "", AudioFormat{}, err
	}

	err = ff.extractToFormat(videoPath, audioPath, audioFormat.Name)
	if err != nil {
		return "", AudioFormat{}, err
	}

	return audioPath, audioFormat, nil
}

// chooseFormat chooses audio format by the content of the file at videoPath.
// The format is chosen so that audio extraction to the format
// needs no reencoding.
func (ff *FFmpeg) chooseFormat(videoPath string) (AudioFormat, error) {
	out, err := ff.ffprobe(videoPath)
	if err != nil {
		return AudioFormat{}, err
	}

	return chooseFormatRaw(out)
}

// ffprobe determines content of the path
func (ff *FFmpeg) ffprobe(path string) (output []byte, err error) {
	cmd := exec.Command(ff.ffprobePath,
		"-loglevel", "quiet", "-print_format", "json", "-show_streams", "-select_streams", "a", path)

	cmd.Env = []string{}

	return cmd.Output()
}

// Choose audio format by video file content
// data - ffprobe output
func chooseFormatRaw(data []byte) (AudioFormat, error) {
	codec_long_name, err := codecLongName(data)
	if err != nil {
		return AudioFormat{}, err
	}

	af, present := audioFormats[codec_long_name]
	if !present {
		return AudioFormat{}, fmt.Errorf("unknown codec_long_name: %q", codec_long_name)
	}

	return af, nil
}

// data - ffprobe output
func codecLongName(data []byte) (string, error) {
	var msg struct {
		Streams []struct {
			Codec_long_name string `json:"codec_long_name"`
		} `json:"streams"`
	}

	err := json.Unmarshal(data, &msg)
	if err != nil {
		return "", err
	}

	if streamCount := len(msg.Streams); streamCount != 1 {
		return "", fmt.Errorf("strange audio stream count: %v", streamCount)
	}

	return msg.Streams[0].Codec_long_name, nil
}

// Extract audio from videoPath to audioPath without reencoding.
// The file at the audioPath must not exist and must have an extension.
// ffmpeg will try to derive audio container format from the extension.
// This derived format must not require reencoding.
func (ff *FFmpeg) extractToExt(videoPath, audioPath string) error {
	arg := []string{
		"-nostdin", "-loglevel", "quiet",
		"-i", videoPath,
		"-vn", "-acodec", "copy",
		audioPath,
	}

	cmd := exec.Command(ff.ffmpegPath, arg...)
	cmd.Env = []string{}

	return cmd.Run()
}

// Extract audio from videoPath to audioPath without reencoding.
// The file at the audioPath must not exist.
// audioFormatName must not require reencoding.
// When running with the "-f" option defined, ffmpeg doesn't use audioPath
// extension (if any).
func (ff *FFmpeg) extractToFormat(videoPath, audioPath, audioFormatName string) error {
	arg := []string{
		"-nostdin", "-loglevel", "quiet",
		"-i", videoPath,
		"-vn", "-acodec", "copy",
		"-f", audioFormatName,
		audioPath,
	}

	cmd := exec.Command(ff.ffmpegPath, arg...)
	cmd.Env = []string{}

	return cmd.Run()
}
