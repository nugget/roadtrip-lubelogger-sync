package roadtrip

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	Raw                string
}

func (rt *CSV) LoadFile(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	rt.Filename = filename
	rt.Raw = string(buf[:])

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
	//	return fmt.Errorf("Vehicle: %w", err)
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

func (rt *CSV) Section(header string) string {
	var (
		target bool = false
		outbuf string
	)

	scanner := bufio.NewScanner(strings.NewReader(rt.Raw))

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			target = false
		}

		if target {
			outbuf += line
			outbuf += "\n"
		}

		if line == header {
			target = true
		}
	}

	return outbuf
}

func (rt *CSV) Moo() {
	var s []Fuel

	data := []byte(`Odometer (mi),Trip Distance,Date,Fill Amount,Fill Units,Price per Unit,Total Price,Partial Fill,MPG,Note,Octane,Location,Payment,Conditions,Reset,Categories,Flags,Currency Code,Currency Rate,Latitude,Longitude,ID,Trip Comp Fuel Economy,Trip Comp Avg. Speed,Trip Comp Temperature,Trip Comp Drive Time,Tank Number
61287,227,"2024-11-10 19:07",13.767,Gal,3.6987,50.92,,16.4887,,"93 Octane","Shell","Apple Pay",,,,0,,1,29.677433,-98.059817,369,,,26,,0
61505,218,"2024-12-5 9:04",13.622,Gal,4.0993,55.84,,16.0035,"Multi

Line

Note","93 Octane","Tri-Star","Apple Pay",,,,0,,1,29.856987,-97.952196,371,,,13,,0

`)

	fmt.Println(data)
	fmt.Println(string(data))

	fmt.Println("Running:")

	result, err := csvlib.Unmarshal(data, &s)

	fmt.Printf("result: %+v\n", result)
	fmt.Printf("err: %+v\n", err)

	fmt.Printf("%+v\n", *result)
	for _, u := range s {
		fmt.Printf("%+v\n", u)
	}
}

func (rt *CSV) Parse(section string, target interface{}) error {
	if _, err := csvlib.Unmarshal([]byte(rt.Section(section)), target); err != nil {
		return err
	}

	return nil
}
