package cache

import (
	"sync"
	"time"
)

// Cache представляет интерфейс для работы с кешем.
type Cache interface {
	Set(key string, value interface{}, duration time.Duration) // Задает значение по ключу с указанной продолжительностью действия.
	Get(key string) (interface{}, bool)                        // Возвращает значение по ключу и флаг, указывающий на его наличие.
	Delete(key string) bool                                    // Удаляет значение по ключу и возвращает флаг, указывающий на успех операции.
	DeleteAll() []string                                       // Удаляет все значения из кеша и возвращает список удаленных ключей.
	Count() int                                                // Возвращает количество элементов в кеше.
}

type item struct {
	value   interface{} // Значение элемента.
	created time.Time   // Время создания элемента.
}

type cache struct {
	sync.RWMutex                      // Используется для блокировки доступа к кешу.
	defaultExpiration time.Duration   // Продолжительность действия элемента по умолчанию.
	items             map[string]item // Мапа для хранения элементов кеша.
}

// NewCache создает и возвращает новый экземпляр Cache.
func NewCache() Cache {
	c := &cache{
		defaultExpiration: 1 * time.Second,
		items:             make(map[string]item),
	}

	go c.startCleanupExpiredItems() // Запуск горутины для регулярной очистки истекших элементов.

	return c
}

// Set задает значение по ключу с указанной продолжительностью действия.
func (c *cache) Set(key string, value interface{}, duration time.Duration) {
	c.Lock()
	defer c.Unlock()

	if duration == 0 {
		duration = c.defaultExpiration
	}

	expiration := time.Now().Add(duration)
	c.items[key] = item{
		value:   value,
		created: expiration,
	}
}

// Get возвращает значение по ключу и флаг, указывающий на его наличие.
func (c *cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]
	if !found || c.isExpired(item) {
		return nil, false
	}

	return item.value, true
}

// Delete удаляет значение по ключу и возвращает флаг, указывающий на успех операции.
func (c *cache) Delete(key string) bool {
	c.Lock()
	defer c.Unlock()

	_, found := c.items[key]
	if !found {
		return false
	}

	delete(c.items, key)
	return true
}

// DeleteAll удаляет все значения из кеша и возвращает список удаленных ключей.
func (c *cache) DeleteAll() []string {
	c.Lock()
	defer c.Unlock()

	deletedKeys := make([]string, 0, len(c.items))
	for key := range c.items {
		deletedKeys = append(deletedKeys, key)
		delete(c.items, key)
	}

	return deletedKeys
}

// Count возвращает количество элементов в кеше.
func (c *cache) Count() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.items)
}

// startCleanupExpiredItems запускает горутину для регулярной очистки истекших элементов.
func (c *cache) startCleanupExpiredItems() {
	ticker := time.NewTicker(c.defaultExpiration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: // Ожидание истечения времени.
			c.cleanupExpiredItems()
		}
	}
}

// cleanupExpiredItems удаляет истекшие элементы из кеша.
func (c *cache) cleanupExpiredItems() {
	c.Lock()
	defer c.Unlock()

	for key, item := range c.items {
		if c.isExpired(item) {
			delete(c.items, key)
		}
	}
}

// isExpired проверяет, истек ли срок действия элемента.
func (c *cache) isExpired(item item) bool {
	return item.created.Before(time.Now())
}
