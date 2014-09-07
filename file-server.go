package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

func fsNoDirList(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL)
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func fsDirList(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL)
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

	if *noDirList {
		http.Handle("/", fsNoDirList(http.FileServer(http.Dir(*path))))
	} else {
		http.Handle("/", fsDirList(http.FileServer(http.Dir(*path))))
	}

	log.Println("Started")
	return http.ListenAndServe(*listen, nil)
}
