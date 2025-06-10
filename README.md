# Telegram Bot on Go with MongoDB and Microservices

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![MongoDB](https://img.shields.io/badge/MongoDB-7.0+-47A248?logo=mongodb)](https://www.mongodb.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Производственный Telegram-бот на Golang с поддержкой:
- **Конкурентной обработки** через worker pool (горутины)
- **MongoDB** для хранения данных
- **Микросервисной архитектуры** (опционально через RabbitMQ/gRPC)
- Деплой в Docker/Kubernetes

## 🚀 Особенности

✔ **Гибкая архитектура**  
- Модульная структура (бот, хранилище, воркеры)  
- Поддержка graceful shutdown  

✔ **Работа с MongoDB**  
- Хранение пользователей, состояний (FSM), логов  
- Индексы для быстрых запросов  

✔ **Конкурентность**  
- Пул воркеров для обработки сообщений  
- Каналы для межгорутинного взаимодействия  

✔ **Масштабируемость**  
- Готовность к переходу на микросервисы (через RabbitMQ или gRPC)  
- Конфиги через environment variables  

## 📦 Установка

### Требования
- Go 1.21+
- MongoDB 7.0+
- Telegram Bot API Token

### Запуск
```bash
# Клонировать репозиторий
git clone https://github.com/yourusername/telegram-bot-go.git
cd telegram-bot-go

# Установить зависимости
go mod download

# Запустить (предварительно настроив .env)
cp .env.example .env
go run cmd/bot/main.go