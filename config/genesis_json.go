package config

import (
	"io"
	"log"
	"net/http"
)

var (
	urlDownloadBase = "https://hydra.iohk.io/job/Cardano/cardano-node/cardano-deployment/" +
		"latest-finished/download/1/"
	shelleyName = "mainnet-shelley-genesis.json"
	byronName   = "mainnet-byron-genesis.json"
)

func getGenesisJson(jsonName string) ([]byte, error) {
	res, err := http.Get(urlDownloadBase + jsonName)
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
