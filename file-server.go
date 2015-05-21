package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

func fsDirList(h http.Handler, allowDirList bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL)
		if !allowDirList && strings.HasSuffix(r.URL.Path, "/") {
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
	dirList := f.Bool("dir-list", false, "enable directory listing")
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

	http.Handle("/", fsDirList(http.FileServer(http.Dir(*path)), *dirList))
	log.Println("Started")
	return http.ListenAndServe(*listen, nil)
}
