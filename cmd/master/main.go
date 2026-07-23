package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yourusername/tg_forward_master/internal/config"
	"github.com/yourusername/tg_forward_master/internal/database"
	"github.com/yourusername/tg_forward_master/internal/master/handlers"
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
	bot, err := tgbotapi.NewBotAPI(cfg.MasterBotToken)
	if err != nil {
		log.Fatalf("创建Bot失败: %v", err)
	}

	bot.Debug = false
	log.Printf("✅ 主控Bot已启动: @%s", bot.Self.UserName)

	// 创建处理器
	handler := handlers.NewMasterHandler(bot, db, cfg)

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
