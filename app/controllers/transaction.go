package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
			err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Description)
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
func (c Controller) AddTransaction(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var transaction models.Transaction
		var errorObj models.Error
		json.NewDecoder(r.Body).Decode(&transaction)
		if transaction.UserID < 0 {
			errorObj.Message = "\"UserId\" が指定されていません"
			utils.Respond(w, http.StatusBadRequest, errorObj)
			return
		}
		if transaction.Description == "" {
			errorObj.Message = "\"Description\" が指定されていません"
			utils.Respond(w, http.StatusBadRequest, errorObj)
			return
		}
		if transaction.Amount < 0 {
			errorObj.Message = "\"Amount\" が指定されていません"
			utils.Respond(w, http.StatusBadRequest, errorObj)
			return
		}
		rows, err := db.Query("select * from transactions where user_id=?", transaction.UserID)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			errorObj.Message = "SQLエラーです"
			utils.Respond(w, http.StatusInternalServerError, errorObj)
			return
		}
		defer rows.Close()
		totalAmount := 0
		amount := transaction.Amount
		for rows.Next() {
			err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Description)
			if err != nil {
				log.Println(err)
				errorObj.Message = "Server error"
				return
			}
			totalAmount += amount
		}
		if totalAmount+transaction.Amount > 1000 {
			log.Println(err)
			errorObj.Message = "amountが1000以上です。登録はできません。"
			utils.Respond(w, http.StatusPaymentRequired, errorObj)
			return
		}
		insert, err := db.Prepare("INSERT INTO transactions (user_id, description, amount) values(?,?,?)")
		if err != nil {
			log.Println(err)
			errorObj.Message = "transactionの準備ができませんでした。"
			utils.Respond(w, http.StatusInternalServerError, errorObj)
			return
		}
		defer insert.Close()
		result, err := insert.Exec(transaction.UserID, transaction.Description, transaction.Amount)
		if err != nil {
			log.Println(err)
			errorObj.Message = "transactionの実行ができませんでした。"
			utils.Respond(w, http.StatusInternalServerError, errorObj)
			return
		}
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			log.Println(err)
			errorObj.Message = "transactionのIDが取得できませんでした。"
			utils.Respond(w, http.StatusInternalServerError, errorObj)
			return
		}
		log.Println("lastInsertID=" + fmt.Sprint(lastInsertID))
		errorObj.Message = "作成しました"
		utils.Respond(w, http.StatusCreated, errorObj)
	}
}
