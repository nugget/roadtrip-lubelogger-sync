package lubelogger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

var (
	apiURI        string
	authorization string
	logger        *slog.Logger
)

func Init(uri, auth string, l *slog.Logger) {
	logger = l
	logger.Debug("Initializing LubeLogger API")
	apiURI = uri
	authorization = auth
}

func FormatDate(t time.Time) string {
	return t.Format("1/2/2006")
}

func Vehicles() ([]Vehicle, error) {
	var (
		response []Vehicle
	)

	body, err := APIGet("vehicles")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling json: %w", err)
	}

	return response, nil
}

func GasRecords(vehicleID int) (VehicleGasRecords, error) {
	var (
		response VehicleGasRecords
	)

	body, err := APIGet(fmt.Sprintf("vehicle/gasrecords?vehicleID=%d", vehicleID))
	if err != nil {
		return VehicleGasRecords{}, err
	}

	err = json.Unmarshal(body, &response.Records)
	if err != nil {
		return VehicleGasRecords{}, fmt.Errorf("unmarshalling json: %w", err)
	}

	logger.Info("Loaded LubeLogger GasRecords",
		"vehicleId", vehicleID,
		"count", len(response.Records),
	)

	return response, nil
}

func AddGasRecord(vehicleID int, gr GasRecord) (PostResponse, error) {
	requestBody := gr.URLValues()

	logger.Debug("AddRecord()",
		"vehicleId", vehicleID,
		"gr", gr,
		"requestBody", requestBody.Encode(),
	)

	// fmt.Printf("%+v\n", requestBody.Encode())

	endpoint := fmt.Sprintf("vehicle/gasrecords/add?vehicleID=%d", vehicleID)

	response, err := APIPostForm(endpoint, requestBody)
	if err != nil {
		logger.Debug("Request Debug",
			"gr", gr,
			"requestBody", requestBody.Encode(),
			"vehicleId", vehicleID,
		)

		return response, fmt.Errorf("AddGasRecord: %w", err)
	}

	return response, nil
}
