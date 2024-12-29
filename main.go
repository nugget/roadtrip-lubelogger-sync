package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
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

	llGasRecords, errGR := lubelogger.GasRecords(v.ID)
	if errGR != nil {
		return errGR
	}

	for i, rtf := range rt.FuelRecords {
		rtComparator := rtfComparator(rtf)

		gr, errFGR := llGasRecords.FindGasRecord(rtComparator)
		if errFGR != nil {
			logger.Error("FindGasRecord failed",
				"error", errFGR,
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

		gr, errR2L := TransformRoadTripFuelToLubeLogger(e)
		if errR2L != nil {
			logger.Error("Failed Adding Road Trip Fillup to LubeLogger", "error", errR2L)
			break
		}

		response, errAGR := lubelogger.AddGasRecord(v.ID, gr)
		if errAGR != nil {
			logger.Error("Failed Adding Road Trip Fillup to LubeLogger",
				"index", i,
				"fuelEntry", e,
				"error", errAGR,
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
	gr.Odometer = strconv.Itoa(int(rtf.Odometer))
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

func setupSecrets() (string, string) {
	// apiURI        string = "https://lubelogger.example.com/api"
	// authorization string = "Basic BASIC_AUTH_TOKEN_GOES_HERE"

	apiURI := os.Getenv("API_URI")
	authorization := os.Getenv("AUTHORIZATION")

	return apiURI, authorization
}

func main() {
	var (
		roadtripCSVPath = flag.String("csvpath", "./testdata/CSV", "Location of Road Trip CSV files")
		debugMode       = flag.Bool("v", false, "Verbose logging")
	)

	apiURI, authorization := setupSecrets()

	setupLogs()
	flag.Parse()

	options := roadtrip.VehicleOptions{
		Logger: logger,
	}

	if *debugMode {
		// AddSource: true here
		slog.SetLogLoggerLevel(slog.LevelDebug)
		logLevel.Set(slog.LevelDebug)
	}

	lubelogger.Init(apiURI, authorization, logger)

	logger.Debug("Loading vehicles from LubeLogger API",
		"uri", apiURI,
	)

	vehicles, errLV := lubelogger.Vehicles()
	if errLV != nil {
		logger.Error("Error loading lubelogger Vehicles", "error", errLV)
		os.Exit(1)
	}

	for _, v := range vehicles {
		filename := v.CSVFilename()

		logger.Info("Evaluating lubelogger vehicle",
			"filename", filename,
		)

		if filename != "" {
			rt, errNV := roadtrip.NewVehicleFromFile(filepath.Join(*roadtripCSVPath, filename), options)

			if errNV != nil {
				logger.Error("Error loading vehicle",
					"filename", filename,
					"error", errNV,
				)
				break
			}

			errSGR := SyncGasRecords(v, rt)
			if errSGR != nil {
				logger.Error("Error synching fuel records",
					"error", errSGR,
				)
				break
			}
		}
	}
}
