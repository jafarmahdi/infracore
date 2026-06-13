-- 000007_alerting.down.sql
DROP TABLE IF EXISTS silences;
DROP TABLE IF EXISTS alert_notifications;
DROP TABLE IF EXISTS alert_events;
DROP TABLE IF EXISTS alert_rule_channels;
DROP TABLE IF EXISTS alert_rules;
DROP TABLE IF EXISTS escalation_steps;
DROP TABLE IF EXISTS escalation_policies;
DROP TABLE IF EXISTS notification_channels;
