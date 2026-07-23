-- 数据库初始化脚本

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id INTEGER UNIQUE NOT NULL,
    username TEXT,
    first_name TEXT,
    last_name TEXT,
    is_admin BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 订阅表
CREATE TABLE IF NOT EXISTS subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    plan_id INTEGER NOT NULL,
    status TEXT DEFAULT 'active',
    expires_at DATETIME NOT NULL,
    max_bots INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 套餐表
CREATE TABLE IF NOT EXISTS plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    duration_days INTEGER NOT NULL,
    max_bots INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 兑换码表
CREATE TABLE IF NOT EXISTS redeem_codes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT UNIQUE NOT NULL,
    plan_id INTEGER NOT NULL,
    status TEXT DEFAULT 'unused',
    used_by INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    used_at DATETIME,
    FOREIGN KEY (plan_id) REFERENCES plans(id),
    FOREIGN KEY (used_by) REFERENCES users(id)
);

-- 子Bot表
CREATE TABLE IF NOT EXISTS worker_bots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    bot_token TEXT NOT NULL,
    bot_username TEXT NOT NULL,
    bot_nickname TEXT,
    owner_telegram_id INTEGER NOT NULL,
    status TEXT DEFAULT 'stopped',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 子Bot配置表
CREATE TABLE IF NOT EXISTS bot_configs (
    bot_id INTEGER PRIMARY KEY,
    owner_telegram_id INTEGER NOT NULL,
    ai_enabled BOOLEAN DEFAULT 1,
    use_custom_ai BOOLEAN DEFAULT 0,
    custom_ai_endpoint TEXT,
    custom_ai_key TEXT,
    custom_ai_model TEXT,
    whitelist_threshold INTEGER DEFAULT 3,
    notify_on_block BOOLEAN DEFAULT 1,
    notify_on_ai_block BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bot_id) REFERENCES worker_bots(id)
);

-- 客户表
CREATE TABLE IF NOT EXISTS customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    bot_id INTEGER NOT NULL,
    telegram_id INTEGER NOT NULL,
    username TEXT,
    first_name TEXT,
    last_name TEXT,
    verified_count INTEGER DEFAULT 0,
    total_messages INTEGER DEFAULT 0,
    is_whitelisted BOOLEAN DEFAULT 0,
    is_blacklisted BOOLEAN DEFAULT 0,
    first_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_message_at DATETIME,
    whitelisted_at DATETIME,
    blacklisted_at DATETIME,
    blacklist_reason TEXT,
    is_manual_trust BOOLEAN DEFAULT 0,
    UNIQUE(bot_id, telegram_id),
    FOREIGN KEY (bot_id) REFERENCES worker_bots(id)
);

-- 拦截日志表
CREATE TABLE IF NOT EXISTS block_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    bot_id INTEGER NOT NULL,
    telegram_id INTEGER NOT NULL,
    username TEXT,
    first_name TEXT,
    last_name TEXT,
    block_reason TEXT,
    ai_confidence REAL,
    message_content TEXT,
    message_type TEXT,
    blocked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_false_positive BOOLEAN DEFAULT 0,
    FOREIGN KEY (bot_id) REFERENCES worker_bots(id)
);

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    key TEXT PRIMARY KEY,
    value TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    plan_id INTEGER NOT NULL,
    order_no TEXT UNIQUE NOT NULL,
    amount REAL NOT NULL,
    status TEXT DEFAULT 'pending',
    payment_url TEXT,
    paid_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (plan_id) REFERENCES plans(id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_expires_at ON subscriptions(expires_at);
CREATE INDEX IF NOT EXISTS idx_worker_bots_user_id ON worker_bots(user_id);
CREATE INDEX IF NOT EXISTS idx_customers_bot_telegram ON customers(bot_id, telegram_id);
CREATE INDEX IF NOT EXISTS idx_customers_whitelist ON customers(bot_id, is_whitelisted);
CREATE INDEX IF NOT EXISTS idx_customers_blacklist ON customers(bot_id, is_blacklisted);
CREATE INDEX IF NOT EXISTS idx_block_logs_bot_time ON block_logs(bot_id, blocked_at);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
