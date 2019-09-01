package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ericmcbride/go-dfw-testing/pkg/clients"
	"github.com/ericmcbride/go-dfw-testing/pkg/logging"
	"github.com/ericmcbride/go-dfw-testing/pkg/models"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

type CarPostPayload struct {
	Make  string `json:"make"`
	Model string `json:"model"`
	Color string `json:"color"`
	Year  int    `json:"year"`
}

func CarsHandler(w http.ResponseWriter, r *http.Request) {
	log := logging.GetLog(r.Context())

	err := ValidateAuthId(r.Header.Get("X-CARS-ID"))
	if err != nil {
		log.Error("Unauthorized Auth Id: ", err)
		http.Error(w, "Unauthorized", 401)
		return
	}

	var (
		statusCode int
	)

	switch r.Method {
	case "POST":
		statusCode, err = PostCar(w, r)
	case "DELETE":
		statusCode, err = DeleteCar(w, r)
	case "GET":
		statusCode, err = GetCar(w, r)
	default:
		http.Error(w, "Invalid Request Method.", 405)
		return
	}

	if err != nil {
		log.Error(err)
		jsonErr := &logging.JsonError{
			Status:  http.StatusText(statusCode),
			Code:    strconv.Itoa(statusCode),
			Message: err.Error(),
		}
		logging.FormatError(r.Context(), w, statusCode, *jsonErr)
	}
}

func PostCar(w http.ResponseWriter, r *http.Request) (int, error) {
	var postPayload CarPostPayload

	ctx := r.Context()
	log := logging.GetLog(ctx)
	log.Info("PostCar: Processing Add Car endpoint...")

	// get db conn
	log.Debug("PostCar: Getting Database Connection...")
	db, err := clients.NewDbConn()
	if err != nil {
		return 500, err
	}
	defer clients.Close(&db)

	log.Debug("PostCar: Decoding request body...")
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&postPayload)
	if err != nil {
		log.WithError(err)
		return 400, err
	}

	log.Debug("PostCar: validating payload...")
	err = ValidateCarPayload(&postPayload)
	if err != nil {
		return 422, err
	}

	log.Debug("PostCar: building model from payload...")
	carId := uuid.NewV4().String()
	carModel := &models.CarModel{
		Id:    carId,
		Model: postPayload.Model,
		Make:  postPayload.Make,
		Color: postPayload.Color,
		Year:  postPayload.Year,
	}

	log.Debug("PostCar: Saving Car Model")
	carId, err = models.SaveCar(&db, carModel)
	if err != nil {
		return 500, err
	}

	// Return valid response if No err comes back from db
	saveCarJson, err := json.Marshal(carModel)
	if err != nil {
		return 400, err
	}

	// Send back response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(saveCarJson)

	return 200, nil
}

func DeleteCar(w http.ResponseWriter, r *http.Request) (int, error) {
	ctx := r.Context()
	log := logging.GetLog(ctx)
	log.Info("DeleteCar: Processing Delete Car endpoint...")

	// get db conn
	log.Debug("DeleteCar: Getting Database Connection...")
	db, err := clients.NewDbConn()
	if err != nil {
		return 500, err
	}
	defer clients.Close(&db)

	log.Debug("DeleteCar: Getting ID query param...")
	carId := r.URL.Query().Get("car_id")
	if carId == "" {
		return 400, errors.New("Need a Car ID to delete...")
	}

	log.Debug("DeleteCar: Deleting car from databse...")
	err = models.DeleteCar(&db, carId)
	if err != nil {
		return 500, err
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	return 200, nil
}

func GetCar(w http.ResponseWriter, r *http.Request) (int, error) {
	ctx := r.Context()
	log := logging.GetLog(ctx)
	log.Info("GetCar: Processing Get Cars endpoint...")

	// get db conn
	log.Debug("GetCar: Getting Database Connection...")
	db, err := clients.NewDbConn()
	if err != nil {
		return 500, err
	}
	defer clients.Close(&db)

	log.Debug("GetCar: Getting ID query param...")
	carId := r.URL.Query().Get("car_id")
	if carId == "" {
		return 400, errors.New("Need a car Id to GET...")
	}

	log.Debug("GetCar: Getting car from databse...")
	car, err := models.GetCar(&db, carId)
	if err != nil {
		return 500, err
	}

	// Return valid response if No err comes back from db
	carJson, err := json.Marshal(car)
	if err != nil {
		return 400, err
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(carJson)
	return 200, nil
}

func ValidateCarPayload(payload *CarPostPayload) error {
	if payload.Make == "" {
		return errors.New("Make must be included in the payload")
	}
	if payload.Model == "" {
		return errors.New("Model must be included in the payload")
	}
	if payload.Color == "" {
		return errors.New("Color must be included in the payload")
	}
	if payload.Year == 0 {
		return errors.New("Year must be included in the payload")
	}

	return nil
}

func ValidateAuthId(auth string) error {
	if auth != "1234" {
		return errors.New("Invalid Authorization Header ID")
	}
	return nil
}
