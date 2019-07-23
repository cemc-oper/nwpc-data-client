package cmd

import (
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/data_service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var (
	ConfigDir = ""
	Address   = ""
)

var rootCmd = &cobra.Command{
	Use:   "nwpc_data_server",
	Short: "A data server to get NWPC data.",
	Run: func(cmd *cobra.Command, args []string) {
		runServer(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&ConfigDir, "config-dir", "", "Config dir")
	rootCmd.Flags().StringVar(&Address, "address", ":30105", "server address")

	_ = rootCmd.MarkFlagRequired("config-dir")
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
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
