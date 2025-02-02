package sql

import (
	"context"
	"fmt"
	"log"
	"os"
)

func ClearDatabaseInfo() {
	log.Println("Clearing database info")
	// language=PostgreSQL
	sqls := []string{
		"DELETE FROM answer;",
		"DELETE FROM availability;",
		"DELETE FROM schedule;",
		"DELETE FROM selection;",
		"DELETE FROM select_option;",
		"DELETE FROM category;",
		"DELETE FROM game;",
		"DELETE FROM submission;",
		"DELETE FROM question;",
		"DELETE FROM marathon;",
		"DELETE FROM moderator;",
		"DELETE FROM user_roles;",
		"DELETE FROM password_resets;",
		"DELETE FROM email_verification;",
		"DELETE FROM social_accounts;",
		"DELETE FROM users;",
	}

	for _, sql := range sqls {
		log.Printf("Running sql: %s\n", sql)
		db := GetConnection()
		_, err := db.Query(context.Background(), sql)
		CloseConnection(db)

		if err == nil {
			log.Printf("Running sql: '%s' success!\n", sql)
		} else {
			fmt.Fprintf(os.Stderr, "Failed execute sql '%s': %v\n", sql, err)
		}
	}

	log.Println("Clearing database info done!")
}
