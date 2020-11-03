package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/leoadamek/irapi"
)

func main() {
	ctx := context.Background()

	ssid := flag.Uint64("s", 0, "Sub-Session ID")
	flag.Parse()

	if ssid == nil || *ssid == 0 {
		flag.Usage()
		log.Fatalln("No Subsession Provided")
	}

	api := irapi.New(ctx, irapi.EnvironmentCredentialsProvider)

	log.Println("Created API")

	if err := api.Login(ctx); err != nil {
		log.Fatalf("Failed to log into iRacing: %s", err)
	}

	result, err := api.GetSubSessionResult(ctx, *ssid)

	if err != nil {
		log.Fatalln("Failed to get session results:", err)
	}

	filteredResults := byFinishPos(filterResultsBySession(result.Results, "RACE"))
	sort.Sort(filteredResults)

	tw := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)

	tw.Write([]byte("#\tName\tStart\tFinish\tCh\tiRating\tiRating Change\n"))

	for _, r := range filteredResults {
		tw.Write([]byte(fmt.Sprintf(
			"%s\t%s\t%d\t%d\t%+3d\t%5d\t%+ 5d\n",
			r.CarNumber,
			r.Name,
			r.StartPosition,
			r.FinishPosition,
			r.StartPosition-r.FinishPosition,
			r.NewIRating,
			r.NewIRating-r.OldIRating,
		)))
	}

	tw.Flush()

}

func filterResultsBySession(results []irapi.DriverResult, session string) []irapi.DriverResult {
	var out []irapi.DriverResult

	for _, r := range results {
		if r.SessionName == session {
			out = append(out, r)
		}
	}

	return out
}

type byFinishPos []irapi.DriverResult

func (b byFinishPos) Less(i, j int) bool {
	return b[i].FinishPosition < b[j].FinishPosition
}

func (b byFinishPos) Len() int { return len(b) }

func (b byFinishPos) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
