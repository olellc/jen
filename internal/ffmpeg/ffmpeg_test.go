package ffmpeg

import (
	"io/ioutil"
	"testing"
)

func TestCodecLongName(t *testing.T) {
	filepath := "testdata/1.mp3.json"
	target_codec_long_name := "MP3 (MPEG audio layer 3)"

	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Errorf("cannot read file %q: %v", filepath, err)
		return
	}

	codec_long_name, err := codecLongName(b)
	if err != nil {
		t.Errorf("%q: codecLongName returns unexpected error: %v", filepath, err)
		return
	}

	if codec_long_name != target_codec_long_name {
		t.Errorf("%q: strange codec long name: %q, want: %q",
			filepath, codec_long_name, target_codec_long_name)
	}
}

func TestChooseFormatRaw(t *testing.T) {
	cases := []struct {
		filepath string
		af       AudioFormat
		ok       bool
	}{
		{"testdata/1.mp3.json", AudioFormat{"mp3", "mp3", "audio/mpeg"}, true},
		{"testdata/err-cln.json", AudioFormat{}, false},
	}

	for _, c := range cases {
		b, err := ioutil.ReadFile(c.filepath)
		if err != nil {
			t.Errorf("cannot read file %q: %v", c.filepath, err)
			continue
		}

		af, err := chooseFormatRaw(b)

		if !c.ok {
			if err == nil {
				t.Errorf("%q: unexpected nil error", c.filepath)
			}
			continue
		}

		if err != nil {
			t.Errorf("%q: unexpected error: %v", c.filepath, err)
			continue
		}

		if af != c.af {
			t.Errorf("%q: unexpected audio format: %q, want: %q", c.filepath, af, c.af)
		}
	}
}
