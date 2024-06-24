package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"oengus-timers/structs"
	"os"
	"strings"
	"time"
)

func FindMarathonsToOpenSubmissions() ([]structs.Marathon, error) {
	// Also check for end date here, we want to make sure that we do not open submissions when they need to be closed
	// language=PostgreSQL
	sql := buildMarathonQuery("submits_open = False AND NOW() >= submissions_start_date AND submissions_end_date >= NOW()")

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
	// Open submissions
	// Enable edits for submissions

	marathonIdsParsed := parseIdsToList(marathonIds)

	// language=PostgreSQL
	sql := fmt.Sprintf("UPDATE marathon SET submits_open = True, can_edit_submissions = True WHERE id IN (%s)", marathonIdsParsed)

	db := GetConnection()
	defer CloseConnection(db)

	_, err := db.Query(context.Background(), sql)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open marathon submissions: %v\n", err)
	}
}

func CloseSubmission(marathonIds []string) {
	// Close submissions
	// Keep edits enabled

	marathonIdsParsed := parseIdsToList(marathonIds)

	// language=PostgreSQL
	sql := fmt.Sprintf("UPDATE marathon SET submits_open = False WHERE id IN (%s)", marathonIdsParsed)

	db := GetConnection()
	defer CloseConnection(db)

	_, err := db.Query(context.Background(), sql)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open marathon submissions: %v\n", err)
	}
}

func parseIdsToList(ids []string) string {
	strs, _ := json.Marshal(ids)
	idsParsed := strings.Trim(string(strs), "[]")

	replacer := strings.NewReplacer("\"", "'")

	idsParsed = replacer.Replace(idsParsed)

	return idsParsed
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
