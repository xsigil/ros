PRAGMA foreign_keys = ON;

-- 1. 必然度と半減期設定 (NULLは半減期無限大 = 永久保存)
CREATE TABLE IF NOT EXISTS inevitability_config (
    inevitability TEXT PRIMARY KEY,
    half_life_days REAL,
    description TEXT NOT NULL
);

INSERT OR IGNORE INTO inevitability_config (inevitability, half_life_days, description) VALUES
    ('transient_noise', 7.0,  'Surface-level noise / gestures (rapid decay)'),
    ('daily',           30.0, 'Baseline habitual behaviors (monthly decay)'),
    ('critical_moment', NULL, 'Core architectural integrity / breach (permanent exemption)');

-- 2. 重みマスタ (type x sign x inevitability)
CREATE TABLE IF NOT EXISTS odds_weights (
    type TEXT NOT NULL CHECK(type IN ('adult_customer', 'adult_non_customer', 'child_customer', 'child_non_customer', 'employee')),
    sign TEXT NOT NULL CHECK(sign IN ('pos', 'neg')),
    inevitability TEXT NOT NULL,
    log_odds_delta REAL NOT NULL,
    PRIMARY KEY (type, sign, inevitability),
    FOREIGN KEY (inevitability) REFERENCES inevitability_config(inevitability)
);

-- 初期シードデータの投入
INSERT OR IGNORE INTO odds_weights VALUES
    -- adult_customer
    ('adult_customer', 'pos', 'daily',            2.0),
    ('adult_customer', 'neg', 'daily',           -3.0),
    ('adult_customer', 'pos', 'critical_moment', 15.0),
    ('adult_customer', 'neg', 'critical_moment', -30.0),
    ('adult_customer', 'pos', 'transient_noise',   0.3),
    ('adult_customer', 'neg', 'transient_noise',  -0.5),
    -- adult_non_customer
    ('adult_non_customer', 'pos', 'daily',            1.5),
    ('adult_non_customer', 'neg', 'daily',           -4.0),
    ('adult_non_customer', 'pos', 'critical_moment', 12.0),
    ('adult_non_customer', 'neg', 'critical_moment', -30.0),
    ('adult_non_customer', 'pos', 'transient_noise',   0.2),
    ('adult_non_customer', 'neg', 'transient_noise',  -1.0),
    -- employee
    ('employee', 'pos', 'daily',            2.5),
    ('employee', 'neg', 'daily',           -3.5),
    ('employee', 'pos', 'critical_moment', 20.0),
    ('employee', 'neg', 'critical_moment', -35.0),
    ('employee', 'pos', 'transient_noise',   0.5),
    ('employee', 'neg', 'transient_noise',  -0.5),
    -- child_customer
    ('child_customer', 'pos', 'daily',            1.0),
    ('child_customer', 'neg', 'daily',           -1.0),
    ('child_customer', 'pos', 'critical_moment', 18.0),
    ('child_customer', 'neg', 'critical_moment', -20.0),
    ('child_customer', 'pos', 'transient_noise',   0.2),
    ('child_customer', 'neg', 'transient_noise',  -0.2),
    -- child_non_customer
    ('child_non_customer', 'pos', 'daily',            0.8),
    ('child_non_customer', 'neg', 'daily',           -0.5),
    ('child_non_customer', 'pos', 'critical_moment', 20.0),
    ('child_non_customer', 'neg', 'critical_moment', -15.0),
    ('child_non_customer', 'pos', 'transient_noise',   0.1),
    ('child_non_customer', 'neg', 'transient_noise',  -0.1);

-- 3. 人物テーブル (状態ステータスを含む)
CREATE TABLE IF NOT EXISTS person (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    current_type TEXT NOT NULL CHECK(current_type IN ('adult_customer', 'adult_non_customer', 'child_customer', 'child_non_customer', 'employee')),
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'monitoring_only', 'quarantined')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 4. 観測イベント履歴 (Append-Only)
CREATE TABLE IF NOT EXISTS history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    person_id TEXT NOT NULL,
    episode TEXT NOT NULL,
    type TEXT NOT NULL,
    sign TEXT NOT NULL,
    inevitability TEXT NOT NULL,
    log_odds_applied REAL NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (person_id) REFERENCES person(id),
    FOREIGN KEY (type, sign, inevitability) REFERENCES odds_weights(type, sign, inevitability)
);

CREATE INDEX IF NOT EXISTS idx_history_person ON history(person_id);
CREATE INDEX IF NOT EXISTS idx_history_timestamp ON history(timestamp);

-- 5. 時間減衰集計ビュー
CREATE VIEW IF NOT EXISTS v_person_effective_odds AS
WITH event_decay AS (
    SELECT 
        h.person_id,
        h.log_odds_applied,
        h.inevitability,
        cfg.half_life_days,
        (julianday('now') - julianday(h.timestamp)) AS days_elapsed,
        CASE 
            WHEN cfg.half_life_days IS NULL THEN 1.0
            ELSE POWER(0.5, (julianday('now') - julianday(h.timestamp)) / cfg.half_life_days)
        END AS retention_factor
    FROM history h
    JOIN inevitability_config cfg ON h.inevitability = cfg.inevitability
)
SELECT 
    p.id,
    p.name,
    p.current_type,
    p.status,
    ROUND(COALESCE(SUM(ed.log_odds_applied * ed.retention_factor), 0.0), 2) AS current_log_odds_db,
    COUNT(ed.log_odds_applied) AS total_events,
    ROUND(COALESCE(SUM(CASE WHEN ed.half_life_days IS NULL THEN ed.log_odds_applied ELSE 0 END), 0.0), 2) AS permanent_anchor_db,
    ROUND(COALESCE(SUM(CASE WHEN ed.half_life_days IS NOT NULL THEN ed.log_odds_applied * ed.retention_factor ELSE 0 END), 0.0), 2) AS transient_active_db
FROM person p
LEFT JOIN event_decay ed ON p.id = ed.person_id
GROUP BY p.id;
