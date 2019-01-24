package main

import (
	"fmt"
	"net/http"
)

const main_page = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Audio Extraction Service</title>
</head>
<body>
    <form enctype="multipart/form-data" action="/extractor" method="POST">
        <input type="file" name="videoFile" />
        <input type="submit" value="Upload" />
    </form>
</body>
</html>`

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, main_page)
}
