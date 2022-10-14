package cash

import (
	"sync"
	"time"
)

type item struct {
	value      interface{}
	created    time.Time
	expiration int64
}

type cash struct {
	sync.RWMutex
	defaultExpiration time.Duration
	items             map[string]item
}

var singleCash *cash

// Создание объекта
func NewCash() interface{} {
	if singleCash == nil {
		singleCash = &cash{
			defaultExpiration: 1 * time.Second,
			items:             make(map[string]item),
		}
	}

	return singleCash
}

// func (obj *cash) Building(f factory.IFactory, data map[string]interface{}) interface{} {
// 	return obj
// }

func (obj *cash) Del(key string) bool {
	obj.RLock()
	defer obj.RUnlock()

	_, found := obj.items[key]
	if !found {
		return false
	}

	delete(obj.items, key)
	return true
}

// Установка значения в кеш duration = 1 * time.Second
func (obj *cash) Set(key string, value interface{}, duration time.Duration) {

	var expiration int64

	if duration == 0 {
		duration = obj.defaultExpiration
	}

	// Устанавливаем время истечения кеша
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	obj.Lock()

	defer obj.Unlock()

	obj.items[key] = item{
		value:      value,
		expiration: expiration,
		created:    time.Now(),
	}

	// Ищем элементы с истекшим временем жизни и удаляем из хранилища
	for k, i := range obj.items {
		if time.Now().UnixNano() > i.expiration && i.expiration > 0 {
			delete(obj.items, k)
		}
	}
}

// Получение значения из кеша
func (obj *cash) Get(key string) (interface{}, bool) {

	obj.RLock()

	defer obj.RUnlock()

	item, found := obj.items[key]

	// ключ не найден
	if !found {
		return nil, false
	}

	if item.expiration > 0 {

		// Если в момент запроса кеш устарел возвращаем nil
		if time.Now().UnixNano() > item.expiration {
			return nil, false
		}

	}

	return item.value, true
}

// Получение кол-ва элементов в кеше
func (obj *cash) Count() int {

	obj.RLock()

	defer obj.RUnlock()

	return len(obj.items)
}
