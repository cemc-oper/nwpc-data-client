package cmd

import (
	"context"
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

func init() {
	rootCmd.AddCommand(serviceCommand)

	serviceCommand.Flags().StringVar(&dataType, "data-type", "",
		"Data type used to locate config file path in config dir.")

	serviceCommand.Flags().StringVar(&locationLevels, "location-level", "",
		"Location levels, split by ',', such as 'runtime,archive'.")

	serviceCommand.Flags().StringVar(&startTimeSting, "start-time", "",
		"start time, YYYYMMDDHH, such as 2020021400")
	serviceCommand.Flags().StringVar(&forecastTimeString, "forecast-time", "",
		"forecast time, FFFh, such as 0h, 120h")

	serviceCommand.Flags().StringVar(&serviceAddress, "address", "",
		"serviceAddress of nwpc_data_server.")

	serviceCommand.Flags().StringVar(&serviceAction, "action", "",
		"service action, such as findDataPath, getDataFileInfo, downloadDataFile")

	serviceCommand.Flags().StringVar(&outputDirectory, "output-dir", "",
		"output file directory, default is work directory.")

	serviceCommand.MarkFlagRequired("data-type")
	serviceCommand.MarkFlagRequired("start-time")
	serviceCommand.MarkFlagRequired("forecast-time")
	serviceCommand.MarkFlagRequired("address")
	serviceCommand.MarkFlagRequired("action")
}

const serviceCommandDocString = `nwpc_data_client service
Find or get data from nwpc_data_server.
`

var serviceCommand = &cobra.Command{
	Use:   "service",
	Short: "Find data path using nwpc_data_server on HPC.",
	Long:  serviceCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if serviceAction == "downloadDataFile" {
			cmd.MarkFlagRequired("output-dir")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		_, err := common.ParseStartTime(startTimeSting)
		if err != nil {
			log.Errorf("check startTime failed: %s", err)
			return
		}

		_, err = common.ParseForecastTime(forecastTimeString)
		if err != nil {
			log.Errorf("check forecastTime failed: %s", err)
			return
		}

		conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
		if err != nil {
			log.WithFields(log.Fields{
				"component": "service",
				"event":     "connect",
			}).Fatalf("can't not connect to service: %v", err)
		}

		defer conn.Close()

		c := data_service.NewNWPCDataServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel()

		locationLevels := strings.Split(locationLevels, ",")

		if serviceAction == "findDataPath" {
			r, err := c.FindDataPath(ctx, &data_service.DataRequest{
				DataType:       dataType,
				LocationLevels: locationLevels,
				StartTime:      startTimeSting,
				ForecastTime:   forecastTimeString,
			})
			if err != nil {
				log.WithFields(log.Fields{
					"component": "service",
					"event":     "remote-run",
				}).Fatalf("could not run FindDataPath: %v", err)
			}

			fmt.Printf("%s\n%s\n", r.LocationType, r.Location)
		} else if serviceAction == "getDataFileInfo" {
			r, err := c.GetDataFileInfo(ctx, &data_service.DataRequest{
				DataType:       dataType,
				LocationLevels: locationLevels,
				StartTime:      startTimeSting,
				ForecastTime:   forecastTimeString,
			})
			if err != nil {
				log.WithFields(log.Fields{
					"component": "service",
					"event":     "remote-run",
				}).Fatalf("could not run GetDataFileInfo: %v", err)
			}

			if r.Status == data_service.StatusCode_Success {
				fmt.Printf("%s\n%d\n", r.GetFilePath(), r.GetFileSize())
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", r.GetErrorMessage())
			}
		} else if serviceAction == "downloadDataFile" {

			r, err := c.GetDataFileInfo(ctx, &data_service.DataRequest{
				DataType:       dataType,
				LocationLevels: locationLevels,
				StartTime:      startTimeSting,
				ForecastTime:   forecastTimeString,
			})
			if err != nil {
				log.WithFields(log.Fields{
					"component": "service",
					"event":     "remote-run",
				}).Fatalf("could not run GetDataFileInfo: %v", err)
			}

			if r.Status != data_service.StatusCode_Success {
				fmt.Fprintf(os.Stderr, "%s\n", r.GetErrorMessage())
				os.Exit(2)
			}

			_, remoteFileName := filepath.Split(r.GetFilePath())
			totalLength := float64(r.GetFileSize())
			outputFilePath := filepath.Join(outputDirectory, remoteFileName)

			common.PrepareLocalDir(outputFilePath)

			stream, err := c.DownloadDataFile(ctx, &data_service.DataRequest{
				DataType:       dataType,
				LocationLevels: locationLevels,
				StartTime:      startTimeSting,
				ForecastTime:   forecastTimeString,
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
					log.WithFields(log.Fields{
						"component": "service",
						"event":     "download",
					}).Fatalf("%v.DownloadFile(_) = _, %v", c, err)
				}
				f.Write(chunk.Chunk)
				currentSize += float64(chunk.ChunkLength)
				log.WithFields(log.Fields{
					"component": "service",
					"event":     "download",
				}).Printf("%0.2f%%", currentSize*100/totalLength)
			}

		} else {
			log.WithFields(log.Fields{
				"component": "service",
				"event":     "remote-run",
			}).Fatalf("service action is not supported: %s\n", serviceAction)
		}

	},
}
