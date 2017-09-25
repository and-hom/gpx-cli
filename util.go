package main

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"github.com/ptrv/go-gpx"
	"io"
)

func withGpxFiles(paths []string, callback func(string, *gpx.Gpx)) {
	for _, name := range paths {
		f, err := os.Open(name)
		if err != nil {
			log.Warnf("Can not open file %s: %s\n", name, err.Error())
			continue
		}
		defer f.Close()
		gpx_data, err := gpx.Parse(f)
		if err != nil {
			log.Warnf("Can not read GPX data from %s: %s\n", name, err.Error())
			continue
		}
		callback(name, gpx_data)
	}
}

func modifyWaypointsBySegment(gpxData *gpx.Gpx, target io.Writer, modifier func(*[]gpx.Wpt) (*[]gpx.Wpt, bool)) {
	for t_i := 0; t_i < len(gpxData.Tracks); t_i++ {
		track := gpxData.Tracks[t_i]
		for s_i := 0; s_i < len(track.Segments); s_i++ {
			seg := track.Segments[s_i]
			segWpts := []gpx.Wpt(seg.Waypoints)
			wpts, changed := modifier(&segWpts)

			if changed {
				seg.Waypoints = *wpts
				track.Segments[s_i] = seg
				gpxData.Tracks[t_i] = track
			}
		}
	}

	xmlBytes := gpxData.ToXML()
	target.Write(xmlBytes)
}