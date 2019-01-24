package ffmpeg

import (
	"fmt"
)

type AudioFormat struct {
	Name      string // audio format for the ffmpeg's -f option
	Ext       string // audio file extension suitable for the format
	MediaType string // media type suitable for serving a file stored in the format
}

func (af AudioFormat) String() string {
	return fmt.Sprintf("{%q, %q, %q}", af.Name, af.Ext, af.MediaType)
}

// codec_long_name to AudioFormat
var audioFormats = map[string]AudioFormat{
	"MP3 (MPEG audio layer 3)":    {"mp3", "mp3", "audio/mpeg"},
	"AAC (Advanced Audio Coding)": {"ipod", "m4a", "audio/mp4"},
	"Windows Media Audio 2":       {"asf", "wma", "audio/x-ms-wma"},
}
