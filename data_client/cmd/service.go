package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/data_service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"time"
)

var (
	ServiceAddress = ""
	ServiceAction  = ""
)

func init() {
	rootCmd.AddCommand(serviceCommand)

	serviceCommand.Flags().StringVar(&DataType, "data-type", "",
		"Data type used to locate config file path in config dir.")
	serviceCommand.Flags().StringVar(&ServiceAddress, "address", "",
		"ServiceAddress of nwpc_data_server.")

	serviceCommand.Flags().StringVar(&ServiceAction, "action", "",
		"service action, such as findDataPath")

	serviceCommand.MarkFlagRequired("data-type")
	serviceCommand.MarkFlagRequired("address")
	serviceCommand.MarkFlagRequired("action")
}

var serviceCommand = &cobra.Command{
	Use:   "service",
	Short: "Find data path using nwpc_data_server on HPC.",
	Long:  "Find data path using nwpc_data_server on HPC.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires two arguments")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(ServiceAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		defer conn.Close()

		c := data_service.NewNWPCDataServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel()

		if ServiceAction == "findDataPath" {
			r, err := c.FindDataPath(ctx, &data_service.DataRequest{
				DataType:     DataType,
				StartTime:    args[0],
				ForecastTime: args[1],
			})
			if err != nil {
				log.Fatalf("could not run FindDataPath: %v", err)
			}

			fmt.Printf("%s\n%s\n", r.LocationType, r.Location)
		} else {
			log.Fatalf("service action is not supported: %s\n", ServiceAction)
		}

	},
}
