package lubelogger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

var (
	api_uri       string
	authorization string
)

func Init(uri, auth string) {
	slog.Debug("Initializing LubeLogger API")
	api_uri = uri
	authorization = auth
}

func FormatDate(t time.Time) string {
	return t.Format("1/2/2006")
}

func Vehicles() (response []Vehicle, err error) {
	body, err := APIGet("vehicles")

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling json: %w", err)
	}

	return response, nil
}

func GasRecords(vehicleID int) (response VehicleGasRecords, err error) {
	body, err := APIGet(fmt.Sprintf("vehicle/gasrecords?vehicleID=%d", vehicleID))

	err = json.Unmarshal(body, &response.Records)
	if err != nil {
		return VehicleGasRecords{}, fmt.Errorf("unmarshalling json: %w", err)
	}

	slog.Info("Loaded LubeLogger GasRecords",
		"vehicleID", vehicleID,
		"count", len(response.Records),
	)

	return response, nil
}

func AddGasRecord(vehicleID int, gr GasRecord) (PostResponse, error) {
	requestBody := gr.URLValues()

	slog.Debug("AddRecord()",
		"vehicleID", vehicleID,
		"gr", gr,
		"requestBody", requestBody.Encode(),
	)

	// fmt.Printf("%+v\n", requestBody.Encode())

	endpoint := fmt.Sprintf("vehicle/gasrecords/add?vehicleID=%d", vehicleID)

	response, err := APIPostForm(endpoint, requestBody)
	if err != nil {
		slog.Debug("Request Debug",
			"gr", gr,
			"requestBody", requestBody.Encode(),
			"vehicleID", vehicleID,
		)

		return response, fmt.Errorf("AddGasRecord: %w", err)
	}

	return response, nil
}
