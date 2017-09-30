package main

import (
	"gopkg.in/urfave/cli.v1"
	log "github.com/Sirupsen/logrus"
	"github.com/ptrv/go-gpx"
	"errors"
	"github.com/and-hom/gpx-cli/util"
)

func trksimpl(c *cli.Context) error {
	if (len(c.Args()) == 0) {
		log.Warn("No input files - exiting")
		return errors.New("Input is missing")
	}
	if (len(c.Args()) > 1) {
		log.Warn("Can not process more then one file")
		return errors.New("Too much files")
	}

	target, err := util.GetTarget(c.String("out"))
	if (err != nil) {
		return err
	}
	defer target.Close()

	minPoints := int(c.Uint("min-points"))
	if minPoints == 0 {
		minPoints = 20
	}
	maxDist := c.Int("max-dist")
	if maxDist == 0 {
		maxDist = 20
	}
	log.Infof("Removing clusters greater then %d points with distance not more then %d meters\n", minPoints, maxDist)

	util.WithGpxFiles(c.Args(), func(fileName string, gpxData *gpx.Gpx) {
		util.ModifyWaypointsBySegment(gpxData, target, func(wpts *[]gpx.Wpt) (*[]gpx.Wpt, bool) {
			size := len(*wpts)
			if (size == 0) {
				return wpts, false
			}

			resultPts := make([]gpx.Wpt, size)
			var count = 1
			var clusterSize = 0;
			var changed = false

			resultPts[0] = (*wpts)[0]
			for i := 1; i < size; i++ {
				jLim := size - i
				clusterSize = 0
				for j := 1; j < jLim; j++ {
					lastPoint := (j == jLim - 1)
					if (*wpts)[i].Distance2D(&((*wpts)[i + j])) < float64(maxDist) && !lastPoint {
						clusterSize++
					} else {
						if (lastPoint) {
							clusterSize++
						}
						lastIdx := (i + j - 2) // -2 to keep last point of the cluster
						if clusterSize >= minPoints {
							log.Infof("Drop points from [%f %f] %s to %s: %d\n",
								(*wpts)[i].Lat,
								(*wpts)[i].Lon,
								(*wpts)[i].Timestamp,
								(*wpts)[lastIdx].Timestamp,
								clusterSize)
							i = lastIdx
							changed = true
						}
						break
					}
				}
				resultPts[count] = (*wpts)[i]
				count++
			}

			finalPts := make([]gpx.Wpt, count)
			for i := 0; i < count; i++ {
				finalPts[i] = resultPts[i]
			}

			return &finalPts, changed

		})
	})
	return nil
}