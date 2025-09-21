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
		"DELETE FROM answer WHERE id in (select id from answer);",
		"DELETE FROM availability WHERE submission_id in (select submission_id from availability);",
		"DELETE FROM schedule_line_runner WHERE schedule_line_id in (select schedule_line_id from schedule_line_runner);",
		"DELETE FROM schedule_line WHERE id in (select id from schedule_line);",
		"DELETE FROM schedule WHERE id in (select id from schedule);",
		"DELETE FROM selection WHERE id in (select id from selection);",
		"DELETE FROM select_option WHERE question_id in (select question_id from select_option);",
		"DELETE FROM opponent WHERE id in (select id from opponent);",
		"DELETE FROM category WHERE id in (select id from category);",
		"DELETE FROM game WHERE id in (select id from game);",
		"DELETE FROM submission WHERE id in (select id from submission);",
		"DELETE FROM question WHERE id in (select id from question);",
		"DELETE FROM moderator WHERE marathon_id in (select marathon_id from moderator);",
		"DELETE FROM marathon WHERE id in (select id from marathon);",
		"DELETE FROM user_roles WHERE user_id in (select user_id from user_roles);",
		"DELETE FROM password_resets WHERE token in (select token from password_resets);",
		"DELETE FROM email_verification WHERE user_id in (select user_id from email_verification);",
		"DELETE FROM social_accounts WHERE id in (select id from social_accounts);",
		"DELETE FROM users WHERE id in (select id from users);",
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
