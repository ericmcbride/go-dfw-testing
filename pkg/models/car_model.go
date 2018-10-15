package models

import (
	"github.com/ericmcbride/go-dfw-testing/pkg/clients"
	_ "github.com/lib/pq"
)

type CarModel struct {
	Id    string `json:"id"`
	Model string `json:"model"`
	Make  string `json:"make"`
	Color string `json:"color"`
	Year  int    `json:"year"`
}

func SaveCar(db *clients.DBClient, car *CarModel) (string, error) {
	sqlStatement := `
		INSERT INTO cars (id, model, make, color, year)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	var id string

	err := db.Db.QueryRow(
		sqlStatement,
		car.Id,
		car.Model,
		car.Make,
		car.Color,
		car.Year,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func DeleteCar(db *clients.DBClient, carId string) error {
	sqlStatement := `
		DELETE FROM cars
		WHERE id = $1;
	`

	_, err := db.Db.Exec(sqlStatement, carId)
	if err != nil {
		return err
	}

	return nil
}

func GetCar(db *clients.DBClient, carId string) (CarModel, error) {
	sqlStatement := `
		SELECT id, model, make, color, year
		FROM "cars" WHERE id = $1;
	`
	var carModel CarModel

	err := db.Db.QueryRow(
		sqlStatement,
		carId,
	).Scan(
		&carModel.Id,
		&carModel.Model,
		&carModel.Make,
		&carModel.Color,
		&carModel.Year,
	)

	if err != nil {
		return CarModel{}, err
	}

	return carModel, nil

}
