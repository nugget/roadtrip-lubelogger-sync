package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/nugget/roadtrip-go/roadtrip"
	"github.com/nugget/roadtrip-lubelogger-sync/lubelogger"
)

var (
	logger   *slog.Logger
	logLevel *slog.LevelVar
)

func rtfComparator(f roadtrip.FuelRecord) string {
	return fmt.Sprintf("%07d", int64(f.Odometer))
}

func SyncGasRecords(v lubelogger.Vehicle, rt roadtrip.Vehicle) error {
	logger.Debug("Synching fillups")

	var (
		rtInsertQueue []roadtrip.FuelRecord
		llInsertQueue []lubelogger.GasRecord
	)

	llGasRecords, err := lubelogger.GasRecords(v.ID)
	if err != nil {
		return err
	}

	for i, rtf := range rt.FuelRecords {
		rtComparator := rtfComparator(rtf)

		gr, err := llGasRecords.FindGasRecord(rtComparator)
		if err != nil {
			logger.Error("FindGasRecord failed",
				"error", err,
			)
			break
		}

		llComparator := gr.Comparator()

		if llComparator == rtComparator {
			logger.Debug("RT Fillup found in LubeLogger",
				"rtIndex", i,
				"comparator", rtComparator,
				"llOdometer", gr.Odometer,
			)

		} else {
			logger.Debug("RT Fillup not in LubeLogger, Enqueing",
				"rtIndex", i,
				"comparator", rtComparator,
			)
			rtInsertQueue = append(rtInsertQueue, rtf)
		}
	}

	logger.Info("Missing Fuel records enqueued",
		"rtCount", len(rtInsertQueue),
		"llCount", len(llInsertQueue),
	)

	for i, e := range rtInsertQueue {
		logger.Debug("Adding Road Trip Fillup to LubeLogger",
			"index", i,
			"fuelEntry", e,
		)

		gr, err := TransformRoadTripFuelToLubeLogger(e)
		if err != nil {
			logger.Error("Failed Adding Road Trip Fillup to LubeLogger", "error", err)
			break
		}

		response, err := lubelogger.AddGasRecord(v.ID, gr)
		if err != nil {
			logger.Error("Failed Adding Road Trip Fillup to LubeLogger",
				"index", i,
				"fuelEntry", e,
				"error", err,
			)
			break
		}
		logger.Info("Added Road Trip Fillup to LubeLogger",
			"index", i,
			"fuelEntry", e,
			"success", response.Success,
			"message", response.Message,
		)
	}

	return nil
}

func TransformRoadTripFuelToLubeLogger(rtf roadtrip.FuelRecord) (lubelogger.GasRecord, error) {
	gr := lubelogger.GasRecord{}
	date, err := rtf.Date.MustParse()
	if err != nil {
		return lubelogger.GasRecord{}, err
	}

	gr.Date = lubelogger.FormatDate(date)
	gr.Odometer = fmt.Sprintf("%d", int(rtf.Odometer))
	gr.FuelConsumed = fmt.Sprintf("%0.3f", rtf.FillAmount)
	gr.Cost = fmt.Sprintf("%0.2f", rtf.TotalPrice)
	gr.FuelEconomy = fmt.Sprintf("%f", rtf.MPG)
	gr.MissedFuelUp = "False"
	gr.Notes = rtf.Note

	if rtf.PartialFill != "" {
		gr.IsFillToFull = "False"
	} else {
		gr.IsFillToFull = "True"
	}

	gr.Notes += fmt.Sprintf("\n%0.02f gallons @ $%0.2f from %s", rtf.FillAmount, rtf.PricePerUnit, rtf.Location)

	gr.Notes = strings.Trim(gr.Notes, " \t\r\n")

	location := lubelogger.ExtraField{}
	location.Name = "Location"
	location.Value = rtf.Location
	gr.ExtraFields = append(gr.ExtraFields, location)

	return gr, nil
}

func setupLogs() {
	logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)

	handlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelInfo)
}

func main() {
	setupLogs()

	var roadtripCSVPath = flag.String("csvpath", "./testdata/CSV", "Location of Road Trip CSV files")
	var debugMode = flag.Bool("v", false, "Verbose logging")

	flag.Parse()

	options := roadtrip.VehicleOptions{
		Logger: logger,
	}

	if *debugMode {
		// AddSource: true here
		slog.SetLogLoggerLevel(slog.LevelDebug)
		logLevel.Set(slog.LevelDebug)
	}

	lubelogger.Init(API_URI, AUTHORIZATION)

	logger.Debug("Loading vehicles from LubeLogger API",
		"uri", API_URI,
	)

	vehicles, err := lubelogger.Vehicles()
	if err != nil {
		logger.Error("Error loading lubelogger Vehicles", "error", err)
		os.Exit(1)
	}

	for _, v := range vehicles {
		filename := v.CSVFilename()

		logger.Info("Evaluating lubelogger vehicle",
			"roadTripCSV", filename,
		)

		if filename != "" {
			rt, err := roadtrip.NewVehicleFromFile(filepath.Join(*roadtripCSVPath, filename), options)

			if err != nil {
				logger.Error("Error loading vehicle",
					"filename", filename,
					"error", err,
				)
				break
			}

			err = SyncGasRecords(v, rt)
			if err != nil {
				logger.Error("Error synching fuel records",
					"error", err,
				)
				break
			}
		}
	}
}
