package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/leoadamek/irapi"
)

func main() {
	ctx := context.Background()

	uid := flag.Uint64("u", 0, "User ID")
	cp := flag.String("c", "", "Credential file path")
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

	log.Println("Created API")

	if err := api.Login(ctx); err != nil {
		log.Fatalf("Failed to log into iRacing: %s", err)
	}

	stats, err := api.GetCareerStats(ctx, *uid)

	if err != nil {
		log.Fatalln("Failed to get stats:", err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)

	tw.Write([]byte("Category\tStarts\tLaps\tAvg. Finish\tWins\n"))

	for _, cat := range stats {
		fmt.Fprintf(tw, "%s\t%8d\t%8d\t%2d\t%4d\n", cat.Category, cat.Starts, cat.TotalLaps, cat.AverageFinish, cat.Wins)
	}

	tw.Flush()

}
