package models

import (
	"context"
	"time"
)

type Repository interface {
	InsertSize(ctx context.Context, user_id int64, fname, lname, username string, dick_size int, chat_id int64, is_group bool) (int, error)
	GetLastMeasureByUserInThisChat(ctx context.Context, user_id int64, chatId int64) (DickSize, error)
	GetUserAllSizesByChatId(ctx context.Context, chatId int64) ([]map[string]string, error)
	DeleteSizesByTime(ctx context.Context)
	SelectOnlyTodaysMeasures(ctx context.Context, chatId int64) ([]DickSize, error)
	GetAllCredentials(ctx context.Context, chatId int64) []UserCredentials
	GetUserData(ctx context.Context, userId int64) (user UserCredentials)
	CreateOrUpdateUser(ctx context.Context, user_id int64, fname, lname, username string, chat_id int64) int
	InsertDuelData(ctx context.Context, duel Duel) int
	IncreaceLastDickSize(ctx context.Context, dickSizeId int, bet int)
	GetLastDuelByUserId(ctx context.Context, userId int64, chatId int64) (time.Time, error)
	GetDuelsStat(ctx context.Context, chatId int64) []map[string]string

	//	CreateTableIfNotExists(ctx context.Context, chatId int64)
}
