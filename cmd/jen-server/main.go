package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
)

type Opts struct {
	FFmpegPath  string `long:"ffmpeg" description:"path to the ffmpeg command" default:"ffmpeg"`
	FFprobePath string `long:"ffprobe" description:"path to the ffprobe command" default:"ffprobe"`

	Addr string `long:"addr" description:"TCP network address to listen on" default:":8080"`
}

func main() {

	var opts Opts

	var parser = flags.NewParser(&opts, flags.Default)

	_, err := parser.Parse()
	if err != nil {
		return
	}

	err = run(opts.FFmpegPath, opts.FFprobePath, opts.Addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func run(ffmpegPath, ffprobePath, addr string) error {
	app, err := NewApp(ffmpegPath, ffprobePath)
	if err != nil {
		return err
	}

	err = ListenAndServeUntilSignal(addr, app.GetRouter())
	if err != http.ErrServerClosed {
		app.Close()
		return err
	}

	return app.Close()
}

/*
Works like http.ListenAndServe().
After receiving a signal performs server shutdown.
On successfull shutdown returns http.ErrServerClosed.

Example:

	http.HandleFunc("/", handler)
	err := ListenAndServeUntilSignal(":8080", nil)
	if err != http.ErrServerClosed {
		log.Println(err)
	}
*/
func ListenAndServeUntilSignal(addr string, handler http.Handler) error {
	sig_ch := make(chan os.Signal, 1)
	signal.Notify(sig_ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// see source for http.ListenAndServe()
	server := &http.Server{Addr: addr, Handler: handler}

	err_ch := make(chan error, 1)
	go func() {
		err_ch <- server.ListenAndServe()
	}()

	select {
	case err := <-err_ch:
		return err
	case <-sig_ch:
		err1 := server.Shutdown(context.Background())
		err2 := <-err_ch
		if err1 != nil {
			return err1
		}
		// if err2 == http.ErrServerClosed {
		// 	return nil
		// }
		return err2
	}
}
