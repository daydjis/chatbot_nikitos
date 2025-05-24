# Chatbot Nikitos

Дерзкий Telegram-бот "Никитос", работающий с локальной LLaMA 3 моделью через Ollama.

---

## 🚀 Возможности
- Поднимается через Docker Compose
- Использует кастомную модель `nikitos`, созданную на базе `llama3`
- Реагирует на сообщения в Telegram и ведёт себя как Никитос
- Сохраняет историю общения для контекста

---

## 📦 Запуск

1. Склонируй репозиторий:
```bash
git clone https://github.com/daydjis/chatbot_nikitos.git
cd chatbot_nikitos
```

2. Создай `.env` файл с токеном:
```env
TELEGRAM_BOT_TOKEN=твой_бот_токен
```

3. Запусти проект:
```bash
docker compose up --build
```

4. Бот автоматически создаст модель и начнёт отвечать в Telegram.

---

## ⚙️ Структура
- `Modelfile` — описание поведения Никитоса
- `main.go` — код Telegram-бота на Go
- `nikitos_replicas.json` — редкие реплики Никитоса
- `docker-compose.yml` — всё, что нужно для запуска

---

## 📜 Пример поведения
**Пользователь:** как дела?

**Никитос:** норм

**Пользователь:** завтра в 10?

**Никитос:** рабские условия

---

## 🛠 TODO
- [ ] Команда `!забудь` для сброса контекста
- [ ] Поддержка нескольких персонализаций
- [ ] Интеграция RAG (памяти)

---

## 🧠 Базируется на
- [Ollama](https://ollama.com/)
- [LLaMA 3 (Meta)](https://ai.meta.com/llama/)
- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Go](https://golang.org)

---

Сделано по приколу ©️ daydjis