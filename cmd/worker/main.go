package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/acnoway/tg_forward_master/internal/config"
	"github.com/acnoway/tg_forward_master/internal/database"
	"github.com/acnoway/tg_forward_master/internal/worker/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.env")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 运行数据库迁移
	if err := db.RunMigrations("migrations/001_init.sql"); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("✅ 数据库初始化完成")

	// 创建Bot
	bot, err := tgbotapi.NewBotAPI(cfg.MasterBotToken) // FIXME: 需要WorkerBotToken配置
	if err != nil {
		log.Fatalf("创建Bot失败: %v", err)
	}

	bot.Debug = false
	log.Printf("✅ 子Bot已启动: @%s", bot.Self.UserName)

	// 获取Bot ID
	workerBotID, err := getBotID(db, bot.Self.UserName)
	if err != nil {
		log.Fatalf("获取Bot ID失败: %v", err)
	}

	// 获取主人ID
	ownerID, err := getOwnerID(db, workerBotID)
	if err != nil {
		log.Fatalf("获取主人ID失败: %v", err)
	}

	// 创建处理器
	handler := handlers.NewWorkerHandler(bot, db, workerBotID, ownerID)

	// 获取更新
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// 优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("收到退出信号，正在关闭...")
		bot.StopReceivingUpdates()
		os.Exit(0)
	}()

	// 处理更新
	log.Println("🚀 开始接收消息...")
	for update := range updates {
		if update.Message != nil {
			handler.HandleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			handler.HandleCallback(update.CallbackQuery)
		}
	}
}

// getBotID 获取子Bot的ID
func getBotID(db *database.DB, botUsername string) (int64, error) {
	query := `SELECT id FROM worker_bots WHERE bot_username = ?`
	var id int64
	err := db.QueryRow(query, botUsername).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// getOwnerID 获取Bot主人的Telegram ID
func getOwnerID(db *database.DB, botID int64) (int64, error) {
	query := `SELECT owner_telegram_id FROM worker_bots WHERE id = ?`
	var ownerID int64
	err := db.QueryRow(query, botID).Scan(&ownerID)
	if err != nil {
		return 0, err
	}
	return ownerID, nil
}