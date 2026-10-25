-- 000005_rename_id_to_activity_id.up.sql
ALTER TABLE activities
RENAME COLUMN id TO activity_id;