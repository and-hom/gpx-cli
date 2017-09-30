package v2

import (
	"github.com/ptrv/go-gpx"
	"sort"
	"github.com/and-hom/gpx-cli/util"
)

func doConcatenation(result *gpx.Gpx, mode Mode) {
	if (mode.preserveTracks) {
		if mode.preserveSegments {
			sort.Sort(util.ByTimestamp(result.Tracks))
		} else {
			for i := 0; i < len(result.Tracks); i++ {
				var wptLen = 0
				for j := 0; j < len(result.Tracks[i].Segments); j++ {
					wptLen += len(result.Tracks[i].Segments[j].Waypoints)
				}
				wpts := make([]gpx.Wpt, wptLen)
				var idx = 0
				for j := 0; j < len(result.Tracks[i].Segments); j++ {
					for k := 0; k < len(result.Tracks[i].Segments[j].Waypoints); k++ {
						wpts[idx] = result.Tracks[i].Segments[j].Waypoints[k]
						idx++
					}
				}
				result.Tracks[i].Segments = []gpx.Trkseg{{
					Waypoints:wpts,
				}}
			}
		}
	} else if mode.preserveSegments {
		var segLen = 0;
		for i := 0; i < len(result.Tracks); i++ {
			segLen += len(result.Tracks[i].Segments)
		}
		segments := make([]gpx.Trkseg, segLen)
		var idx = 0
		for i := 0; i < len(result.Tracks); i++ {
			for j := 0; j < len(result.Tracks[i].Segments); j++ {
				segments[idx] = result.Tracks[i].Segments[j]
				idx++
			}
		}
		sort.Sort(util.ByTimestampSegment(segments))
		result.Tracks = []gpx.Trk{{
			Name:"United track",
			Segments:segments,
		}}
	} else {
		var wptLen = 0
		for i := 0; i < len(result.Tracks); i++ {
			for j := 0; j < len(result.Tracks[i].Segments); j++ {
				wptLen += len(result.Tracks[i].Segments[j].Waypoints)
			}
		}
		wpts := make([]gpx.Wpt, wptLen)
		var idx = 0
		for i := 0; i < len(result.Tracks); i++ {
			for j := 0; j < len(result.Tracks[i].Segments); j++ {
				for k := 0; k < len(result.Tracks[i].Segments[j].Waypoints); k++ {
					wpts[idx] = result.Tracks[i].Segments[j].Waypoints[k]
					idx++
				}
			}
		}
		result.Tracks = []gpx.Trk{{
			Name:"United track",
			Segments:[]gpx.Trkseg{{
				Waypoints:wpts,
			}},
		}}
	}
}
