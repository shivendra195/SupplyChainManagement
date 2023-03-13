ALTER TABLE IF EXISTS ordered_items_scan_history
    ADD COLUMN item_id int references ordered_items(id),
    ADD COLUMN scanned_address text