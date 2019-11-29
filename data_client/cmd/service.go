package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	"github.com/nwpc-oper/nwpc-data-client/data_service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ServiceAddress  = ""
	ServiceAction   = ""
	OutputDirectory = "."
)

func init() {
	rootCmd.AddCommand(serviceCommand)

	serviceCommand.Flags().StringVar(&DataType, "data-type", "",
		"Data type used to locate config file path in config dir.")

	serviceCommand.Flags().StringVar(&LocationLevels, "location-level", "",
		"Location levels, split by ',', such as 'runtime,archive'.")

	serviceCommand.Flags().StringVar(&ServiceAddress, "address", "",
		"ServiceAddress of nwpc_data_server.")

	serviceCommand.Flags().StringVar(&ServiceAction, "action", "",
		"service action, such as findDataPath, getDataFileInfo, downloadDataFile")

	serviceCommand.Flags().StringVar(&OutputDirectory, "output-dir", "",
		"output file directory, default is work directory.")

	serviceCommand.MarkFlagRequired("data-type")
	serviceCommand.MarkFlagRequired("address")
	serviceCommand.MarkFlagRequired("action")
}

const serviceCommandDocString = `nwpc_data_client service
Find or get data from nwpc_data_server.

Args:
    start_time: YYYYMMDDHH, such as 2018080100
    forecast_time: time duration, such as 0h, 120h`

var serviceCommand = &cobra.Command{
	Use:   "service",
	Short: "Find data path using nwpc_data_server on HPC.",
	Long:  serviceCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires two arguments")
		}

		if ServiceAction == "downloadDataFile" {
			cmd.MarkFlagRequired("output-dir")
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

		locationLevels := strings.Split(LocationLevels, ",")

		if ServiceAction == "findDataPath" {
			r, err := c.FindDataPath(ctx, &data_service.DataRequest{
				DataType:       DataType,
				LocationLevels: locationLevels,
				StartTime:      args[0],
				ForecastTime:   args[1],
			})
			if err != nil {
				log.Fatalf("could not run FindDataPath: %v", err)
			}

			fmt.Printf("%s\n%s\n", r.LocationType, r.Location)
		} else if ServiceAction == "getDataFileInfo" {
			r, err := c.GetDataFileInfo(ctx, &data_service.DataRequest{
				DataType:       DataType,
				LocationLevels: locationLevels,
				StartTime:      args[0],
				ForecastTime:   args[1],
			})
			if err != nil {
				log.Fatalf("could not run GetDataFileInfo: %v", err)
			}

			if r.Status == data_service.StatusCode_Success {
				fmt.Printf("%s\n%d\n", r.GetFilePath(), r.GetFileSize())
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", r.GetErrorMessage())
			}
		} else if ServiceAction == "downloadDataFile" {

			r, err := c.GetDataFileInfo(ctx, &data_service.DataRequest{
				DataType:       DataType,
				LocationLevels: locationLevels,
				StartTime:      args[0],
				ForecastTime:   args[1],
			})
			if err != nil {
				log.Fatalf("could not run GetDataFileInfo: %v", err)
			}

			if r.Status != data_service.StatusCode_Success {
				fmt.Fprintf(os.Stderr, "%s\n", r.GetErrorMessage())
				os.Exit(2)
			}

			_, remoteFileName := filepath.Split(r.GetFilePath())
			totalLength := float64(r.GetFileSize())
			outputFilePath := filepath.Join(OutputDirectory, remoteFileName)

			common.PrepareLocalDir(outputFilePath)

			stream, err := c.DownloadDataFile(ctx, &data_service.DataRequest{
				DataType:       DataType,
				LocationLevels: locationLevels,
				StartTime:      args[0],
				ForecastTime:   args[1],
			})

			f, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY, 0644)
			defer f.Close()

			var currentSize float64 = 0
			for {
				chunk, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("%v.DownloadFile(_) = _, %v", c, err)
				}
				f.Write(chunk.Chunk)
				currentSize += float64(chunk.ChunkLength)
				log.Printf("%0.2f%%", currentSize*100/totalLength)
			}

		} else {
			log.Fatalf("service action is not supported: %s\n", ServiceAction)
		}

	},
}
