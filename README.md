# # 🚲 BikeStoreGolang — Development Guide

Добро пожаловать в руководство по разработке проекта \*\*BikeStoreGolang\*\*. Здесь собраны все инструкции для локальной настройки, запуска и тестирования микросервисов.

\---

\## 🛠️ Настройка среды

Перед началом работы убедитесь, что установлены следующие зависимости:

\- \*\*Go\*\* \`1.22+\`
\- \*\*Docker\*\* + \`docker-compose\`
\- \*\*Node.js\*\* \`18+\` (для фронтенда)
\- \*\*protoc\*\* (Protocol Buffers — генерация gRPC-кода)

\### 🔽 Клонирование репозитория

\`\`\`bash
git clone -b development https://github.com/your-repo/BikeStoreGolang.git
cd BikeStoreGolang
\`\`\`

\### 🚀 Запуск инфраструктуры

\`\`\`bash
docker-compose -f deployments/docker-compose.dev.yml up -d postgres redis nats
\`\`\`

\---

\## 👨‍💻 Работа с кодом

\### 🧬 Генерация gRPC-кода

\`\`\`bash
make generate-proto
\`\`\`

\### ▶️ Запуск сервисов локально

Например, для \`auth-service\`:

\`\`\`bash
cd services/auth-service
go run cmd/main.go
\`\`\`

\---

\## 🔃 Правила внесения изменений

Создайте новую ветку от \`development\`:

\`\`\`bash
git checkout -b feature/auth-service-login
\`\`\`

После завершения работы:

\`\`\`bash
git add .
git commit -m "feat(auth): add login endpoint"
git push origin feature/auth-service-login
\`\`\`

Откройте \*\*Pull Request\*\* в ветку \`development\`.

\---

\## 🧪 Тестирование

\### ✅ Юнит-тесты

\`\`\`bash
cd services/auth-service
go test ./... -v
\`\`\`

\### 🔗 Интеграционные тесты (требуется Docker)

\`\`\`bash
go test ./... -tags=integration
\`\`\`

\---

\## 🔌 Полезные команды

| Команда | Описание |
|----|----|
| \`make migrate\` | Запуск миграций |
| \`make generate-proto\` | Генерация gRPC-кода |
| \`docker-compose logs -f\` | Просмотр логов |

\---

\## ⚠️ Troubleshooting

\- \*\*Ошибки NATS\*\*: Убедитесь, что NATS запущен (\`docker-compose up nats\`).
\- \*\*Проблемы с gRPC\*\*: Перегенерируйте код (\`make generate-proto\`).

\---

\## 📌 Ключевые различия

\### \`main\`:

\- Описание продукта (что это, стек, структура)
\- Инструкции для деплоя
\- Минимальные технические детали

\### \`development\`:

\- Пошаговая настройка среды
\- Git-воркфлоу и стандарты коммитов
\- Локальный запуск и тестирование
\- Частые проблемы и решения

> Такой подход разделяет аудиторию:
>
> * \*\*main\*\* — для пользователей и менеджеров
> * \*\*development\*\* — только для разработчиков


