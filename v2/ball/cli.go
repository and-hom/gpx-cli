package ball

import (
	"github.com/and-hom/gpx-cli/v2/model"
	"strconv"
	"log"
)

func ParseBallModelParams(params string) model.Model {
	model := BallModel{
		maxDist:15,
		minCount:15,
		dimension:2,
	}

	paramsFound := ballParamExpr.FindAllStringSubmatch(params, -1)
	for _, kv := range paramsFound {
		switch kv[1] {
		case "maxDist":
			maxDist, e := strconv.ParseInt(kv[2], 10, 64)
			if e != nil {
				log.Fatal(e.Error())
			}
			model.maxDist = maxDist
		case "minCount":
			minCount, e := strconv.ParseInt(kv[2], 10, 64)
			if e != nil {
				log.Fatal(e.Error())
			}
			model.minCount = minCount
		case "dimensions":
			dInt, e := strconv.ParseInt(kv[2], 10, 8)
			if e != nil {
				log.Fatal(e.Error())
			}
			dim := BallModelDimension(dInt)
			switch dim {
			case _2D:
				model.dimension = dim;
			case _3D:
				model.dimension = dim;
			default:
				log.Fatalf("Unknown dimension value %d. Can be 2 or 3", dim)
			}
		default:
			log.Fatalf("Unknown model param %s. Available are \"minCount\" and \"maxDist\" in meters and \"dimension\"=2 or 3", kv[1])
		}
	}
	return model
}

