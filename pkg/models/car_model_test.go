package models_test

import (
	"github.com/ericmcbride/go-dfw-testing/pkg/clients"
	"github.com/ericmcbride/go-dfw-testing/pkg/harness"
	"github.com/ericmcbride/go-dfw-testing/pkg/models"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	harness.Run(m)
}

func TestSaveCar(t *testing.T) {
	carId := uuid.NewV4().String()

	db, err := clients.NewDbConn()
	if err != nil {
		t.Errorf("failed to establish a connection to the database")
	}
	defer clients.Close(&db)

	carModel := &models.CarModel{
		Id:    carId,
		Model: "prius-c",
		Make:  "corolla",
		Color: "black",
		Year:  2018,
	}

	id, err := models.SaveCar(&db, carModel)
	if err != nil {
		t.Fatalf("There was an error saving the carModel %s", err)
	}
	if id == "" {
		t.Errorf("Expected: %s, got: %s", carId, id)
	}

	lookupStmt := `
		SELECT id FROM "cars" WHERE id = $1;
	`
	var lookupId string
	err = db.Db.QueryRow(
		lookupStmt,
		id,
	).Scan(&lookupId)

	if lookupId != id {
		t.Fatalf("Did not save correctly, Expected: %s, got: %s", id, lookupId)
	}
	harness.Truncate()
}

func TestLookupCar(t *testing.T) {
	carId := uuid.NewV4().String()

	db, err := clients.NewDbConn()
	if err != nil {
		t.Errorf("failed to establish a connection to the database")
	}
	defer clients.Close(&db)

	carModel := &models.CarModel{
		Id:    carId,
		Model: "Corolla",
		Make:  "Toyota",
		Color: "white",
		Year:  2018,
	}

	id, err := models.SaveCar(&db, carModel)
	if err != nil {
		t.Fatalf("There was an error saving the carModel")
	}
	if id == "" {
		t.Errorf("Expected: %s, got: %s", carId, id)
	}

	got, err := models.GetCar(&db, carId)
	if err != nil {
		t.Fatalf("Failed to lookup the database info %v", err)
	}

	if got.Id != carId {
		t.Errorf("Expected: %s, got: %s", carId, got.Id)
	}
	if strings.EqualFold(got.Model, carModel.Model) {
		t.Errorf("Expected: %s, got: %s", carModel.Model, got.Model)
	}
	if strings.EqualFold(got.Make, carModel.Make) {
		t.Errorf("Expected: %s, got: %s", carModel.Make, got.Make)
	}

	harness.Truncate()
}

func TestDeleteCar(t *testing.T) {
	var got models.CarModel
	carId := uuid.NewV4().String()

	db, err := clients.NewDbConn()
	if err != nil {
		t.Errorf("failed to establish a connection to the database")
	}
	defer clients.Close(&db)

	carModel := &models.CarModel{
		Id:    carId,
		Model: "Corolla",
		Make:  "Toyota",
		Color: "white",
		Year:  2018,
	}

	id, err := models.SaveCar(&db, carModel)
	if err != nil {
		t.Fatalf("There was an error saving the carModel")
	}
	if id == "" {
		t.Errorf("Expected: %s, got: %s", carId, id)
	}

	err = models.DeleteCar(&db, carId)
	if err != nil {
		t.Fatalf("failed to delete car %v", err)
	}

	got, err = models.GetCar(&db, carId)
	if err == nil {
		t.Fatalf("Should be empty but got %v", got)
	}

	if len(got.Id) != 0 {
		t.Fatalf("Car should of been deleted %+v", got)
	}
	harness.Truncate()

}
