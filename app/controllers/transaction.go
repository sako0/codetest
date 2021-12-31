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
		// jsonからparamを受け取りバリデーションを行う
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
		// 同じuser_idのtransactionsを全て検索
		rows, err := tx.Query("select * from transactions where user_id=?", transaction.UserID)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			errorObj.Message = "SQLエラーです"
			utils.Respond(w, http.StatusInternalServerError, errorObj)
			return
		}
		defer rows.Close()
		// 同じuser_idのtransactionsのamountを足していく
		totalAmount := 0
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
		}
		err = rows.Err()
		if err != nil {
			log.Println(err)
			errorObj.Message = "rows.Next()後のエラー"
			return
		}
		// 同じuser_idのtransactionsのamountの合計値が規定値を超えていた場合エラーを吐くようにする
		if totalAmount+transaction.Amount > 1000 {
			log.Println(err)
			errorObj.Message = "amountが1000以上です。登録はできません。"
			utils.Respond(w, http.StatusPaymentRequired, errorObj)
			return
		}
		// transactionsの最終レコードのIDを取得
		var lastRowTransactionId int
		err = tx.QueryRow("select id from transactions ORDER BY id DESC LIMIT 1").Scan(&lastRowTransactionId)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			errorObj.Message = "SQLエラーです。"
			utils.Respond(w, http.StatusInternalServerError, errorObj)
			return
		} else if err == sql.ErrNoRows {
			lastRowTransactionId = 0
		}

		// レコードの追加(現時点で最新のレコードのIDを指定して追加) https://blog.suganoo.net/entry/2019/01/25/190200
		transaction.ID = lastRowTransactionId + 1
		log.Println(transaction.ID)
		result, err := tx.Exec("INSERT INTO transactions (id, user_id, description, amount) values(?, ?, ?, ?)", transaction.ID, transaction.UserID, transaction.Description, transaction.Amount)
		if err != nil {
			errorObj.Message = "INSERT時にエラーが発生しました。"
			log.Println(err)
			utils.Respond(w, http.StatusPaymentRequired, errorObj)
			return
		}
		lastInsertID, insertErr := result.LastInsertId() // INSERTした行のIDを取得する
		if insertErr != nil {
			errorObj.Message = "lastInsertIdの取得エラーです"
			log.Println(err)
			utils.Respond(w, http.StatusPaymentRequired, errorObj)
			return
		}
		// 現時点で最新のレコードのID+1のレコードだった場合のみ追加ができる https://sourjp.github.io/posts/go-db/
		if lastRowTransactionId+1 != int(lastInsertID) {
			tx.Rollback()
			log.Println(err)
			log.Println("lastInsertID=" + fmt.Sprint(lastInsertID) + ", totalAmount=" + fmt.Sprint(totalAmount) + ", last_id=" + fmt.Sprint(lastRowTransactionId))
			errorObj.Message = "レコードが同時に登録されました。再度実行してください。"
			utils.Respond(w, http.StatusPaymentRequired, errorObj)
			return
		} else {
			tx.Commit()
			log.Println("lastInsertID=" + fmt.Sprint(lastInsertID) + ", totalAmount=" + fmt.Sprint(totalAmount) + ", last_id=" + fmt.Sprint(lastRowTransactionId))
			errorObj.Message = "レコードを作成しました"
			utils.Respond(w, http.StatusCreated, errorObj)
		}
	}
}
