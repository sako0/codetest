package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"codetest-docker/app/models"
	"codetest-docker/app/utils"
)

func (c Controller) GetTransactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var transaction models.Transaction
		transactions := make([]models.Transaction, 0)

		rows, err := db.Query("select * from transactions")
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&transaction.ID, &transaction.UserId, &transaction.Amount, &transaction.Description)
			if err != nil {
				log.Println(err)
			}
			transactions = append(transactions, transaction)
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}
		utils.Respond(w, http.StatusOK, transactions)
	}
}
