ALTER TABLE IF EXISTS ordered_items_scan_history
    ADD COLUMN qr_id TEXT NOT NULL,
    DROP COLUMN order_id;



DROP TABLE IF EXISTS items_qr_data;

CREATE TABLE IF NOT EXISTS ordered_items
(
    id          BIGSERIAL PRIMARY KEY,
    order_id    INT REFERENCES orders (id) NOT NULL,
    qr_id       TEXT                       NOT NULL,
    data        json                       NOT NULL,
    scanned_at  TIMESTAMP WITH TIME ZONE   NOT NULL,
    archived_at TIMESTAMP
)
