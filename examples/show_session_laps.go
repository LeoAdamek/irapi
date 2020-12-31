package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/leoadamek/irapi"
)

func main() {
	ctx := context.Background()

	sid := flag.Uint64("s", 0, "Session ID")
	eid := flag.Uint64("e", 0, "Entrant ID")
	ph := flag.Uint64("p", 0, "Phase number")
	cp := flag.String("c", "", "Credential file path")

	debug := flag.Bool("d", false, "Enable Debug Output")

	flag.Parse()

	if eid == nil || *eid == 0 {
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
		api.BeforeRequest(func(_ context.Context, req *http.Request) error {

			log.Println("Request")
			req.Header.Write(os.Stdout)
			return nil
		})

		api.AfterResponse(func(_ context.Context, _ *http.Request, res *http.Response) error {
			log.Println("Response")
			log.Println(res.Status)
			res.Header.Write(os.Stdout)
			res.Write(os.Stdout)
			return nil
		})
	}

	laps, err := api.GetLaps(ctx, *sid, *eid, *ph)

	if err != nil {
		log.Fatalln("Unable to get laps:", err)
	}

	for _, l := range laps {
		fmt.Printf("%+#v\n", l)
	}

}
