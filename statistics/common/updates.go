package common

import (
	"encoding/json"
	"github.com/adarocket/informer/helpers"
	pb "github.com/adarocket/proto/proto-gen/common"
	"log"
	"time"
)

// GetUpdates -
func (commonStatistic *CommonStatistic) GetUpdates() *pb.Updates {
	var updates pb.Updates

	jsonBody, err := helpers.DownloadFile("",
		"https://api.github.com/repos/adarocket/informer/releases/latest")
	if err != nil {
		log.Println(err)
		return &updates
	}

	var info infoAboutLibGitHub
	err = json.Unmarshal(jsonBody, &info)
	if err != nil {
		log.Println(err)
		return &updates
	}

	updates.InformerActual = "v1.0.0"
	updates.InformerAvailable = info.TagName
	updates.UpdaterActual = ""
	updates.UpdaterAvailable = ""
	updates.PackagesAvailable = 0

	return &updates
}

type infoAboutLibGitHub struct {
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
}
