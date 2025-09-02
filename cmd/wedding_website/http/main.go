package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"wedding_website/internal/app/telegram"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"wedding_website/internal/app/handlers"
	"wedding_website/internal/config"
	"wedding_website/internal/lib/logger"
	"wedding_website/internal/lib/logger/sl"
	"wedding_website/internal/storage/postgres"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	indexPath := filepath.Join(wd, "internal", "templates", "index.html")

	log.Printf("Serving index.html from: %s", indexPath)
	http.ServeFile(w, r, indexPath)
}

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	cfg.Telegram.BotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	cfg.Telegram.ChatID = os.Getenv("TELEGRAM_CHAT_ID")

	myLogger := logger.SetupLogger()

	myLogger.Info("Конфигурация логера подгружена")

	// Проверяем настройки Telegram
	if cfg.Telegram.BotToken == "" || cfg.Telegram.ChatID == "" {
		myLogger.Error("Токен Telegram-бота или идентификатор чата отсутствуют в конфигурации")
		return
	}

	storage, err := postgres.New(&cfg.DB)
	if err != nil {
		myLogger.Error("Ошибка подключения к БД", sl.Err(err))
		return
	}
	defer storage.Close()

	err = storage.Ping()
	if err != nil {
		myLogger.Error("БД не пингуется", sl.Err(err))
		return
	}

	myLogger.Info("БД подключена")

	telegramService := telegram.NewTelegramService(cfg.Telegram.BotToken, cfg.Telegram.ChatID, myLogger)

	// Инициализация handler с Telegram
	rsvpHandler := handlers.NewRSVPHandler(storage.DB, telegramService)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", HomeHandler)
	router.Post("/", rsvpHandler.HandleRSVP)

	fs := http.FileServer(http.Dir("internal/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	httpAddr := fmt.Sprintf(":%s", cfg.HTTP.Port)
	myLogger.Info("Запуск сервера на " + httpAddr)

	srv := &http.Server{
		Addr:         httpAddr,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	myLogger.Info("запуск сервера...")

	err = srv.ListenAndServe()
	if err != nil {
		myLogger.Error("ошибка запуска сервера", sl.Err(err))
		return
	}
}
