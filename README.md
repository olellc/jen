# jen

Couple of applications for audio extraction from video files without re-encoding. Both applications use FFmpeg for the actual work.

## jen-dir

jen-dir recoursively extracts audio from all video files in a directory. Usage example:
```
$ jen-dir -i ~/video -o ~/audio
```

## jen-server

jen-server is an HTTP server providing audio extraction via browser UI. Usage example:
```
$ jen-server
```
By default it runs HTTP server on address http://localhost:8080.

## Note

Both applications use only `ffmpeg` and `ffprobe` commands from the FFmpeg distribution. These commands must be available on your `$PATH`. If they are not, then you can specify their location using `--ffmpeg` and `--ffprobe` command line options.
