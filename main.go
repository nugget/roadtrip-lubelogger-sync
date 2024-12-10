package main

import (
	"fmt"
	"strings"

	"github.com/nugget/roadtrip-lubelogger-sync/lubelogger"
	"github.com/nugget/roadtrip-lubelogger-sync/roadtrip"

	log "github.com/sirupsen/logrus"
)

func SyncGasRecords(v lubelogger.Vehicle, rt roadtrip.CSV) error {
	v.Logrus().Trace("Synching fillups")

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
			log.WithError(err).Error("FindGasRecord failed")
			break
		}

		llComparator := gr.Comparator()

		if llComparator == rtComparator {
			log.WithFields(log.Fields{
				"rtIndex":    i,
				"comparator": rtComparator,
				"llOdometer": gr.Odometer,
			}).Trace("RT Fillup found in LubeLogger")

		} else {
			log.WithFields(log.Fields{
				"rtIndex":    i,
				"comparator": rtComparator,
			}).Trace("RT Fillup not in LubeLogger, Enqueing")
			rtInsertQueue = append(rtInsertQueue, rtf)
		}
	}

	log.WithFields(log.Fields{
		"rtCount": len(rtInsertQueue),
		"llCount": len(llInsertQueue),
	}).Info("Missing Fuel records enqueued")

	for i, e := range rtInsertQueue {
		e.Logrus().WithFields(log.Fields{
			"index": i,
		}).Trace("Adding Road Trip Fillup to LubeLogger")

		gr, err := TransformRoadTripFuelToLubeLogger(e)
		if err != nil {
			e.Logrus().WithError(err).Error("Failed Adding Road Trip Fillup to LubeLogger")
			break
		}

		response, err := lubelogger.AddGasRecord(v.ID, gr)
		if err != nil {
			e.Logrus().WithFields(log.Fields{
				"index": i,
				"error": err,
			}).Info("Failed Adding Road Trip Fillup to LubeLogger")
			break
		}
		e.Logrus().WithFields(log.Fields{
			"index":   i,
			"success": response.Success,
			"message": response.Message,
		}).Info("Added Road Trip Fillup to LubeLogger")
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
	log.SetLevel(log.DebugLevel)

	lubelogger.Init(API_URI, AUTHORIZATION)

	vehicles, err := lubelogger.Vehicles()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range vehicles {
		filename := v.CSVFilename()

		v.Logrus().WithFields(log.Fields{
			"roadTripCSV": filename,
		}).Info("Evaluating lubelogger vehicle")

		if filename != "" {
			var rt roadtrip.CSV
			err := rt.LoadFile(fmt.Sprintf("data/CSV/%s", filename))
			if err != nil {
				v.Logrus().WithFields(log.Fields{
					"filename": filename,
					"error":    err,
				}).Error("Error loading vehicle")
				break
			}

			err = SyncGasRecords(v, rt)
			if err != nil {
				v.Logrus().WithError(err).Error("Error synching fuel records")
				break
			}
		}
	}
}
