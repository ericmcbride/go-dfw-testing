package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ericmcbride/go-dfw-testing/pkg/clients"
	"github.com/ericmcbride/go-dfw-testing/pkg/handlers"
	"github.com/ericmcbride/go-dfw-testing/pkg/harness"
	"github.com/ericmcbride/go-dfw-testing/pkg/models"
	"github.com/ericmcbride/go-dfw-testing/pkg/server"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	harness.Run(m)
}

func TestCarsHandlerPatch(t *testing.T) {
	req, err := http.NewRequest("PATCH", "/cars", nil)
	if err != nil {
		t.Errorf("Error while reading request JSON: %s", err)
	}
	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 405 {
		t.Errorf("Expected: %d, but got: %d", 405, rr.Code)
	}
}

func TestCarsMissingHeaders(t *testing.T) {
	req, err := http.NewRequest("PUT", "/cars", nil)
	if err != nil {
		t.Errorf("Error while reading request JSON: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 401 {
		t.Errorf("Expected: %d, but got: %d", 401, rr.Code)
	}
}
func TestCarsHandlerPut(t *testing.T) {
	req, err := http.NewRequest("PUT", "/cars", nil)
	if err != nil {
		t.Errorf("Error while reading request JSON: %s", err)
	}
	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 405 {
		t.Errorf("Expected: %d, but got: %d", 405, rr.Code)
	}
}

func TestCarsHandlerHead(t *testing.T) {
	req, err := http.NewRequest("HEAD", "/cars", nil)
	if err != nil {
		t.Errorf("Error while reading request JSON: %s", err)
	}
	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 405 {
		t.Errorf("Expected: %d, but got: %d", 405, rr.Code)
	}
}

func TestCarPostHandlerUnmarshalError(t *testing.T) {
	payload := []byte(``)

	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}
	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	_, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("Error while converting body to string")
	}

	if rr.Code != 400 {
		t.Errorf("Expected: %d, but got: %d", 400, rr.Code)
	}
}

func TestCarPostHandler(t *testing.T) {
	var got models.CarModel
	payload := []byte(`{"make": "Toyota", "model": "Camry", "color": "green", "year": 2005}`)

	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}
	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected: 200, but got: %d ", rr.Code)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("Error while converting body to string")
	}

	json.Unmarshal([]byte(body), &got)
	if got.Make != "Toyota" {
		t.Errorf("Expected: Toyota, got: %s", got.Make)
	}
	if got.Model != "Camry" {
		t.Errorf("Expected: Camry, got: %s", got.Model)
	}
	if got.Color != "green" {
		t.Errorf("Expected: green, got: %s", got.Color)
	}

	harness.Truncate()
}

func TestCarPostHandlerInvalidPostBodyError(t *testing.T) {
	payload := []byte(`{"foo": "bar"}`)

	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}
	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 422 {
		t.Errorf("Expected: %d, but got: %d", 422, rr.Code)
	}
}

func TestCarPostHandlerMissingModelError(t *testing.T) {
	payload := []byte(`{"make": "toyota", "color": "red", "year": 1999}`)

	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}

	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 422 {
		t.Errorf("Expected: %d, but got: %d", 422, rr.Code)
	}
}

func TestCarPostHandlerMissingMakeError(t *testing.T) {
	payload := []byte(`{"model": "Camry", "color": "green", "year": 2018 }`)

	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}

	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 422 {
		t.Errorf("Expected: %d, but got: %d", 422, rr.Code)
	}
}

func TestCarPostHandlerMissingYearError(t *testing.T) {
	payload := []byte(`{"model": "Camry", "color": "green", "make": "Toyota" }`)

	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}

	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 422 {
		t.Errorf("Expected: %d, but got: %d", 422, rr.Code)
	}
}

func TestCarPostHandlerMissingColorError(t *testing.T) {
	payload := []byte(`{"model": "Camry", "year": 1999, "make": "Toyota" }`)

	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}

	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 422 {
		t.Errorf("Expected: %d, but got: %d", 422, rr.Code)
	}
}

func TestValidateAuthorizationIdPositive(t *testing.T) {
	got := handlers.ValidateAuthId("1234")
	if got != nil {
		t.Errorf("Expected: %s, but got: %s ", "nil", got)
	}
}

func TestDeleteHandler(t *testing.T) {
	// Save a model to delete
	carId := uuid.NewV4().String()
	db, err := clients.NewDbConn()
	if err != nil {
		t.Fatalf("Couldnt connect to db")
	}

	defer clients.Close(&db)
	carModel := &models.CarModel{
		Id:    carId,
		Model: "Corolla",
		Make:  "Toyota",
		Color: "White",
		Year:  2018,
	}

	id, err := models.SaveCar(&db, carModel)
	if err != nil {
		t.Fatalf("Couldn't save model")
	}
	if id == "" {
		t.Errorf("Expected: %s, got: %s", carId, id)
	}

	carIdStr := fmt.Sprintf("/cars?car_id=%s", id)
	req, err := http.NewRequest("DELETE", carIdStr, nil)
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}

	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected: %d, but got: %d", 200, rr.Code)
	}
	harness.Truncate()
}

func TestDeleteHandlerMissingCarId(t *testing.T) {
	carIdStr := fmt.Sprintf("/cars?car_id=")
	req, err := http.NewRequest("DELETE", carIdStr, nil)
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}

	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 400 {
		t.Errorf("Expected: %d, but got: %d", 400, rr.Code)
	}
	harness.Truncate()
}

func TestGetHandler(t *testing.T) {
	// Save a model to delete
	carId := uuid.NewV4().String()
	db, err := clients.NewDbConn()
	if err != nil {
		t.Fatalf("Couldnt connect to db")
	}

	defer clients.Close(&db)
	carModel := &models.CarModel{
		Id:    carId,
		Model: "Corolla",
		Make:  "Toyota",
		Color: "White",
		Year:  2018,
	}

	id, err := models.SaveCar(&db, carModel)
	if err != nil {
		t.Fatalf("Couldn't save model")
	}
	if id == "" {
		t.Errorf("Expected: %s, got: %s", carId, id)
	}

	carIdStr := fmt.Sprintf("/cars?car_id=%s", id)
	req, err := http.NewRequest("GET", carIdStr, nil)
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}

	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected: %d, but got: %d", 200, rr.Code)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Couldnt read response body")
	}
	var got models.CarModel

	json.Unmarshal([]byte(body), &got)

	if strings.EqualFold(got.Model, "Corolla") {
		t.Fatalf("Expected: Corolla, got: %s", got.Model)
	}
	if strings.EqualFold(got.Make, "Toyota") {
		t.Fatalf("Expected: Toyota, got: %s", got.Make)
	}
	if strings.EqualFold(got.Color, "White") {
		t.Fatalf("Expected: White, got %s", got.Color)
	}

	harness.Truncate()
}

func TestGetCarMissingCarId(t *testing.T) {
	carIdStr := fmt.Sprintf("/cars?car_id=")
	req, err := http.NewRequest("GET", carIdStr, nil)
	if err != nil {
		t.Errorf("Error while reading request payload: %s", err)
	}

	req.Header.Set("X-CARS-ID", "1234")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := server.New()
	handler.ServeHTTP(rr, req)

	if rr.Code != 400 {
		t.Errorf("Expected: %d, but got: %d", 400, rr.Code)
	}
}
