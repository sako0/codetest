package controllers

// https://www.sunapro.com/mysql-lock/

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"codetest-docker/app/models"
	"codetest-docker/app/utils"
)

func (c Controller) AddTransaction(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var transaction models.Transaction
		var errorObj models.Error

		tx, err := db.Begin() // トランザクション開始
		if err != nil {
			errorObj.Message = "トランザクション開始できませんでした。"
			utils.Respond(w, http.StatusBadRequest, errorObj)
		}

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
		rows, err := tx.Query("select * from transactions where user_id=?", transaction.UserID)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			errorObj.Message = "SQLエラーです"
			utils.Respond(w, http.StatusInternalServerError, errorObj)
			return
		}
		defer rows.Close()

		totalAmount := 0
		var lastId int = 0
		for rows.Next() {
			var id int
			var user_id int
			var amount int
			var description string
			// カラムを変数に格納していく https://golang.shop/post/go-databasesql-04-retrieving-ja/
			err := rows.Scan(&id, &user_id, &amount, &description)
			if err != nil {
				log.Println(err)
				errorObj.Message = "Server error"
				return
			}
			totalAmount += amount
			lastId = id
		}
		err = rows.Err()
		if err != nil {
			log.Println(err)
			errorObj.Message = "rows.Next()後のエラー"
			return
		}
		if totalAmount+transaction.Amount > 1000 {
			log.Println(err)
			errorObj.Message = "amountが1000以上です。登録はできません。"
			utils.Respond(w, http.StatusPaymentRequired, errorObj)
			return
		}
		// https://blog.suganoo.net/entry/2019/01/25/190200
		transaction.ID = lastId + 1
		log.Println(transaction.ID)
		result, err := tx.Exec("INSERT INTO transactions (id, user_id, description, amount) values(?, ?, ?, ?)", transaction.ID, transaction.UserID, transaction.Description, transaction.Amount)
		if err != nil {
			errorObj.Message = "execエラー"
			log.Println(err)
			utils.Respond(w, http.StatusPaymentRequired, errorObj)
			return
		}
		lastInsertID, insertErr := result.LastInsertId()
		if insertErr != nil {
			errorObj.Message = "lastInsertIdの取得エラーです"
			log.Println(err)
			utils.Respond(w, http.StatusPaymentRequired, errorObj)
			return
		}
		// https://sourjp.github.io/posts/go-db/
		if lastId+1 != int(lastInsertID) {
			tx.Rollback()
			log.Println(err)
			log.Println("lastInsertID=" + fmt.Sprint(lastInsertID) + ", totalAmount=" + fmt.Sprint(totalAmount) + ", last_id=" + fmt.Sprint(lastId))
			errorObj.Message = "同時に複数回の登録はできません。"
			utils.Respond(w, http.StatusPaymentRequired, errorObj)
			return
		} else {
			tx.Commit()
			log.Println("lastInsertID=" + fmt.Sprint(lastInsertID) + ", totalAmount=" + fmt.Sprint(totalAmount) + ", last_id=" + fmt.Sprint(lastId))
			errorObj.Message = "作成しました"
			utils.Respond(w, http.StatusCreated, errorObj)
		}
	}
}
