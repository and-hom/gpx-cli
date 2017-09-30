package v2

import (
	"gopkg.in/urfave/cli.v1"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"github.com/ptrv/go-gpx"
	"os"
	"time"
	"github.com/and-hom/gpx-cli/util"
	"github.com/and-hom/gpx-cli/v2/ball"
	"github.com/and-hom/gpx-cli/v2/model"
)

func parseDoNothingModelParams(_ string) model.Model {
	return DoNothingModel{}
}

func Filter(c *cli.Context) error {
	models := make(map[string]func(string) model.Model)
	models["none"] = parseDoNothingModelParams
	models["ball"] = ball.ParseBallModelParams

	mode, e := ParseMode(c.String(FILTER_MODE_FLAG))
	if e != nil {
		log.Fatal(e.Error())
		return e
	}

	modelId := c.String(FILTER_MODEL_FLAG)
	mType, ok := models[modelId]
	if !ok {
		e := fmt.Errorf("Unknown mode %s", modelId)
		log.Fatal(e.Error())
		return e
	}
	currentModel := mType(c.String(FILTER_MODEL_PARAMS_FLAG))

	// Load
	result := gpx.NewGpx()
	util.WithGpxFiles(c.Args(), func(s string, g *gpx.Gpx) {
		result.Tracks = append(result.Tracks, g.Tracks...)
	})
	// Sort
	doConcatenation(result, mode)

	// process data
	for i := 0; i < len(result.Tracks); i++ {
		for j := 0; j < len(result.Tracks[i].Segments); j++ {
			features := currentModel.Features(result.Tracks[i].Segments[j].Waypoints)

			featuresDumpFile, err := os.Create(fmt.Sprintf("%d", time.Now().Unix()))
			if err != nil {
				log.Fatal(err.Error())
			}
			defer featuresDumpFile.Close()
			for k := 0; k < len(features); k++ {
				featuresDumpFile.WriteString(features[k].CsvLine())
			}

			result.Tracks[i].Segments[j].Waypoints = currentModel.Cut(result.Tracks[i].Segments[j].Waypoints, features)
		}
	}

	xmlBytes := result.ToXML()
	_, e = os.Stdout.Write(xmlBytes)
	if e != nil {
		log.Fatal(e.Error())
	}

	return e
}
