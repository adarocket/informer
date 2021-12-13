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
			f.shelleyGenesis = file
			if err != nil {
				log.Println(err)
			}

			file, err = helpers.DownloadFile(byronName, urlGenesisBase)
			f.byronGenesis = file
			if err != nil {
				log.Println(err)
			}

			file, err = helpers.DownloadFile(alonzoName, urlGenesisBase)
			f.alonzoGenesis = file
			if err != nil {
				log.Println(err)
			}

			time.Sleep(time.Hour * time.Duration(timeoutHours))
		}
	}()
}
