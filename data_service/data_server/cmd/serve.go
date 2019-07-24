package cmd

import (
	"github.com/nwpc-oper/nwpc-data-client/data_service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
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
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	data_service.RegisterNWPCDataServiceServer(s, &data_service.NWPCDataServer{
		ConfigDir: configDir,
	})
	log.Printf("begin to serve...\n")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
