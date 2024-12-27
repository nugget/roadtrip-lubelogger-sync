package lubelogger

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
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
	PurchasePrice         float64      `json:"purchasePrice"`
	SoldPrice             float64      `json:"soldPrice"`
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

// LogValue is the handler for [log.slog] to emit structured output for a
// [Vehicle] object when logging.
func (v *Vehicle) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("year", v.Year),
		slog.String("make", v.Make),
		slog.String("model", v.Model),
	)
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

// LogValue is the handler for [log.slog] to emit structured output for a
// [GasRecord] object when logging.
func (gr *GasRecord) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("date", gr.Date),
		slog.String("Odometer", gr.Odometer),
	)
}

func (gr *GasRecord) Comparator() (comparator string) {
	if gr.Odometer == "" {
		return comparator
	}

	i, err := strconv.Atoi(gr.Odometer)
	if err != nil {
		slog.Error("Unable to parse Odometer value",
			"error", err,
			"gasRecord", gr,
		)
	}

	comparator = fmt.Sprintf("%07d", i)

	slog.Debug("Calculated GasRecord comparator",
		"gr", gr,
		"comparator", comparator,
	)

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
