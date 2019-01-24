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

	"github.com/olellc/jen/internal/ffmpeg"
)

type Opts struct {
	FFmpegPath  string `long:"ffmpeg" description:"path to the ffmpeg command" default:"ffmpeg"`
	FFprobePath string `long:"ffprobe" description:"path to the ffprobe command" default:"ffprobe"`

	VideoDir string `short:"i" long:"videodir" description:"input directory with videos to extract" required:"true"`

	OutDir string `short:"o" long:"outdir" description:"output directory for audio. If exists, it will be removed before extraction." required:"true"`
}

func main() {

	var opts Opts

	var parser = flags.NewParser(&opts, flags.Default)

	_, err := parser.Parse()
	if err != nil {
		return
	}

	app := App{
		cmd: &ffmpeg.FFmpeg{
			FFmpegPath:  opts.FFmpegPath,
			FFprobePath: opts.FFprobePath,
		},
		videoRoot: opts.VideoDir,
		outRoot:   opts.OutDir,
	}

	err = app.Extract()
	if err != nil {
		fmt.Println(err)
		return
	}
}
