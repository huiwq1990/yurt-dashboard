package main

import (
	"flag"
	"fmt"
	"kubeapi/cmd/app"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	version    = false
	configFile = ""
	rootCmd    = &cobra.Command{
		Use:   "kubeapi",
		Short: "",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run()
		},
	}
)

func main() {
	versionInfoCmd := &cobra.Command{
		Use:   "version",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%#v\n", "aaa")

		},
	}
	rootCmd.AddCommand(versionInfoCmd)
	flags := rootCmd.Flags()
	flags.AddFlagSet(InitFlag())
	local := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(local)
	flags.AddGoFlagSet(local)

	if err := rootCmd.Execute(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}

func Run() error {
	//stopCh := setupSignalHandler()

	err := app.InitRouter()
	if err != nil {
		klog.ErrorS(err, "init router fail")
		return err
	}

	app.Router.Run(fmt.Sprintf("%s:%d", "0.0.0.0", 8080))
	return nil
}

var onlyOneSignalHandler = make(chan struct{})
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func setupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()
	return stop
}

func InitFlag() *pflag.FlagSet {
	flags := pflag.NewFlagSet("kubeapi", pflag.ContinueOnError)
	flags.BoolVarP(&version, "version", "", false, "print app version")
	return flags
}
