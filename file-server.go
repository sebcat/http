package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

func fsDirList(h http.Handler, blockDirList bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL)
		if blockDirList && strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func fsMain(args []string) (err error) {
	var f flag.FlagSet
	path := f.String("path", ".", "path to HTTP root")
	listen := f.String("listen", ":8080", "listening directive")
	noDirList := f.Bool("no-dir-list", false, "disable directory listing")
	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		} else {
			return err
		}
	}

	if len(*path) == 0 {
		f.PrintDefaults()
		return nil
	}

	http.Handle("/", fsDirList(http.FileServer(http.Dir(*path)), *noDirList))
	log.Println("Started")
	return http.ListenAndServe(*listen, nil)
}
