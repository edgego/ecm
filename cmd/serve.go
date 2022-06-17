package cmd

import (
	"fmt"
	"github.com/cnrancher/autok3s/pkg/server/handler"
	"github.com/cnrancher/autok3s/pkg/server/ui"
	"github.com/gorilla/mux"
	responsewriter "github.com/rancher/apiserver/pkg/middleware"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run as daemon and serve HTTP/HTTPS request",
	}

	bindPort    = "8080"
	bindAddress = "0.0.0.0"
)

func init() {
	serveCmd.Flags().StringVar(&bindPort, "bind-port", bindPort, "HTTP/HTTPS bind port")
	serveCmd.Flags().StringVar(&bindAddress, "bind-address", bindAddress, "HTTP/HTTPS bind address")
}

// ServeCommand serve command.
func ServeCommand() *cobra.Command {
	serveCmd.Run = func(cmd *cobra.Command, args []string) {
		//router := server.Start()
		router := mux.NewRouter()

		//add by edgego
		router.HandleFunc("/cluster/{cluster}/node/{node}/instance/{instance}", handler.DeleteNode).Methods(http.MethodDelete)
		router.HandleFunc("/cluster/{cluster}/provider/{provider}", handler.DeleteCluster).Methods(http.MethodDelete)
		router.HandleFunc("/cluster/{cluster}", handler.CreateCluster).Methods(http.MethodPost)
		router.HandleFunc("/clusters/{provider}", handler.GetClustersList).Methods(http.MethodGet)
		router.HandleFunc("/cluster/{cluster}/nodes", handler.ListClusterDetail).Methods(http.MethodGet)
		router.HandleFunc("/cluster/join/{cluster}", handler.JoinCluster).Methods(http.MethodPost)
		router.HandleFunc("/cluster/ping/{cluster}/{host}", handler.PingHost).Methods(http.MethodPost)

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

		// start kube-explorer for K3s clusters
		//go common.InitExplorer()
		logrus.Infof("run as daemon, listening on %s:%s", bindAddress, bindPort)
		logrus.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", bindAddress, bindPort), router))
	}

	return serveCmd
}
