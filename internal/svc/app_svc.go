package svc

type AppService struct {
	// Добавьте зависимости (репозиторий, кэш и т.д.)
}

func NewAppService() *AppService {
	return &AppService{}
}

func (s *AppService) CheckUserAccess(userID int64) (bool, error) {
	// Проверяем подписку/баланс/доступ
	// Например, запрос в БД или Redis
	return true, nil // или false, errors.New("no access")
}
