package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/leoadamek/irapi"
)

func main() {

	ctx := context.Background()

	api := irapi.New(ctx, irapi.EnvironmentCredentialsProvider)

	if err := api.Login(ctx); err != nil {
		log.Fatalln("Failed to log into iRacing", err)
	}

	log.Println("Logged into iRacing")

	log.Println("Getting List of Seasons")

	seasons, err := api.GetSeasons(ctx, false)

	if err != nil {
		log.Fatalln("Unable to get seasons info", err)
	}

	log.Printf("Got details for %d seasons", len(seasons))

	tw := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)

	tw.Write([]byte("Year\tQ\tWk\tID\tLIC\tCategory\tSeries\tName\tStart\tEnd\n"))

	sort.Sort(seasons)

	for _, s := range seasons {
		line := fmt.Sprintf(
			"%04d\t%01d\t%02d\t%d\t%s\t%s\t%d\t%s\t%s\t%s\n",
			s.Year,
			s.Quarter,
			s.Week,
			s.SeasonID,
			s.LicenceGroup,
			s.Category,
			s.SeriesID,
			s.ShortName,
			time.Time(s.Start).Format("_2 Jan 2006"),
			time.Time(s.End).Format("_2 Jan 2006"),
		)

		tw.Write([]byte(line))
	}

	tw.Flush()
}
