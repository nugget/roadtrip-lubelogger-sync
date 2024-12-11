package lubelogger

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	api_uri       string
	authorization string
)

func Init(uri, auth string) {
	log.Trace("Initializing logrus api")
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

	log.WithFields(log.Fields{
		"vehicleID": vehicleID,
		"count":     len(response.Records),
	}).Info("Loaded LubeLogger GasRecords")

	return response, nil
}

func AddGasRecord(vehicleID int, gr GasRecord) (PostResponse, error) {
	requestBody := gr.URLValues()

	log.WithFields(log.Fields{
		"vehicleID":   vehicleID,
		"gr":          gr,
		"requestBody": requestBody.Encode(),
	}).Trace("AddRecord()")

	// fmt.Printf("%+v\n", requestBody.Encode())

	endpoint := fmt.Sprintf("vehicle/gasrecords/add?vehicleID=%d", vehicleID)

	response, err := APIPostForm(endpoint, requestBody)
	if err != nil {
		gr.Logrus().WithFields(log.Fields{
			"requestBody": requestBody.Encode(),
			"vehicleID":   vehicleID,
		}).Debug("Request Trace")

		return response, fmt.Errorf("AddGasRecord: %w", err)
	}

	return response, nil
}
