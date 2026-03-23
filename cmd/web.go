package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gdisw/resume/pkg/env"
	xhttp "github.com/gdisw/resume/pkg/http"
	"github.com/gdisw/resume/pkg/http/session"
	"github.com/gdisw/resume/pkg/http/view"
	"github.com/meehow/securebytes"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessions := session.NewStore(
			securebytes.New(
				[]byte(os.Getenv("COOKIE_HASH_KEY")+os.Getenv("COOKIE_BLOCK_KEY")),
				securebytes.GOBSerializer{},
			),
			os.Getenv("COOKIE_NAME"),
		)

		view.PutStaticConfig()
		view.LoadBase(".")

		if env.IsLocal() {
			view.SetRefreshViewEnabled()
		}

		router := xhttp.NewRouter(sessions)
		router.Attach(xhttp.Health{})
		router.AttachProtected(xhttp.Home{})

		port, err := strconv.Atoi(os.Getenv("WEB_PORT"))
		if err != nil {
			port = 8080
		}
		bind := fmt.Sprintf(":%d", port)
		slog.Info("Listening on", "bind", bind)
		return http.ListenAndServe(bind, router)
	},
}
