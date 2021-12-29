package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"codetest-docker/app/models"
	"codetest-docker/app/utils"
)

func (c Controller) GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		users := make([]models.User, 0)

		rows, err := db.Query("select * from users")
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&user.ID, &user.Name, &user.ApiKey)
			if err != nil {
				log.Println(err)
			}
			users = append(users, user)
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}
		utils.Respond(w, http.StatusOK, users)
	}
}
