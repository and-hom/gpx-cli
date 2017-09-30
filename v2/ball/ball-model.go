package ball

import (
	"regexp"
	"fmt"
	"github.com/ptrv/go-gpx"
	"container/list"
	"github.com/and-hom/gpx-cli/v2/model"
)

type BallModelDimension int8

const _2D BallModelDimension = 2
const _3D BallModelDimension = 3

var ballParamExpr = regexp.MustCompile("(\\w+)=(\\d+)")

type BallFeature struct {
	size      int
	threshold bool
	timestamp string
	alt       float64
	lDist     []float64
	rDist     []float64
}

func (b BallFeature) CsvLine() string {
	return fmt.Sprintf("%d\t%t\t%s\t%f\t%.0f\t%.0f\n", b.size, b.threshold, b.timestamp, b.alt, b.lDist, b.rDist)
}

type BallModel struct {
	maxDist   int64
	minCount  int64
	dimension BallModelDimension
}

func (m BallModel) distL(wpts []gpx.Wpt, i int, ballSize int) (float64, bool) {
	pointIdx := i - ballSize
	if pointIdx < 0 {
		return 0, false
	} else {
		return m.dist(wpts, i, pointIdx), true
	}
}

func (m BallModel) distR(wpts []gpx.Wpt, i int, ballSize int) (float64, bool) {
	pointIdx := i + ballSize
	if pointIdx >= len(wpts) {
		return 0, false
	} else {
		return m.dist(wpts, i, pointIdx), true
	}
}

func (m BallModel) dist(wpts []gpx.Wpt, i int, j int) float64 {
	if m.dimension == _2D {
		return wpts[i].Distance2D(&(wpts[j]))
	} else {
		return wpts[i].Distance3D(&(wpts[j]))
	}
}

func list2SliceFLoat64(lst *list.List) []float64 {
	s := make([]float64, lst.Len())
	var i = 0;
	for e := lst.Front(); e != nil; e = e.Next() {
		s[i] = e.Value.(float64)
		i++
	}
	return s
}

func (m BallModel)Features(wpts []gpx.Wpt) []model.Feature {
	features := make([]model.Feature, len(wpts))
	for i := 0; i < len(wpts); i++ {
		distsL := list.New()
		distsR := list.New()
		ballSize := 1
		for ; ballSize < len(wpts); ballSize++ {
			lDist, lBoundOk := m.distL(wpts, i, ballSize)
			rDist, rBoundOk := m.distR(wpts, i, ballSize)
			if lDist > float64(m.maxDist) || rDist > float64(m.maxDist) {
				break
			}
			if lBoundOk {
				distsL.PushBack(lDist)
			}
			if rBoundOk {
				distsR.PushBack(rDist)
			}
		}
		features[i] = BallFeature{ballSize, ballSize >= int(m.minCount), wpts[i].Timestamp, wpts[i].Ele, list2SliceFLoat64(distsL), list2SliceFLoat64(distsR), }
	}
	return features
}

func (m BallModel)Cut(wpts []gpx.Wpt, _ []model.Feature) []gpx.Wpt {
	return wpts
}
