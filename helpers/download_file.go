package helpers

import (
	"io"
	"log"
	"net/http"
)

func DownloadFile(fileName, urlBase string) ([]byte, error) {
	res, err := http.Get(urlBase + fileName)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Printf("Response failed with status code: %d and\nbody: %s\n",
			res.StatusCode, body)
		return []byte{}, err
	}

	return body, nil
}
