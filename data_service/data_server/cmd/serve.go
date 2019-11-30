package cmd

import (
	"github.com/nwpc-oper/nwpc-data-client/data_service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
)

var (
	ConfigDir = ""
	Address   = ""
)

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "run data server",
	Long:  "run data server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(serveCommand)
	serveCommand.Flags().StringVar(&ConfigDir, "config-dir", "", "Config dir")
	serveCommand.Flags().StringVar(&Address, "address", ":33383", "server address")

	_ = serveCommand.MarkFlagRequired("config-dir")
}

func runServer(cmd *cobra.Command, args []string) {
	runGRPCServer(ConfigDir, Address)
}

func runGRPCServer(configDir string, address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "serve",
			"event":     "connect",
		}).Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	data_service.RegisterNWPCDataServiceServer(s, &data_service.NWPCDataServer{
		ConfigDir: configDir,
	})

	log.WithFields(log.Fields{
		"component": "serve",
		"event":     "connect",
	}).Printf("nwpc_data_server begin to serve at %s ...\n", address)
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{
			"component": "serve",
			"event":     "connect",
		}).Fatalf("nwpc_data_server failed to serve: %v", err)
	}
}
