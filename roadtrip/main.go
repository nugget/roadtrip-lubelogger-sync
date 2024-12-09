package roadtrip

import (
	"bytes"
	"fmt"
	"os"

	"github.com/tiendc/go-csvlib"

	log "github.com/sirupsen/logrus"
)

type CSV struct {
	Delimiters         string
	Version            int
	Language           string
	Filename           string
	Vehicle            []Vehicle
	FuelRecords        []Fuel
	MaintenanceRecords []Maintenance
	RoadTrips          []RoadTrip
	TireLogs           []Tire
	Valuations         []Valuation
	Raw                []byte
}

func (rt *CSV) LoadFile(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	rt.Filename = filename
	rt.Raw = buf

	if err := rt.Parse("FUEL RECORDS", &rt.FuelRecords); err != nil {
		return fmt.Errorf("FuelRecords: %w", err)
	}

	if err := rt.Parse("MAINTENANCE RECORDS", &rt.MaintenanceRecords); err != nil {
		return fmt.Errorf("MaintenanceRecords: %w", err)
	}

	if err := rt.Parse("ROAD TRIPS", &rt.RoadTrips); err != nil {
		return fmt.Errorf("RoadTrips: %w", err)
	}

	// if err := rt.Parse("VEHICLE", &rt.Vehicle); err != nil {
	// 	return fmt.Errorf("Vehicle: %w", err)
	// }

	if err := rt.Parse("TIRE LOG", &rt.TireLogs); err != nil {
		return fmt.Errorf("TireLogs: %w", err)
	}

	if err := rt.Parse("VALUATIONS", &rt.Valuations); err != nil {
		return fmt.Errorf("Valuations: %w", err)
	}

	log.WithFields(log.Fields{
		"filename":          rt.Filename,
		"bytes":             len(buf),
		"lines":             len(rt.Raw),
		"vehicleRecords":    len(rt.Vehicle),
		"fuelRecords":       len(rt.FuelRecords),
		"mainteanceRecords": len(rt.MaintenanceRecords),
		"roadTrips":         len(rt.RoadTrips),
		"tireLogs":          len(rt.TireLogs),
		"valuations":        len(rt.Valuations),
	}).Info("Loaded Road Trip CSV")

	return nil
}

func (rt *CSV) Section(sectionHeader string) (outbuf []byte) {
	sectionStart := make(map[string]int)

	for index, element := range HEADERS {
		i := bytes.Index(rt.Raw, []byte(HEADERS[index]))
		sectionStart[element] = i
		fmt.Printf("Section %s starts at position %d\n", element, i)
	}

	startPosition := sectionStart[sectionHeader]
	endPosition := len(rt.Raw)

	for _, e := range sectionStart {
		if e > startPosition && e < endPosition {
			endPosition = e - 1
		}
	}

	fmt.Printf("Section %s: %6d - %6d\n", sectionHeader, startPosition, endPosition)

	// Don't include the section header line in the outbuf
	startPosition = startPosition + len(sectionHeader) + 1

	outbuf = rt.Raw[startPosition:endPosition]

	fmt.Println(string(outbuf))

	return
}

func (rt *CSV) Parse(sectionHeader string, target interface{}) error {
	if _, err := csvlib.Unmarshal(rt.Section(sectionHeader), target); err != nil {
		return err
	}

	return nil
}
