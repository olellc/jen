package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/olellc/jen/internal/ffmpeg"
)

type App struct {
	cmd       *ffmpeg.FFmpeg
	videoRoot string // root directory with videos
	outRoot   string // root directory for extracted sound
}

func NewApp(ffmpegDir, videoDir, outDir string) *App {
	return &App{
		cmd:       ffmpeg.New(ffmpegDir),
		videoRoot: videoDir,
		outRoot:   outDir,
	}
}

func (app *App) Extract() error {
	err := os.RemoveAll(app.outRoot)
	if err != nil {
		return err
	}

	err = os.MkdirAll(app.outRoot, os.ModePerm)
	if err != nil {
		return err
	}

	return filepath.Walk(app.videoRoot,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return app.extractDir(path)
			}
			if info.Mode().IsRegular() {
				err = app.extractFile(path)
				if err != nil {
					fmt.Printf("%q: %v\n", path, err)
				}
				return nil
			}
			return nil
		})
}

// path must be rooted at app.videoRoot
func (app *App) extractDir(path string) error {
	newpath, err := app.switchRoot(path)
	if err != nil {
		return err
	}
	if newpath == app.outRoot {
		return nil
	}

	return os.Mkdir(newpath, os.ModePerm)
}

// Switch path root from app.videoRoot to app.outRoot.
// path must be rooted at app.videoRoot
func (app *App) switchRoot(path string) (string, error) {
	return switchRootImpl(app.videoRoot, app.outRoot, path)
}

// Switch path root from videoRoot to outRoot.
// path must be rooted at videoRoot
func switchRootImpl(videoRoot, outRoot, path string) (string, error) {
	rel, err := filepath.Rel(videoRoot, path)
	if err != nil {
		return "", err
	}

	return filepath.Join(outRoot, rel), nil
}

// Filter files by extension
var extFilter = map[string]struct{}{
	"avi": {},
	"flv": {},
	"wmv": {},
}

// path must be rooted at app.videoRoot
func (app *App) extractFile(path string) error {
	videoExt := filepath.Ext(path)
	if len(videoExt) == 0 {
		return errors.New("no extension: scipping file")
	}

	videoExt = videoExt[1:] // scipping dot
	if _, present := extFilter[videoExt]; !present {
		return errors.New("unknown extension: scipping file")
	}

	_, _, err := app.cmd.Extract(path, app.getAudioPath)

	return err
}

// videoPath must be a file rooted at app.videoRoot
func (app *App) getAudioPath(videoPath string, af ffmpeg.AudioFormat) (string, error) {
	return getAudioPathImpl(app.videoRoot, app.outRoot, videoPath, af.Ext)
}

// videoPath must be a file rooted at videoRoot
func getAudioPathImpl(videoRoot, outRoot, videoPath, audioExt string) (string, error) {
	videoExt := filepath.Ext(videoPath)
	audioPath := videoPath[:len(videoPath)-len(videoExt)] + "." + audioExt

	audioPath, err := switchRootImpl(videoRoot, outRoot, audioPath)
	if err != nil {
		return "", err
	}

	return audioPath, nil
}
