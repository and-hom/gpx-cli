package main

import (
	"gopkg.in/urfave/cli.v1"
	log "github.com/Sirupsen/logrus"
	"github.com/ptrv/go-gpx"
	"fmt"
	"github.com/and-hom/gpx-cli/util"
)

func printRow(rType string, id interface{}, l float64) {
	fmt.Printf("%s\t%v\t%.3f\n", rType, id, l)
}

func trklen(c *cli.Context) error {
	log.Warn("Calculate track length without height")
	if (len(c.Args()) == 0) {
		log.Warn("No input files - exiting")
		return nil
	}
	use2d := c.Bool("2d")
	util.WithGpxFiles(c.Args(), func(fileName string, gpxData *gpx.Gpx) {
		var sumLen = float64(0)
		for _, track := range gpxData.Tracks {
			for sIdx, seg := range track.Segments {
				sLen := seg.Length3D() / 1000
				printRow("segment", sIdx, sLen)
			}
			var tLen float64
			if use2d {
				tLen = track.Length2D() / 1000
			} else {
				tLen = track.Length3D() / 1000
			}
			printRow("track", track.Name, tLen)
			sumLen += tLen
		}
		printRow("file", fileName, sumLen)
	})
	return nil
}