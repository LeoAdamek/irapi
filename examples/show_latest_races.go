package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/leoadamek/irapi"
)

func main() {
	ctx := context.Background()

	uid := flag.Uint64("u", 0, "User ID")
	cp := flag.String("c", "", "Credential file path")
	debug := flag.Bool("d", false, "Enable Debug Output")

	flag.Parse()

	if uid == nil || *uid == 0 {
		flag.Usage()
		log.Fatalln("No User ID Provided")
	}

	var creds irapi.CredentialsProvider

	if *cp != "" {
		creds = irapi.FileCredentialsProvider(*cp)
	} else {
		creds = irapi.EnvironmentCredentialsProvider
	}

	api := irapi.New(creds)

	if *debug {

		api.BeforeRequest(func(ctx context.Context, req *http.Request) error {

			return nil
		})

		api.AfterResponse(func(_ context.Context, req *http.Request, res *http.Response) error {

			log.Println("Response")
			res.Write(log.Writer())

			log.Println()

			return nil
		})
	}

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

	log.Println("Searching for the last 100 results from the past two weeks...")

	params.UserID = *uid

	results, err := api.SearchResults(ctx, params)

	if err != nil {
		log.Fatalln("Unable to search results:", err)
	}

	log.Printf("Got %d session reults\n", len(results))
	tw := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)

	tw.Write([]byte("ID\tTime\tStart\tFinish\tG/L\tInc.\tWinner\n"))

	for _, r := range results {
		fmt.Fprintf(tw,
			"%d\t%s\t%2d\t%2d\t%+ 3d\t%4d\t%s\n",
			r.SubsessionID,
			time.Time(r.RawStartTime).Format(time.RFC822),
			r.StartingPos,
			r.FinishPos,
			r.StartingPos-r.FinishPos,
			r.Incidents,
			r.WinnerName,
		)
	}

	tw.Flush()
}
