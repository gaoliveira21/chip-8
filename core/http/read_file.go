package http

import (
	"errors"
	"io"
	"net/http"
)

func ReadFile(path string) ([]byte, error) {
	response, err := http.Get(path)
	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []byte{}, errors.New("InvalidHttpStatusCode")
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
