package sql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"oengus-timers/structs"
	"os"
	"time"
)

func FindMarathonsToOpenSubmissions() ([]structs.Marathon, error) {
	// language=PostgreSQL
	sql := buildMarathonQuery("submits_open = False AND NOW() >= submissions_start_date")

	db := GetConnection()
	defer CloseConnection(db)

	rows, err := db.Query(context.Background(), sql)

	if err != nil {
		return nil, err
	}

	foundMarathons := buildMarathonFromRows(rows)

	return foundMarathons, nil
}

func FindMarathonsToCloseSubmissions() ([]structs.Marathon, error) {
	// language=PostgreSQL
	sql := buildMarathonQuery("submits_open = True AND NOW() >= submissions_end_date")

	db := GetConnection()
	defer CloseConnection(db)

	rows, err := db.Query(context.Background(), sql)

	if err != nil {
		return nil, err
	}

	foundMarathons := buildMarathonFromRows(rows)

	return foundMarathons, nil
}

func OpenSubmission(marathonIds []string) {
	for _, marathonId := range marathonIds {
		log.Println("Opening submissions for marathon", marathonId)
	}
}

func CloseSubmission(marathonIds []string) {
	for _, marathonId := range marathonIds {
		log.Println("Closing submissions for marathon", marathonId)
	}
}

func buildMarathonQuery(wherePart string) string {
	// language=PostgreSQL
	return "SELECT id, name, start_date, end_date, submissions_start_date, submissions_end_date FROM marathon WHERE " + wherePart
}

func buildMarathonFromRows(rows pgx.Rows) []structs.Marathon {
	var foundMarathons []structs.Marathon
	var Id string
	var Name string
	var StartDate time.Time
	var EndDate time.Time
	var SubmissionsStartDate time.Time
	var SubmissionsEndDate time.Time

	for rows.Next() {
		// Scan is positional, not name based
		err2 := rows.Scan(&Id, &Name, &StartDate, &EndDate, &SubmissionsStartDate, &SubmissionsEndDate)

		if err2 != nil {
			fmt.Println("Error scanning rows", err2)
			continue
		}

		parsed := structs.Marathon{
			Id:                   Id,
			Name:                 Name,
			StartDate:            StartDate,
			EndDate:              EndDate,
			SubmissionsStartDate: SubmissionsStartDate,
			SubmissionsEndDate:   SubmissionsEndDate,
		}

		foundMarathons = append(foundMarathons, parsed)
	}

	return foundMarathons
}

func GetConnection() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func CloseConnection(db *pgx.Conn) {
	err := db.Close(context.Background())

	if err != nil {
		log.Println(err)
	}
}
