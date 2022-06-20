package cardano

import (
	"log"
	"time"

	"github.com/adarocket/informer/helpers"
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
	fDownload := func() {
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
	}

	fDownload()
	go func() {
		// c := time.Tick(time.Duration(timeoutHours))
		// for range c {
		// 	fDownload()
		// }

		for {
			time.Sleep(time.Duration(timeoutHours) * time.Hour)
			fDownload()
		}
	}()
}
