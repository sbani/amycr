package cmd

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/sbani/gcr/cmd/cli/server"
	"github.com/sbani/gcr/pkg"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Run:   runServerStart,
}

// Port holds the server's listening port'
var Port int32

func init() {
	serverCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().Int32VarP(&Port, "port", "p", 4444, "server port")
}

func runServerStart(cmd *cobra.Command, args []string) {
	router := httprouter.New()
	serverHandler := &server.Handler{}
	serverHandler.Start(c, router)

	http.Handle("/", router)

	var srv = http.Server{
		Addr: c.GetAddress(),
	}

	logrus.Infof("Starting server on %s", srv.Addr)
	err := srv.ListenAndServe()

	pkg.Must(errors.Wrap(err, "Could not start server"))
}
