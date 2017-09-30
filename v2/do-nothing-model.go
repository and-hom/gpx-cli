package v2

import (
	"github.com/ptrv/go-gpx"
	"github.com/and-hom/gpx-cli/v2/model"
)

type DoNothingModel struct {

}

func (m DoNothingModel)Features([]gpx.Wpt) []model.Feature {
	return []model.Feature{}
}

func (m DoNothingModel)Cut(wpts []gpx.Wpt, _ []model.Feature) []gpx.Wpt {
	return wpts
}

