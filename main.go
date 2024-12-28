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

func SyncGasRecords(v lubelogger.Vehicle, rt roadtrip.CSV) error {
	slog.Debug("Synching fillups")

	var (
		rtInsertQueue []roadtrip.Fuel
		llInsertQueue []lubelogger.GasRecord
	)

	llGasRecords, err := lubelogger.GasRecords(v.ID)
	if err != nil {
		return err
	}

	for i, rtf := range rt.FuelRecords {
		rtComparator := rtf.Comparator()

		gr, err := llGasRecords.FindGasRecord(rtComparator)
		if err != nil {
			slog.Error("FindGasRecord failed",
				"error", err,
			)
			break
		}

		llComparator := gr.Comparator()

		if llComparator == rtComparator {
			slog.Debug("RT Fillup found in LubeLogger",
				"rtIndex", i,
				"comparator", rtComparator,
				"llOdometer", gr.Odometer,
			)

		} else {
			slog.Debug("RT Fillup not in LubeLogger, Enqueing",
				"rtIndex", i,
				"comparator", rtComparator,
			)
			rtInsertQueue = append(rtInsertQueue, rtf)
		}
	}

	slog.Info("Missing Fuel records enqueued",
		"rtCount", len(rtInsertQueue),
		"llCount", len(llInsertQueue),
	)

	for i, e := range rtInsertQueue {
		slog.Debug("Adding Road Trip Fillup to LubeLogger",
			"index", i,
			"fuelEntry", e,
		)

		gr, err := TransformRoadTripFuelToLubeLogger(e)
		if err != nil {
			slog.Error("Failed Adding Road Trip Fillup to LubeLogger", "error", err)
			break
		}

		response, err := lubelogger.AddGasRecord(v.ID, gr)
		if err != nil {
			slog.Error("Failed Adding Road Trip Fillup to LubeLogger",
				"index", i,
				"fuelEntry", e,
				"error", err,
			)
			break
		}
		slog.Info("Added Road Trip Fillup to LubeLogger",
			"index", i,
			"fuelEntry", e,
			"success", response.Success,
			"message", response.Message,
		)
	}

	return nil
}

func TransformRoadTripFuelToLubeLogger(rtf roadtrip.Fuel) (lubelogger.GasRecord, error) {
	gr := lubelogger.GasRecord{}
	date := roadtrip.ParseDate(rtf.Date)

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

func main() {
	_ = slog.SetLogLoggerLevel(slog.LevelInfo)

	var roadtripCSVPath = flag.String("csvpath", "./testdata/CSV", "Location of Road Trip CSV files")
	var debugMode = flag.Bool("v", false, "Verbose logging")

	flag.Parse()

	if *debugMode {
		_ = slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	lubelogger.Init(API_URI, AUTHORIZATION)

	vehicles, err := lubelogger.Vehicles()
	if err != nil {
		slog.Error("Error loading lubelogger Vehicles", "error", err)
		os.Exit(1)
	}

	for _, v := range vehicles {
		filename := v.CSVFilename()

		slog.Info("Evaluating lubelogger vehicle",
			"roadTripCSV", filename,
		)

		if filename != "" {
			rt, err := roadtrip.NewFromFile(filepath.Join(*roadtripCSVPath, filename))

			if err != nil {
				slog.Error("Error loading vehicle",
					"filename", filename,
					"error", err,
				)
				break
			}

			err = SyncGasRecords(v, rt)
			if err != nil {
				slog.Error("Error synching fuel records",
					"error", err,
				)
				break
			}
		}
	}
}
