package main

import (
	"gopkg.in/urfave/cli.v1"
	"sort"
	"github.com/ptrv/go-gpx"
	log "github.com/Sirupsen/logrus"
	"io"
	"encoding/xml"
	"errors"
	"github.com/and-hom/gpx-cli/util"
)

type ByName []gpx.Trk

func (a ByName) Len() int {
	return len(a)
}
func (a ByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByName) Less(i, j int) bool {
	n1 := a[i].Name
	n2 := a[j].Name
	return n1 < n2
}

func writeXml(w io.Writer, data interface{}) error {
	content, err := xml.MarshalIndent(data, "", "	")
	if (err != nil) {
		return err
	}
	_, err = w.Write(content)
	return err
}

func concat(c *cli.Context) error {
	_, err := getOrder(c)
	if (err != nil) {
		return err
	}

	if (len(c.Args()) == 0) {
		log.Warn("No input files - exiting")
		return nil
	}

	target, err := util.GetTarget(c.String("out"))
	if (err != nil) {
		return err
	}
	defer target.Close()

	preserveSegments := c.Bool("preserve-segments")

	target.Write([]byte(xml.Header))
	target.Write([]byte("<gpx><trk><name>Concatenated</name>"))
	if !preserveSegments {
		target.Write([]byte("<trkseg>"))
	}
	util.WithGpxFiles(c.Args(), func(_ string, gpxData *gpx.Gpx) {
		if len(gpxData.Waypoints) > 0 {
			log.Warn("Some waypoints detected - will not be copied to target file")
		}
		if len(gpxData.Routes) > 0 {
			log.Warn("Some routes detected - can not be processed and will not be copied to target file")
		}

		tracks := gpxData.Tracks
		//sort.Sort(ByDate(tracks))
		sort.Sort(ByName(tracks))
		for _, trk := range tracks {
			log.Infof("Importing track %s", trk.Name)
			for _, seg := range trk.Segments {
				if (preserveSegments) {
					writeXml(target, &seg)
				} else {
					writeXml(target, trkpts(seg.Waypoints))
				}
			}
		}
	})

	if !preserveSegments {
		target.Write([]byte("</trkseg>"))
	}
	target.Write([]byte("</trk></gpx>"))
	return nil
}

func trkpts(wpts []gpx.Wpt) interface{} {
	result := make([]Trkpt, len(wpts))
	for i, wpt := range wpts {
		result[i] = Trkpt{wpt, struct {

		}{}}
	}
	return result
}

type Trkpt struct {
	gpx.Wpt
	XMLName struct{}    `xml:"trkpt"`
}

func getOrder(c *cli.Context) (string, error) {
	orderBy := c.String("order-by")
	if (orderBy == "") {
		return "files", nil
	}
	if (orderBy != "files") {
		log.Error("Only 'files' order supported now")
		return "", errors.New("Only 'files' order supported now")
	}
	return orderBy, nil
}