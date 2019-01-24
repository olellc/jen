/*
Jen-dir extracts audio from all video files in a directory without reencoding.
Under cover it uses FFmpeg for the actual work.
See
	jen-dir --help
for the command line options.
*/
package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type Opts struct {
	FFmpegDir string `long:"ffmpeg-dir" description:"Root directory of FFmpeg distribution" required:"true"`

	VideoDir string `short:"i" long:"videodir" description:"Input directory with videos to extract" required:"true"`

	OutDir string `short:"o" long:"outdir" description:"Output directory for audio. If exists, it will be removed before extraction." required:"true"`
}

func main() {

	var opts Opts

	var parser = flags.NewParser(&opts, flags.Default)

	_, err := parser.Parse()
	if err != nil {
		return
	}

	app := NewApp(opts.FFmpegDir, opts.VideoDir, opts.OutDir)

	err = app.Extract()
	if err != nil {
		fmt.Println(err)
		return
	}
}
