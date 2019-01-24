package platform

import (
	"io"
	"io/ioutil"
)

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
