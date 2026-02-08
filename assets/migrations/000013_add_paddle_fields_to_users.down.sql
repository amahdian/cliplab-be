DROP TABLE IF EXISTS subscriptions;

ALTER TABLE users DROP COLUMN paddle_customer_id;
ALTER TABLE users DROP COLUMN subscription_id;
ALTER TABLE users DROP COLUMN price_id;
ALTER TABLE users DROP COLUMN subscription_status;
