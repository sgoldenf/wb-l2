package pattern

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
	Цепочка обязанностей — это поведенческий паттерн проектирования, который позволяет передавать запросы последовательно по цепочке обработчиков.
	Каждый последующий обработчик решает, может ли он обработать запрос сам и стоит ли передавать запрос дальше по цепи.

	Применимость:
	- Когда программа должна обрабатывать разнообразные запросы несколькими способами,
		но заранее неизвестно, какие конкретно запросы будут приходить и какие обработчики для них понадобятся
	- Когда важно, чтобы обработчики выполнялись один за другим в строгом порядке
	- Когда набор объектов, способных обработать запрос, должен задаваться динамически.

	Плюсы и минусы:
	+ Уменьшает зависимость между клиентом и обработчиками
	+ Реализует принцип единственной обязанности
	+ Реализует принцип открытости/закрытости
	- Запрос может остаться никем не обработанным

	Примеры использования на практике:
	Middleware - цепочка обработчиков запроса
	Например, пользователь делает запрос на контент, доступный премиум-подписчикам.
	Middleware в цепочке может проверить, авторизирован ли он, затем есть ли у него премиум-подписка, и только после этого вернуть ответ с контентом
	Если на каком-то из этапов выясняется, что запрос не проходит проверку, цепочка прервется и выведется страница с соответствующим сообщением
*/

// RequestInfo представляет информацию о клиенте (auth -- авторизация, premium -- наличие подписки).
type RequestInfo struct {
	auth    bool
	premium bool
}

// Handler - интерфейс обработчика.
type Handler interface {
	Handle(u *RequestInfo) string
}

// CheckAuthenticationHandler - конкретный обработчик для проверки авторизации.
type CheckAuthenticationHandler struct {
	Next Handler
}

// Handle - метод обработки события (проверка авторизации).
func (h *CheckAuthenticationHandler) Handle(r *RequestInfo) string {
	if !r.auth {
		return "You should login/signup to access this page"
	}
	return h.Next.Handle(r)
}

// CheckPremiumHandler - конкретный обработчик для проверки наличия подписки.
type CheckPremiumHandler struct {
	Next Handler
}

// Handle - метод обработки события (проверка наличия подписки).
func (h *CheckPremiumHandler) Handle(r *RequestInfo) string {
	if !r.premium {
		return "You should get premium subscription to access this page"
	}
	return h.Next.Handle(r)
}

// GetPremiumContentHandler - конкретный обработчик запроса на получение закрытого контента.
type GetPremiumContentHandler struct {
	Next Handler
}

// Handle - метод обработки события (получение закрытого контента).
func (h *GetPremiumContentHandler) Handle(r *RequestInfo) string {
	return "Here's your page with premium content"
}
