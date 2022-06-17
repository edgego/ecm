package server

import (
	"github.com/cnrancher/autok3s/pkg/server/handler"
	"github.com/cnrancher/autok3s/pkg/server/ui"
	"net/http"

	"github.com/gorilla/mux"
	responsewriter "github.com/rancher/apiserver/pkg/middleware"
	"net/http/pprof"
)

// Start starts daemon.
func Start() http.Handler {
	
	router := mux.NewRouter()
	router.Use(webRequestMiddleware)

	middleware := responsewriter.Chain{
		responsewriter.Gzip,
		responsewriter.FrameOptions,
		responsewriter.CacheMiddleware("json", "js", "css", "svg", "png", "woff", "woff2"),
		ui.ServeNotFound,
	}
	router.PathPrefix("/").Handler(middleware.Handler(http.StripPrefix("/", ui.Serve())))
	router.Path("/").HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, "/", http.StatusFound)
	})

	//add by edgego
	router.HandleFunc("/cluster/{cluster}/node/{node}/instance/{instance}", handler.DeleteNode).Methods(http.MethodDelete)
	router.HandleFunc("/cluster/{cluster}/provider/{provider}", handler.DeleteCluster).Methods(http.MethodDelete)
	router.HandleFunc("/cluster/{cluster}", handler.CreateCluster).Methods(http.MethodPost)
	router.HandleFunc("/clusters/{provider}", handler.GetClustersList).Methods(http.MethodGet)
	router.HandleFunc("/cluster/join/{cluster}", handler.JoinCluster).Methods(http.MethodPost)

	// profiling handlers for pprof under /debug/pprof.
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)

	// Manually add support for paths linked to by index page at /debug/pprof/
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))

	return router
}

func webRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
