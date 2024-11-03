package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const telegramBotToken = "7442090726:AAHH8K2wbk6gfxCAxOl1Vwbj_CPlKxlZD7Q"

// Вопросы и ответы
var questions = map[string]string{
	"Какой квадрат числа 7?":                       "49",
	"Сколько будет 5 + 3?":                         "8",
	"Какой цвет получается при смешивании красного и синего?": "Фиолетовый",
	"Сколько планет в солнечной системе?":          "8",
	"Кто написал 'Войну и мир'?":                    "Лев Толстой",
}

var correctAnswers int

func main() {
	rand.Seed(time.Now().UnixNano())

	// Создание нового бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Создание обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u) // Получение обновлений

	for update := range updates {
		if update.Message == nil { // игнорировать не сообщения
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать в GriBotEduLev! Используйте /question для получения учебного вопроса или /help для помощи.")
			bot.Send(msg)

		case "/question":
			askQuestion(update.Message.Chat.ID, bot)

		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Команды:\n/question - получить учебный вопрос\n/help - получить список команд\n/statistics - увидеть количество правильных ответов")
			bot.Send(msg)

		case "/statistics":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Количество правильных ответов: %d", correctAnswers))
			bot.Send(msg)

		default:
			checkAnswer(update.Message.Chat.ID, update.Message.Text, bot)
		}
	}
}

// askQuestion выбирает случайный вопрос и отправляет его пользователю
func askQuestion(chatID int64, bot *tgbotapi.BotAPI) {
	// Выбор случайного вопроса
	questionKeys := make([]string, 0, len(questions))
	for q := range questions {
		questionKeys = append(questionKeys, q)
	}
	randomIndex := rand.Intn(len(questionKeys))
	question := questionKeys[randomIndex]

	msg := tgbotapi.NewMessage(chatID, question)
	bot.Send(msg)
}

// checkAnswer проверяет ответ пользователя и обновляет статистику
func checkAnswer(chatID int64, userAnswer string, bot *tgbotapi.BotAPI) {
	for _, correctAnswer := range questions {
		if userAnswer == correctAnswer {
			correctAnswers++
			msg := tgbotapi.NewMessage(chatID, "Правильно! Отличная работа.")
			bot.Send(msg)
			return
		}
	}
	msg := tgbotapi.NewMessage(chatID, "К сожалению, это неправильный ответ. Попробуйте еще раз!")
	bot.Send(msg)
}

