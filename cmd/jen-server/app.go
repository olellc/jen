package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/olellc/jen/ffmpeg"
)

// App is an audio extraction service
type App struct {
	cmd     *ffmpeg.FFmpeg
	cluster *TempCluster
	audios  *AudioRegistry
}

func NewApp(ffmpegDir string) (*App, error) {
	cmd := ffmpeg.New(ffmpegDir)

	cluster, err := NewTempCluster()
	if err != nil {
		return nil, err
	}

	audios := NewAudioRegistry()

	app := App{
		cmd:     cmd,
		cluster: cluster,
		audios:  audios,
	}

	return &app, nil
}

func (app *App) Close() error {
	return app.cluster.Remove()
}

func (app *App) GetRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", mainPage)

	r.
		With(middleware.AllowContentType("multipart/form-data")).
		Post("/extractor", func(w http.ResponseWriter, r *http.Request) {
			err := app.extractor(w, r)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			}
		})

	r.Get("/audio/{id}/{friendly_name}", func(w http.ResponseWriter, r *http.Request) {
		err := app.audio(w, r)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}
	})

	return r
}

// TempCluster represents temp directory with a couple of subdirectories
type TempCluster struct {
	rootDir  string
	VideoDir string
	AudioDir string
}

func NewTempCluster() (*TempCluster, error) {
	var cluster TempCluster
	var err error

	cluster.rootDir, err = ioutil.TempDir("", "jen-")
	if err != nil {
		return nil, err
	}

	cluster.VideoDir = filepath.Join(cluster.rootDir, "video")
	err = os.Mkdir(cluster.VideoDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	cluster.AudioDir = filepath.Join(cluster.rootDir, "audio")
	err = os.Mkdir(cluster.AudioDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func (cluster *TempCluster) Remove() error {
	return os.RemoveAll(cluster.rootDir)
}

type AudioRegistry struct {
	sync.Mutex
	m map[string]AudioHandle // request_id to AudioHandle
}

func NewAudioRegistry() *AudioRegistry {
	return &AudioRegistry{m: map[string]AudioHandle{}}
}

func (audios *AudioRegistry) Add(id string, h AudioHandle) {
	audios.Lock()
	audios.m[id] = h
	audios.Unlock()
}

func (audios *AudioRegistry) Get(id string) (h AudioHandle, present bool) {
	audios.Lock()
	h, present = audios.m[id]
	audios.Unlock()

	return h, present
}

type AudioHandle struct {
	Path         string
	FriendlyName string // audio file name to be presented to a client
	MediaType    string
}
