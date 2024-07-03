package server

import (
	"main/internal/configs"
	"main/internal/pages"
	"net/http"
)

func HandleRequest() {
	http.HandleFunc("/", pages.Home_page)       // отслеживаем переход по URL (/ - переход на главную страницу)
	http.HandleFunc("/users", pages.Users_page) // ВАЖНО - в конце URL дописываем /, чтобы он корректно обрабатывался

	http.ListenAndServe(configs.Port, nil) // запускаем локальный сервер на порту 8080 (параметры: порт и настройки запуска)
}
