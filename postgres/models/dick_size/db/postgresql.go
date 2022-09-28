package db

import (
	. "DickSizeBot/logger"
	"DickSizeBot/postgres"
	models "DickSizeBot/postgres/models/dick_size"
	"context"
	"database/sql"
	"errors"
)

type repo struct {
	client postgres.Client
}

func (r *repo) SelectOnlyTodaysMeasures(ctx context.Context, chatId int64) ([]models.DickSize, error) {
	//	todayDate := time.Now()

	//	query := fmt.Sprintf("select * from postgres.public.dick_size ds where date(measure_date) = '%d-%d-%d' and chat_id = $1 order by dick_size desc ", todayDate.Year(), (todayDate.Month()), todayDate.Day())
	query := `
			select *
			from postgres.public.dick_size ds
			where measure_date in (select max(measure_date) as measure_date
									from postgres.public.dick_size ds_in
									where chat_id = $1
									group by user_id , fname , lname , username)
			order by measure_date desc`

	Log.Debugf("SelectOnlyTodaysMeasures query: %s", query)

	var dicks []models.DickSize

	rows, err := r.client.Query(ctx, query, chatId)
	if err != nil {
		Log.Errorf("SQL error while exec SelectOnlyTodaysMeasures: %s", err.Error())
	} else {
		indx := 0
		for rows.Next() {
			dicks = append(dicks, models.DickSize{})
			err := rows.Scan(&dicks[indx].Id,
				&dicks[indx].UsedId,
				&dicks[indx].Fname,
				&dicks[indx].Lname,
				&dicks[indx].Username,
				&dicks[indx].Dick_size,
				&dicks[indx].Measure_date,
				&dicks[indx].Chat_id,
				&dicks[indx].Is_group)
			if err != nil {
				//	return model, err
				Log.Errorf("Failed to parse row: %d", indx+1)
			}
			indx++
		}
	}
	return dicks, nil
}

func (r *repo) DeleteSizesByTime(ctx context.Context) {
	query := `delete from dick_size`

	Log.Debugf("DeleteSizesByTime query: %s", query)

	var count int

	err := r.client.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		Log.Debugf(err.Error())
	}
}

func (r *repo) InsertSize(ctx context.Context, user_id int64, fname, lname, username string, dick_size int, chat_id int64, is_group bool) (int, error) {
	query := `insert into public.dick_size (user_id, fname, lname, username, dick_size, measure_date, chat_id, is_group)
			values ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6, $7)
			returning id`
	Log.Debugf("InsertSize func, query %s", query)

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

	Log.Debugf("GetLastMeasureByUserInThisChat func, query %s", query)

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

	Log.Debugf("select avg(dick_size) as \"average\", fname, lname , username \n\t\t\t\tfrom public.dick_size ds \n\t\t\t\twhere chat_id = %v\n\t\t\t\tgroup by fname, lname, username \n\t\t\t\torder by \"average\" DESC", chatId)

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
