package util

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"github.com/ptrv/go-gpx"
	"io"
	"errors"
)

type ByTimestamp []gpx.Trk

func (a ByTimestamp) Len() int {
	return len(a)
}
func (a ByTimestamp) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByTimestamp) Less(i, j int) bool {
	t1, _ := a[i].TimeBounds()
	t2, _ := a[j].TimeBounds()
	return t1.Before(t2)
}

type ByTimestampSegment []gpx.Trkseg

func (a ByTimestampSegment) Len() int {
	return len(a)
}
func (a ByTimestampSegment) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByTimestampSegment) Less(i, j int) bool {
	t1, _ := a[i].Waypoints.TimeBounds()
	t2, _ := a[j].Waypoints.TimeBounds()
	return t1.Before(t2)
}

func WithGpxFiles(paths []string, callback func(string, *gpx.Gpx)) {
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

func ModifyWaypointsBySegment(gpxData *gpx.Gpx, target io.Writer, modifier func(*[]gpx.Wpt) (*[]gpx.Wpt, bool)) {
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

func GetTarget(targetName string) (io.WriteCloser, error) {
	if targetName == "-" {
		return os.Stdout, nil
	} else {
		// detect if file exists
		var _, err = os.Stat(targetName)
		if os.IsNotExist(err) {
			target, err := os.Create(targetName)
			if err != nil {
				log.Errorf("Can not open target file %s: %s\n", targetName, err.Error())
				return nil, err
			}
			return target, nil
		} else {
			log.Errorf("File %s already exists", targetName)
			return nil, errors.New("File exists: " + targetName)
		}
	}
}
