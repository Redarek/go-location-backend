package utils

import "time"

// Выполняет переданную fn функцию attemtps попыток через delay интервал
func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}

		return nil
	}

	return
}

// Возвращает OFFSET для PAGE и SIZE
func GetOffset(page, size int) (offset int) {
	return (page - 1) * size
}
