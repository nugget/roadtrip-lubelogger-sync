package roadtrip

import (
	"bytes"
	"fmt"
	"os"
	"time"

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

	if true {
		// Remove erroneous header fields for VEHICLE section
		// per Darren Stone 9-Dec-2024 via email
		omitHeaders := []byte(",Tank 1 Type,Tank 2 Type,Tank 2 Units")
		rt.Raw = bytes.Replace(buf, omitHeaders, []byte{}, 1)
	} else {
		rt.Raw = buf
	}

	if err := rt.Parse("FUEL RECORDS", &rt.FuelRecords); err != nil {
		return fmt.Errorf("FuelRecords: %w", err)
	}

	if err := rt.Parse("MAINTENANCE RECORDS", &rt.MaintenanceRecords); err != nil {
		return fmt.Errorf("MaintenanceRecords: %w", err)
	}

	if err := rt.Parse("ROAD TRIPS", &rt.RoadTrips); err != nil {
		return fmt.Errorf("RoadTrips: %w", err)
	}

	if err := rt.Parse("VEHICLE", &rt.Vehicle); err != nil {
		return fmt.Errorf("Vehicle: %w", err)
	}

	if err := rt.Parse("TIRE LOG", &rt.TireLogs); err != nil {
		return fmt.Errorf("TireLogs: %w", err)
	}

	if err := rt.Parse("VALUATIONS", &rt.Valuations); err != nil {
		return fmt.Errorf("Valuations: %w", err)
	}

	log.WithFields(log.Fields{
		"filename":          rt.Filename,
		"bytes":             len(buf),
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
	log.WithFields(log.Fields{
		"sectionHeader": sectionHeader,
	}).Trace("Fetching Section from Raw")

	sectionStart := make(map[string]int)

	for index, element := range HEADERS {
		i := bytes.Index(rt.Raw, []byte(HEADERS[index]))
		sectionStart[element] = i
		log.WithFields(log.Fields{
			"element":      element,
			"sectionStart": i,
		}).Trace("Section Start detected")
	}

	startPosition := sectionStart[sectionHeader]
	endPosition := len(rt.Raw)

	for _, e := range sectionStart {
		if e > startPosition && e < endPosition {
			endPosition = e - 1
		}
	}

	// Don't include the section header line in the outbuf
	startPosition = startPosition + len(sectionHeader) + 1

	outbuf = rt.Raw[startPosition:endPosition]

	log.WithFields(log.Fields{
		"sectionHeader": sectionHeader,
		"startPosition": startPosition,
		"endPosition":   endPosition,
		"sectionBytes":  len(outbuf),
	}).Trace("Section Range calculated")

	return
}

func (rt *CSV) Parse(sectionHeader string, target interface{}) error {
	if _, err := csvlib.Unmarshal(rt.Section(sectionHeader), target); err != nil {
		return err
	}

	return nil
}

func ParseDate(dateString string) (t time.Time) {
	t, err := time.Parse("2006-1-2 15:04", dateString)
	if err != nil {
		t, err = time.Parse("2006-1-2", dateString)
		if err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"dateString": dateString,
			}).Debug("Can't parse Road Trip date string")
		}
	}

	return t
}
