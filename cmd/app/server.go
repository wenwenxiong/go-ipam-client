package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
	"wenwenxiong/go-ipam-client/pkg/apiserver"
	"wenwenxiong/go-ipam-client/pkg/client/goipam"
)

var (
	version    bool
	endpoint   string
	listenPort string
	logger     = logrus.New()
	// Version shows the go-ipam-client binary version.
	Version string
	// GitSHA shows the  go-ipam-client binary code commit SHA on git.
	GitSHA string
)

func printVersionInfo() {
	logger.Infof("go-ipam-client Version: %s", Version)
	logger.Infof("Git SHA: %s", GitSHA)
	logger.Infof("Go Version: %s", runtime.Version())
	logger.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
}

func NewIpamClientCommand() *cobra.Command {
	var RootCmd = &cobra.Command{
		Use:   "go-ipam-client",
		Short: "go-ipam-client serve for provide restful api to go-ipam ",
		Long:  `The go-ipam-client,mini web server, is serve to provide restful api for go-ipam, it encapsulation  cidr ip usage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServer(listenPort, endpoint)
		},
		SilenceUsage: true,
	}
	RootCmd.Flags().BoolVarP(&version, "version", "v", false, "print version info")
	RootCmd.PersistentFlags().StringVarP(&listenPort, "listenPort", "p", "9111", "listen port, default is 9111")
	RootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:9090", "go-ipam endpoint, default is http://localhost:9090")

	usageFmt := "Usage:\n  %s\n"
	RootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
	})

	var version = &cobra.Command{
		Use:   "version",
		Short: "Version and Release information about the go-ipam-client",
		Run: func(cmd *cobra.Command, args []string) {
			printVersionInfo()
		},
	}
	RootCmd.AddCommand(version)
	return RootCmd
}

func RunServer(listenPort string, endpoint string) error {
	server, err := NewApiServer(listenPort, endpoint)
	if err != nil {
		return err
	}

	return server.Run()
}

func NewApiServer(listenPort string, endpoint string) (*apiserver.APIServer, error) {
	port, err := strconv.Atoi(listenPort)
	if err != nil {
		return nil, fmt.Errorf("listenPort must be right int format, error: %v", err)
	}
	s := &apiserver.APIServer{
		ListenPort: port,
		Endpoint:   endpoint,
	}
	if (strings.TrimSpace(s.Endpoint)) == "" {
		return nil, fmt.Errorf("go-ipam service address MUST not be empty, please check config endpoint")
	} else {
		goipamClient := goipam.NewGoipamClient(s.Endpoint)
		s.GoipamClient = goipamClient
	}
	router := gin.Default()
	apiserver.RegisterRoutes(router, s.GoipamClient)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", s.ListenPort),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}
	s.Server = srv

	return s, nil
}
