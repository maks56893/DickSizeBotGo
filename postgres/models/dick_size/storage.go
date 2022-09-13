package models

import (
	"context"
)

type Repository interface {
	InsertSize(ctx context.Context, user_id int64, fname, lname, username string, dick_size int, chat_id int64, is_group bool) (int, error)
	GetLastMeasureByUserInThisChat(ctx context.Context, user_id int64, chatId int64) (DickSize, error)
	GetUserAllSizesByChatId(ctx context.Context, chatId int64) ([]map[string]string, error)
	DeleteSizesByTime(ctx context.Context)
	//	CreateTableIfNotExists(ctx context.Context, chatId int64)
}
