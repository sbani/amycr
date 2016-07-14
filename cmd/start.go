package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/sbani/gcr/cmd/cli/server"
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
	serverHandler := &server.Handler{}
	serverHandler.Start(c, echo.New())

	logrus.Infof("Starting server on %s", srv.Addr)
	e.Run(standard.New(c.GetAddress()))
}
