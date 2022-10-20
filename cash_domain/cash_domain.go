package cash_domain

import "time"

// Интерфейс объекта
type ICash interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
	Del(key string) bool
	DelAll() []string
	Count() int
}
