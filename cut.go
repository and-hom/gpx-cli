package main

import (
	"gopkg.in/urfave/cli.v1"
	log "github.com/Sirupsen/logrus"
	"github.com/ptrv/go-gpx"
	"errors"
	"time"
	"regexp"
	"container/list"
	"os"
	"github.com/and-hom/gpx-cli/util"
)

// todo: make free args order
var queryPattern = regexp.MustCompile("[[:space:]]*(from=([0-9:\\-TZ]+))?[[:space:]]*(to=([0-9:\\-TZ]+))?[[:space:]]*(to-file=([^;]+))?[[:space:]]*(;|$)")

type  Interval struct {
	From   time.Time
	To     time.Time
	ToFile string
}

func (i Interval) contains(t time.Time) bool {
	return (i.From.IsZero() || !t.Before(i.From)) && (i.To.IsZero() || !t.After(i.To)) && !(i.isEmpty())
}

func (i Interval) isEmpty() bool {
	return i.From.IsZero() && i.To.IsZero()
}

type Query []Interval

func (q Query) getSubquery(wpts *[]gpx.Wpt) Query {
	var intervals = []Interval{}

	for i := 0; i < len(q); i++ {
		if !q[i].isEmpty() && (q[i].To.IsZero() || !q[i].To.Before((*wpts)[0].Time())) && (q[i].From.IsZero() || !q[i].From.After((*wpts)[len(*wpts) - 1].Time())) {
			intervals = append(intervals, q[i])
		}
	}

	return intervals
}

func parseQ(s string) (Query, error) {
	log.Infof("Query is %v", s)
	found := queryPattern.FindAllStringSubmatch(s, -1)
	if found == nil {
		return []Interval{}, nil
	}
	q := make([]Interval, len(found))
	log.Infof("Found is %v", found)
	for i, result := range found {
		var from = time.Time{}
		if result[1] != "" {
			fromStr := result[2]
			f, err := time.Parse(time.RFC3339, fromStr)
			if err != nil {
				return nil, err
			}
			from = f
		}

		var to = time.Time{}
		if result[3] != "" {
			toStr := result[4]
			f, err := time.Parse(time.RFC3339, toStr)
			if err != nil {
				return nil, err
			}
			to = f
		}

		var toFile = ""
		if result[5] != "" {
			toFile = result[6]
		}

		q[i] = Interval{
			From: from,
			To: to,
			ToFile:toFile,
		}
	}
	return q, nil
}

func trkcut(c *cli.Context) error {
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

	query, err := parseQ(c.String("query"))
	if err != nil {
		log.Errorf("Can not parse query %s: %v", query, err)
		return errors.New("Can not parse query")
	}

	log.Infof("Query is %v", query)

	util.WithGpxFiles(c.Args(), func(fileName string, gpxData *gpx.Gpx) {
		util.ModifyWaypointsBySegment(gpxData, target, func(wpts *[]gpx.Wpt) (*[]gpx.Wpt, bool) {
			sq := query.getSubquery(wpts)

			size := len(*wpts)

			resultPts := make([]gpx.Wpt, size)
			var count = 0
			var skip = false

			dump_files := make(map[Interval]*list.List)

			for i := 0; i < size; i++ {
				skip = false
				for j := 0; j < len(sq); j++ {
					if sq[j].contains((*wpts)[i].Time()) {
						skip = true
						//log.Infof("Cut %s", (*wpts)[i].Timestamp)
						if sq[j].ToFile != "" {
							lst, found := dump_files[sq[j]]
							if !found {
								dump_files[sq[j]] = list.New()
								dump_files[sq[j]].PushBack((*wpts)[i])
							} else {
								lst.PushBack((*wpts)[i])
							}

						}
					}
				}
				if !skip {
					resultPts[count] = (*wpts)[i]
					count++
				}
			}

			err := writeCutFiles(dump_files)
			if err!=nil {
				log.Error("Can not write files:", err.Error())
			}

			if size == count {
				return wpts, false
			}

			finalPts := make([]gpx.Wpt, count)
			for i := 0; i < count; i++ {
				finalPts[i] = resultPts[i]
			}
			return &finalPts, true

		})

	})

	return nil
}
func writeCutFiles(subTracks map[Interval]*list.List) error {
	for i := range subTracks {
		f,err := os.Create(i.ToFile)
		if err!=nil {
			return err
		}
		defer f.Close()

		log.Infof("Dump %v-%v to file %s", i.From, i.To, i.ToFile)

		points := subTracks[i]
		wpts := make([]gpx.Wpt, points.Len())
		var j=0
		for element := points.Front(); element != nil; element = element.Next() {

			wpts[j] = element.Value.(gpx.Wpt)
			j++
		}

		seg:=gpx.Trkseg{
			Waypoints: wpts,
		}

		trk := gpx.Trk{
			Name:i.ToFile,
			Segments: []gpx.Trkseg{seg},
		}

		gpxData := gpx.NewGpx()
		gpxData.Tracks = []gpx.Trk{trk}

		bytes := gpxData.ToXML()
		f.Write(bytes)
	}
	return nil
}