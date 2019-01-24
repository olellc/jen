package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/olellc/jen/ffmpeg"
)

const downloadPage = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Audio Extraction Service</title>
</head>
<body>
    <a href="%s">Download Audio</a>
</body>
</html>`

// POST /extractor
// Content-Type: multipart/form-data
// First part treated as a videofile.
func (app *App) extractor(w http.ResponseWriter, r *http.Request) error {
	multi_reader, err := r.MultipartReader()
	if err != nil {
		return err
	}

	part, err := multi_reader.NextPart()
	if err != nil {
		return err
	}

	id, audioPath, audioFormat, err := app.extractReader(part)
	if err != nil {
		return err
	}

	friendly_name := nameConv(part.FileName(), audioFormat.Ext)

	app.audios.Add(id, AudioHandle{
		Path:         audioPath,
		FriendlyName: friendly_name,
		MediaType:    audioFormat.MediaType})

	// link to the extracted audio
	url_path := "/audio/" + id + "/" + url.PathEscape(friendly_name)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, downloadPage, url_path)

	return nil
}

// videoFN - kinda video file name
// audioFN - kinda audio file name
func nameConv(videoFN, audioExt string) (audioFN string) {
	videoExt := filepath.Ext(videoFN)
	audioFN = videoFN[:len(videoFN)-len(videoExt)] + "." + audioExt

	return audioFN
}

// id must look like "951134894"
func (app *App) extractReader(reader io.Reader) (
	id, audioPath string, audioFormat ffmpeg.AudioFormat, err error) {

	videoPath, err := Reader2TempFile(reader, app.cluster.VideoDir)
	if err != nil {
		return "", "", ffmpeg.AudioFormat{}, err
	}

	// [not quite] using video file name as id
	_, id = filepath.Split(videoPath)

	audioPath, audioFormat, err = app.cmd.Extract(videoPath, app.getAudioPath)
	if err != nil {
		return "", "", ffmpeg.AudioFormat{}, err
	}

	return id, audioPath, audioFormat, nil
}

func (app *App) getAudioPath(videoPath string, af ffmpeg.AudioFormat) (string, error) {
	return getAudioPathImpl(videoPath, app.cluster.AudioDir), nil
}

func getAudioPathImpl(videoPath, audioDir string) string {
	// videoFN must look like "951134894"
	_, videoFN := filepath.Split(videoPath)

	return filepath.Join(audioDir, videoFN)
}

// Reader2TempFile copies data from reader to a new temporary file until
// either EOF is reached on reader or an error occurs.
// The temporary file is created in the directory dir.
// Reader2TempFile returns the path to the temporary file.
// A successful call returns err == nil, not err == io.EOF
func Reader2TempFile(reader io.Reader, dir string) (path string, err error) {
	file, err := ioutil.TempFile(dir, "")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		file.Close()
		return "", err
	}

	path = file.Name()

	err = file.Close()
	if err != nil {
		return "", err
	}

	return path, nil
}
