package lubelogger

import (
	"fmt"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type PostResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ExtraField struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	IsRequired bool   `json:"isRequired"`
}

type Vehicle struct {
	ID                    int          `json:"id"`
	ImageLocation         string       `json:"imageLocation"`
	Year                  int          `json:"year"`
	Make                  string       `json:"make"`
	Model                 string       `json:"model"`
	LicensePlate          string       `json:"licensePlate"`
	PurchaseDate          string       `json:"purchaseDate"`
	SoldDate              string       `json:"soldDate"`
	PurchasePrice         int          `json:"purchasePrice"`
	SoldPrice             int          `json:"soldPrice"`
	IsElectric            bool         `json:"isElectric"`
	IsDiesel              bool         `json:"isDiesel"`
	UseHours              bool         `json:"useHours"`
	OdometerOptional      bool         `json:"odometerOptional"`
	ExtraFields           []ExtraField `json:"extraFields"`
	Tags                  []string     `json:"tags"`
	HasOdometerAdjustment bool         `json:"hasOdometerAdjustment"`
	OdometerMultiplier    string       `json:"odometerMultiplier"`
	OdometerDifference    string       `json:"odometerDifference"`
	DashboardMetrics      []string     `json:"dashboardMetrics"`
	vehicleIdentifier     string       `json:"vehicleIdentifier"`
}

func (v *Vehicle) CSVFilename() string {
	for _, f := range v.ExtraFields {
		if f.Name == "Road Trip CSV" {
			return f.Value
		}
	}
	return ""
}

func (v *Vehicle) Logrus() *log.Entry {
	return log.WithFields(log.Fields{
		"year":  v.Year,
		"make":  v.Make,
		"model": v.Model,
	})
}

func (v *Vehicle) FindGasRecord(comparator string) (GasRecord, error) {
	records, err := GasRecords(v.ID)
	if err != nil {
		return GasRecord{}, err
	}

	return records.FindGasRecord(comparator)
}

type GasRecord struct {
	Date         string       `json:"date"`
	Odometer     string       `json:"odometer"`
	FuelConsumed string       `json:"fuelConsumed"`
	Cost         string       `json:"cost"`
	FuelEconomy  string       `json:"fuelEconomy"`
	IsFillToFull string       `json:"isFillToFull"`
	MissedFuelUp string       `json:"missedFuelUp"`
	Notes        string       `json:"notes"`
	Tags         string       `json:"tags"`
	ExtraFields  []ExtraField `json:"extraFields"`
}

func (gr *GasRecord) Logrus() *log.Entry {
	return log.WithFields(log.Fields{
		"date":     gr.Date,
		"odometer": gr.Odometer,
	})
}

func (gr *GasRecord) Comparator() (comparator string) {
	if gr.Odometer == "" {
		return comparator
	}

	i, err := strconv.Atoi(gr.Odometer)
	if err != nil {
		gr.Logrus().WithError(err).Error("Unable to parse Odometer value")
	}

	comparator = fmt.Sprintf("%07d", i)

	gr.Logrus().WithFields(log.Fields{
		"comparator": comparator,
	}).Trace("Calculated comparator")

	return comparator
}

func (gr *GasRecord) URLValues() url.Values {
	v := url.Values{}

	v.Add("Date", gr.Date)
	v.Set("Odometer", gr.Odometer)
	v.Set("FuelConsumed", gr.FuelConsumed)
	v.Set("Cost", gr.Cost)
	v.Set("FuelEconomy", gr.FuelEconomy)
	v.Set("IsFillToFull", gr.IsFillToFull)
	v.Set("MissedFuelUp", gr.MissedFuelUp)
	v.Set("Notes", gr.Notes)

	for i, ef := range gr.ExtraFields {
		v.Set(fmt.Sprintf("Extrafields[%d][Name]", i), ef.Name)
		v.Set(fmt.Sprintf("Extrafields[%d][Value]", i), ef.Value)
	}

	return v
}

type VehicleGasRecords struct {
	Records []GasRecord
}

func (r *VehicleGasRecords) FindGasRecord(comparator string) (GasRecord, error) {
	for _, gr := range r.Records {
		if comparator == gr.Comparator() {
			return gr, nil
		}
	}

	return GasRecord{}, nil
}
