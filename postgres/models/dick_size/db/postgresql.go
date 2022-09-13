package db

import (
	"DickSizeBot/postgres"
	models "DickSizeBot/postgres/models/dick_size"
	"context"
	"database/sql"
	"errors"
	"log"
)

type repo struct {
	client postgres.Client
}

func (r *repo) DeleteSizesByTime(ctx context.Context) {
	query := `delete from dick_size`

	log.Printf("DeleteSizesByTime query: %s", query)

	var count int

	err := r.client.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		log.Println(err.Error())
	}
}

//func (r *repo) CreateTableIfNotExists(ctx context.Context, chatId int64) {
//	query := `create table if not exists dick_size_$1 (
//	id SERIAL  PRIMARY KEY,
//	user_id INT NOT NULL,
//	fname VARCHAR(50),
//	lname VARCHAR(50),
//	username VARCHAR(50),
//	dick_size BIGINT,
//	measure_date TIMESTAMP,
//	chat_id BIGINT,
//	is_group BOOL
//)`
//
//	log.Printf("InsertSize func, query %s", query)
//
//	var success interface{}
//
//	err := r.client.QueryRow(ctx, query, chatId).Scan(&success)
//	if err != nil {
//		log.Println(err)
//		if success == nil || success == "" {
//			log.Println("Table already presents!")
//		}
//	}
//}

func (r *repo) InsertSize(ctx context.Context, user_id int64, fname, lname, username string, dick_size int, chat_id int64, is_group bool) (int, error) {
	query := `insert into public.dick_size (user_id, fname, lname, username, dick_size, measure_date, chat_id, is_group)
			values ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6, $7)
			returning id`
	log.Printf("InsertSize func, query %s", query)

	var id int
	err := r.client.QueryRow(ctx, query, user_id, fname, lname, username, dick_size, chat_id, is_group).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, sql.ErrNoRows
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (r *repo) GetLastMeasureByUserInThisChat(ctx context.Context, user_id int64, chatId int64) (models.DickSize, error) {
	query := `select dick_size, measure_date
				from public.dick_size ds
				where user_id = $1 and chat_id = $2
				order by measure_date desc 
				limit 1`

	log.Printf("GetLastMeasureByUserInThisChat func, query %s", query)

	model := models.DickSize{}

	rows, err := r.client.Query(ctx, query, user_id, chatId)
	if err != nil {
		return model, err
	} else {
		for rows.Next() {
			err := rows.Scan(&model.Dick_size, &model.Measure_date)
			if err != nil {
				return model, err
			}
		}
	}

	return model, nil
}

func (r *repo) GetUserAllSizesByChatId(ctx context.Context, chatId int64) ([]map[string]string, error) {
	query := `select avg(dick_size) as "average", fname, lname , username 
				from public.dick_size ds 
				where chat_id = $1
				group by fname, lname, username 
				order by "average" DESC`

	var result []map[string]string

	log.Printf("select avg(dick_size) as \"average\", fname, lname , username \n\t\t\t\tfrom public.dick_size ds \n\t\t\t\twhere chat_id = %v\n\t\t\t\tgroup by fname, lname, username \n\t\t\t\torder by \"average\" DESC", chatId)

	rows, err := r.client.Query(ctx, query, chatId)
	if err != nil {
		return result, err
	} else {
		for rows.Next() {
			var fname, lname, username, average string
			err := rows.Scan(&average, &fname, &lname, &username)
			if err != nil {
				return result, err
			}

			oneRow := make(map[string]string)
			oneRow["fname"] = fname
			oneRow["lname"] = lname
			oneRow["username"] = username
			oneRow["average"] = average[:2]

			result = append(result, oneRow)
		}
	}
	return result, nil
}

func NewRepo(client postgres.Client) models.Repository {
	return &repo{
		client: client,
	}
}
