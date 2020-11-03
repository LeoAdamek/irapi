package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/leoadamek/irapi"
)

func main() {
	ctx := context.Background()

	uid := flag.Uint64("u", 0, "User ID")
	flag.Parse()

	if uid == nil || *uid == 0 {
		flag.Usage()
		log.Fatalln("No User ID Provided")
	}

	api := irapi.New(ctx, irapi.EnvironmentCredentialsProvider)

	log.Println("Created API")

	if err := api.Login(ctx); err != nil {
		log.Fatalf("Failed to log into iRacing: %s", err)
	}

	params := irapi.DefaultSearchResultsOptions()
	params.Season = nil

	params.DateRange = &irapi.DateRange{
		Upper: time.Now(),
		Lower: time.Now().Add(-24 * 14 * time.Hour),
	}

	log.Println("Searching for the last 100 results from this week...")
	results, err := api.SearchResults(ctx, *uid, params)

	if err != nil {
		log.Fatalln("Unable to search results:", err)
	}

	log.Printf("Got %d session reults\n", len(results))
	tw := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)

	tw.Write([]byte("ID\tTime\tStart\tFinish\tG/L\tInc.\tWinner\n"))

	for _, r := range results {
		tw.Write([]byte(fmt.Sprintf(
			"%d\t%s\t%2d\t%2d\t%+ 3d\t%4d\t%s\n",
			r.SubsessionID,
			time.Time(r.RawStartTime).Format(time.RFC822),
			r.StartingPos,
			r.FinishPos,
			r.StartingPos-r.FinishPos,
			r.Incidents,
			r.WinnerName,
		)))
	}

	tw.Flush()
}
