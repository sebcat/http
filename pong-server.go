package main

import (
	"flag"
	"log"
	"net/http"
)

func psMain(args []string) error {
	var f flag.FlagSet
	var listenStr = f.String("listen", ":8080", "listen directive")
	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		} else {
			return err
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong\n"))
	})

	log.Println("Started")
	return http.ListenAndServe(*listenStr, nil)
}
