package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
	"sort"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

func main() {
	log.SetOutput(ioutil.Discard)
	for _, arg := range os.Args {
		if arg == "--verbose" {
			log.SetOutput(os.Stderr)
		}
	}

	app := cli.NewApp()
	app.Name = "Gpx Helper"
	app.Usage = "Simple command line interface to GPX"
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:    "concat",
			Aliases: []string{"c"},
			Usage:   "concatenate GPX tracks",
			Action: concat,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "order-by",
					Usage: "Sort order for points in new track - by input files order or by point timestamp",
				},
				cli.StringFlag{
					Name: "out",
					Usage: "Target file. Use '-' to print to stdout",
				},
				cli.BoolFlag{
					Name: "preserve-segments",
					Usage: "Preserve source track segments (by-default no)",
				},
			},
		},
		{
			Name:           "length",
			Aliases:        []string{"len", "l"},
			Usage:          "Calculate track length",
			Action:                trklen,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "2d",
					Usage: "Do not use altitude in calculations",
				},
			},
		},
		{
			Name:           "simplify",
			Usage:          "Remove track artifacts when navigator does not change position but writes track",
			Action:                trksimpl,
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:"min-points",
					Usage:"Minimum count of near-placed points to delete",
				},
				cli.UintFlag{
					Name: "max-dist",
					Usage:"Maximum distance between any points in cluster to remove",
				},
				cli.BoolFlag{
					Name:"interactive",
					Usage:"Ask for every cluster remove action [NOT IMPLEMENTED YET]",
				},
				cli.StringFlag{
					Name: "out",
					Usage: "Target file. Use '-' to print to stdout",
				},
			},
		},
		{
			Name:        "cut",
			Usage:       "Remove parts from track according to query",
			Action:      trkcut,
			Flags:       []cli.Flag{
				cli.StringFlag{
					Name: "query",
					Usage: "Target file. Use '-' to print to stdout",
				},
				cli.StringFlag{
					Name: "out",
					Usage: "Target file. Use '-' to print to stdout",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:"verbose",
			Usage:"Switch logging on",
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}