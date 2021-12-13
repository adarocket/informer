package cardano

import (
	"github.com/adarocket/informer/helpers"
	"log"
	"time"
)

const (
	urlGenesisBase = "https://hydra.iohk.io/job/Cardano/cardano-node/cardano-deployment/" +
		"latest-finished/download/1/"
	shelleyName = "mainnet-shelley-genesis.json"
	byronName   = "mainnet-byron-genesis.json"
	alonzoName  = "mainnet-alonzo-genesis.json"
)

type genesisJsonFiles struct {
	shelleyGenesis []byte
	alonzoGenesis  []byte
	byronGenesis   []byte
}

// startUpdateTimeoutCycle - open new goroutine and update struct fields
func (f *genesisJsonFiles) startUpdateTimeoutCycle(timeoutHours int) {
	go func() {
		for {
			file, err := helpers.DownloadFile(shelleyName, urlGenesisBase)
			if err != nil {
				log.Println(err)
			} else {
				f.shelleyGenesis = file
			}

			file, err = helpers.DownloadFile(byronName, urlGenesisBase)
			if err != nil {
				log.Println(err)
			} else {
				f.byronGenesis = file
			}

			file, err = helpers.DownloadFile(alonzoName, urlGenesisBase)
			if err != nil {
				log.Println(err)
			} else {
				f.alonzoGenesis = file
			}

			time.Sleep(time.Hour * time.Duration(timeoutHours))
		}
	}()
}
