-- Manual usage

-- 1. Позволяет создавать новые объекты (таблицы, индексы, последовательности)
GRANT CREATE ON SCHEMA audit_db TO svc_audit;

-- 2. Позволяет видеть схему
GRANT USAGE ON SCHEMA audit_db TO svc_audit;

-- 3. Если таблицы уже созданы другим юзером (например, postgres),
-- нужно сделать test их владельцем или дать полные права:
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA audit_db TO svc_audit;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA audit_db TO svc_audit;

-- 4. Чтобы будущие объекты тоже были под контролем:
ALTER DEFAULT PRIVILEGES IN SCHEMA audit_db
    GRANT ALL PRIVILEGES ON TABLES TO svc_audit;
ALTER DEFAULT PRIVILEGES IN SCHEMA audit_db
    GRANT ALL PRIVILEGES ON SEQUENCES TO svc_audit;
