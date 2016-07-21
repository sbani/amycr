package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
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
	e := echo.New()

	serverHandler := new(server.Handler)
	serverHandler.Start(c, e)

	addr := c.GetAddress()

	logrus.Infof("Starting server on %s", addr)
	e.Run(fasthttp.New(addr))
}
