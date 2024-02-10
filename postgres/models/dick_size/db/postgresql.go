package db

import (
	. "DickSizeBot/logger"
	"DickSizeBot/postgres"
	models "DickSizeBot/postgres/models/dick_size"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type repo struct {
	client postgres.Client
}

func (r *repo) GetDuelsStat(ctx context.Context, chatId int64) []map[string]string {
	query := `
select winner, ud.fname , ud.username , ud.lname , count(winner) as "wins"
from postgres.public.duels d 
inner join postgres.public.user_data ud on d.winner = ud.user_id 
where d.chat_id = %d
group by winner, ud.fname , ud.username , ud.lname 
order by count(winner) desc 
`

	var result []map[string]string

	query = fmt.Sprintf(query, chatId)

	Log.Debugf("query GetDuelsStat: %v", query)

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		Log.Errorf("query nothing returned")
		return result
	} else {
		for rows.Next() {
			var fname, lname, username string
			var winnerId, wins int
			err := rows.Scan(&winnerId, &fname, &username, &lname, &wins)
			if err != nil {
				Log.Printf("Error while scanning duels wins stat: %v", err)
				return result
			}

			oneRow := make(map[string]string)
			oneRow["fname"] = fname
			oneRow["lname"] = lname
			oneRow["username"] = username
			oneRow["wins"] = strconv.Itoa(wins)

			result = append(result, oneRow)
		}
	}
	return result
}

func (r *repo) GetLastDuelByUserId(ctx context.Context, userId int64, chatId int64) (time.Time, error) {
	query := `
	select duel_time 
	from postgres.public.duels duels
	where caller_user_id = %d and chat_id = %d
	order by duel_time desc
	limit 1
`
	query = fmt.Sprintf(query, userId, chatId)

	Log.Debugf("Exec query GetLastDuelByUserId: %s", query)

	var duelTime time.Time

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		Log.Errorf("Error while exec query GetLastDuelByUserId: %s", query)
		return duelTime, err
	} else {
		for rows.Next() {
			err := rows.Scan(
				//&duel.DuelId,
				//&duel.CallerUserId,
				//&duel.CallerRoll,
				//&duel.CalledUserId,
				//&duel.CalledRoll,
				//&duel.ChatID,
				//&duel.Bet,
				//&duel.Winner,
				//&duel.DuelDate,
				&duelTime,
			)

			if err != nil {
				Log.Errorf("Error while parsing query result: %v", err)
				return duelTime, err
			}
		}
	}

	return duelTime, nil
}

func (r *repo) IncreaceLastDickSize(ctx context.Context, dickSizeID int, bet int) {
	query := `
	update postgres.public.dick_size 
	set dick_size = (select dick_size + (%v) as "updated_dick"
				 from dick_size 
				 where id = %v
				)
	where id  = %v
	returning id
`

	query = fmt.Sprintf(query, bet, dickSizeID, dickSizeID)
	Log.Debugf("Exec IncreaceLastDickSize query: %v", query)

	var quer interface{}

	err := r.client.QueryRow(ctx, query).Scan(quer)
	if err != nil {
		Log.Errorf("Error while exec IncreaceLastDickSize query: %v", err)
	}
}

func (r *repo) GetUserData(ctx context.Context, userId int64) (user models.UserCredentials) {
	query := `
select user_id, fname, username, lname
from public.user_data
where user_id = $1
`

	rows, err := r.client.Query(ctx, query, userId)
	if err != nil {
		Log.Errorf("Error while exec GetUserData query: %v", err)
		return user
	} else {
		counter := 0
		for rows.Next() {
			counter++
			if counter > 1 {
				Log.Errorf("GetUserData query returned more than one profile")
				break
			}

			err = rows.Scan(
				&user.UserId,
				&user.Fname,
				&user.Username,
				&user.Lname,
			)

			if err != nil {
				Log.Errorf("Error while parsing query result: %v", err)
			}
		}
	}
	return user
}

func (r *repo) InsertDuelData(ctx context.Context, duel models.Duel) int {
	query := `
insert into public.duels  (caller_user_id , caller_roll , called_user_id , called_roll , chat_id , bet , winner, duel_time) values
(%d, %d, %d , %d , %d,  %d , %d , CURRENT_TIMESTAMP) returning duel_id
`

	query = fmt.Sprintf(query, duel.CallerUserId, duel.CallerRoll, duel.CalledUserId, duel.CalledRoll, duel.ChatID, duel.Bet, duel.Winner)

	Log.Debugf("Query: %s", query)

	var insertedId int

	err := r.client.QueryRow(ctx, query).Scan(&insertedId)
	if err != nil {
		Log.Errorf("Error while exec InsertDuelData query: %v", err)
		return 0
	}
	return insertedId
}

func (r *repo) CreateOrUpdateUser(ctx context.Context, user_id int64, fname, lname, username string, chat_id int64) int {
	query := `
	insert into postgres.public.user_data (user_id, fname, username, lname) values 
	($1, $2, $3, $4)
	on conflict (user_id) do update 
		set fname = excluded.fname,
			username = excluded.username,
			lname  = excluded.lname;
`

	Log.Debugf("Query: \n%s", query)
	var insertedId int

	err := r.client.QueryRow(ctx, query, user_id, fname, username, lname).Scan(&insertedId)
	if err != nil {
		Log.Errorf("Error while exec CreateOrUpdateUser query: %v", err)
		return 0
	}
	return insertedId
}

// TODO изменить архитектуру бд: добавить отдельную таблицу с именами юзеров и брать их оттуда. Колонки: user_id, fname, username, lname, chat_id
// TODO Перед каждым запросом проверять существование сначала id, потом всех данных, если id есть, а данные не совпадают, то удалить старые и добавить новые

func (r *repo) GetAllCredentials(ctx context.Context, chatId int64) []models.UserCredentials {
	query := `
	select distinct ud.user_id , ud.fname , ud.username , ud.lname 
	from postgres.public.user_data ud 
	inner join postgres.public.dick_size ds on ud.user_id = ds.user_id 
	where ds.chat_id = %d
	`

	query = fmt.Sprintf(query, chatId)

	Log.Debugf("GetAllCredentials query: %s", query)

	var usersData []models.UserCredentials

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		Log.Errorf("SQL error while exec SelectOnlyTodaysMeasures: %s", err.Error())
	} else {
		indx := 0
		for rows.Next() {
			usersData = append(usersData, models.UserCredentials{})
			err := rows.Scan(
				&usersData[indx].UserId,
				&usersData[indx].Fname,
				&usersData[indx].Username,
				&usersData[indx].Lname,
			)
			if err != nil {
				Log.Errorf("Failed to parse row: %d", indx+1)
			}
			indx++
		}
	}
	return usersData
}

func (r *repo) SelectOnlyTodaysMeasures(ctx context.Context, chatId int64) ([]models.DickSize, error) {
	query := `
select ds.id , ds.user_id , ud.fname , ud.username , ud.lname , ds.dick_size , ds.measure_date , ds.chat_id , ds.is_group 
from postgres.public.dick_size ds
inner join postgres.public.user_data ud on ds.user_id = ud.user_id
where measure_date in (select max(measure_date) as measure_date
						from postgres.public.dick_size ds_in
						inner join postgres.public.user_data ud on ds_in.user_id = ud.user_id 
						where ds.chat_id = %d
						group by ds.user_id , ud.fname , ud.lname , ud.username)
order by dick_size desc
`

	query = fmt.Sprintf(query, chatId)

	Log.Debugf("SelectOnlyTodaysMeasures query: %s", query)

	var dicks []models.DickSize

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		Log.Errorf("SQL error while exec SelectOnlyTodaysMeasures: %s", err.Error())
	} else {
		indx := 0
		for rows.Next() {
			dicks = append(dicks, models.DickSize{})
			err := rows.Scan(&dicks[indx].Id,
				&dicks[indx].UsedId,
				&dicks[indx].Fname,
				&dicks[indx].Username,
				&dicks[indx].Lname,
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
	query := `
	insert into postgres.public.user_data (user_id, fname, username, lname) values 
	($1, $2, $3, $4)
	on conflict (user_id) do update 
		set fname = excluded.fname,
			username = excluded.username,
			lname  = excluded.lname,
			chat_id = excluded.chat_id
	returning user_id;
`
	var insertedId int

	err := r.client.QueryRow(ctx, query, user_id, fname, username, lname).Scan(&insertedId)
	if err != nil {
		Log.Errorf("Error while exec CreateOrUpdateUser query: %v", err)

	}

	query = `insert into public.dick_size (user_id, dick_size, measure_date, chat_id, is_group)
			values ($1,  $2, CURRENT_TIMESTAMP, $3, $4)
			returning id`
	Log.Debugf("InsertSize func, query %s", query)

	var id int
	err = r.client.QueryRow(ctx, query, user_id, dick_size, chat_id, is_group).Scan(&id)
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
	query := `
	select ds.id, ds.user_id , ud.fname , ud.username , ud.lname , ds.dick_size , ds.measure_date , ds.chat_id , ds.is_group 
	from postgres.public.dick_size ds 
	inner join postgres.public.user_data ud on ud.user_id = ds.user_id 
	where ds.user_id = $1 and ds.chat_id = $2
	order by ds.measure_date desc 
	limit 1
`

	Log.Debugf("GetLastMeasureByUserInThisChat func, query %s", query)

	model := models.DickSize{}

	rows, err := r.client.Query(ctx, query, user_id, chatId)
	if err != nil {
		return model, err
	} else {
		for rows.Next() {
			err := rows.Scan(
				&model.Id,
				&model.UsedId,
				&model.Fname,
				&model.Username,
				&model.Lname,
				&model.Dick_size,
				&model.Measure_date,
				&model.Chat_id,
				&model.Is_group)
			if err != nil {
				return model, err
			}
		}
	}

	return model, nil
}

func (r *repo) GetUserAllSizesByChatId(ctx context.Context, chatId int64) ([]map[string]string, error) {
	query := `
	select avg(ds.dick_size) as "average", ud.fname, ud.lname , ud.username, ud.user_id  
	from postgres.public.dick_size ds 
	inner join postgres.public.user_data ud on ud.user_id = ds.user_id 
	where ds.chat_id = %d
	group by ud.fname, ud.lname , ud.username, ud.user_id 
	order by "average" desc 
`

	query = fmt.Sprintf(query, chatId)

	var result []map[string]string

	Log.Debugf(query)

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return result, err
	} else {
		for rows.Next() {
			var fname, lname, username, average string
			var userId int
			err := rows.Scan(&average, &fname, &lname, &username, &userId)
			if err != nil {
				Log.Errorf("Error while scanning average: %v", err)
				return result, err
			}

			oneRow := make(map[string]string)
			oneRow["fname"] = fname
			oneRow["lname"] = lname
			oneRow["username"] = username

			averFloat, _ := strconv.ParseFloat(average, 64)
			oneRow["average"] = strconv.Itoa(int(averFloat))

			result = append(result, oneRow)
		}
	}
	return result, nil
}

func initRepo(client postgres.Client) error {
	userDataTableInitQuery := `
	CREATE TABLE public.user_data (
		user_id int4 NOT NULL,
		fname varchar(100) NULL,
		lname varchar(100) NULL,
		username varchar(100) NULL,
		chat_id int8 NULL,
		CONSTRAINT user_data_pkey PRIMARY KEY (user_id)
	);
	`
	_, err := client.Exec(context.TODO(), userDataTableInitQuery)
	if err != nil {
		return err
	}

	DuelsTableInitQuery := `
	CREATE table if not exists duels (
		duel_id SERIAL PRIMARY key,
		caller_user_id INT NOT null,
		caller_roll	 int not null,
		called_user_id int not null,
		called_roll int not null,
		chat_id bigint not null,
		bet int not null,
		winner int,
		duel_time TIMESTAMP
	);
	`
	_, err = client.Exec(context.TODO(), DuelsTableInitQuery)
	if err != nil {
		return err
	}

	BotTableTableInitQuery := `
	CREATE table if not exists dick_size (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		fname VARCHAR(50),
		lname VARCHAR(50),
		username VARCHAR(50),
		dick_size BIGINT,
		measure_date TIMESTAMP,
		chat_id BIGINT,
		is_group BOOL
	);
	`
	_, err = client.Exec(context.TODO(), BotTableTableInitQuery)
	if err != nil {
		return err
	}

	return nil
}

func NewRepo(client postgres.Client) (models.Repository, error) {
	err := initRepo(client)
	if err != nil {
		return nil, err
	}
	return &repo{
		client: client,
	}, nil
}
