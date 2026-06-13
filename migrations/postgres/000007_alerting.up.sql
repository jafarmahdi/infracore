-- 000007_alerting.up.sql
-- Alerting & Notification Engine schema

-- ============================================================
-- NOTIFICATION CHANNELS
-- ============================================================
CREATE TABLE notification_channels (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    name            VARCHAR(255) NOT NULL,
    channel_type    VARCHAR(50)  NOT NULL
                        CHECK (channel_type IN ('email','sms','whatsapp','webhook',
                                                'slack','teams','pagerduty','telegram')),
    -- Stored as encrypted JSONB; decrypted at runtime
    -- email:     {"to": ["ops@co.com"], "from": "alerts@co.com", "smtp_id": "uuid"}
    -- sms:       {"to": ["+1234567890"], "provider": "twilio", "account_sid": "..."}
    -- webhook:   {"url": "https://...", "method": "POST", "headers": {...}, "secret": "..."}
    -- slack:     {"webhook_url": "https://hooks.slack.com/...", "channel": "#alerts"}
    configuration   JSONB        NOT NULL DEFAULT '{}',
    is_enabled      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_notification_channels_updated_at
    BEFORE UPDATE ON notification_channels
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- ESCALATION POLICIES
-- ============================================================
CREATE TABLE escalation_policies (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, name)
);

CREATE TRIGGER set_escalation_policies_updated_at
    BEFORE UPDATE ON escalation_policies
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- ESCALATION STEPS
-- ============================================================
CREATE TABLE escalation_steps (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id               UUID        NOT NULL REFERENCES escalation_policies(id) ON DELETE CASCADE,
    step_order              INTEGER      NOT NULL CHECK (step_order >= 1),
    escalate_after_minutes  INTEGER      NOT NULL CHECK (escalate_after_minutes > 0),
    channel_ids             JSONB        NOT NULL DEFAULT '[]',
    user_ids                JSONB        NOT NULL DEFAULT '[]',
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (policy_id, step_order)
);

-- ============================================================
-- ALERT RULES
-- ============================================================
CREATE TABLE alert_rules (
    id                          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                   UUID        NOT NULL REFERENCES tenants(id),
    escalation_policy_id        UUID        REFERENCES escalation_policies(id),
    name                        VARCHAR(255) NOT NULL,
    description                 TEXT,
    rule_type                   VARCHAR(50)  NOT NULL
                                    CHECK (rule_type IN ('threshold','anomaly','absence','change','composite')),
    -- Scope: which hosts/groups this rule applies to
    scope_type                  VARCHAR(50)
                                    CHECK (scope_type IN ('host','agent_group','site','tag','all')),
    scope_ids                   JSONB        NOT NULL DEFAULT '[]',

    -- Metric targeting
    metric_name                 VARCHAR(100),
    check_type                  VARCHAR(100),
    label_filters               JSONB        NOT NULL DEFAULT '{}',

    -- Thresholds
    condition                   VARCHAR(10)
                                    CHECK (condition IN ('gt','lt','gte','lte','eq','neq')),
    threshold_warning           DOUBLE PRECISION,
    threshold_critical          DOUBLE PRECISION,
    threshold_emergency         DOUBLE PRECISION,

    -- Evaluation
    evaluation_interval_seconds INTEGER      NOT NULL DEFAULT 60,
    evaluation_periods          INTEGER      NOT NULL DEFAULT 1,  -- consecutive periods to confirm

    -- Absence check: fire if no data for N seconds
    absence_window_seconds      INTEGER,

    severity                    VARCHAR(20)  NOT NULL DEFAULT 'warning'
                                    CHECK (severity IN ('info','warning','critical','emergency')),
    is_enabled                  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_alert_rules_updated_at
    BEFORE UPDATE ON alert_rules
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- ALERT RULE ↔ NOTIFICATION CHANNELS  (M:N)
-- ============================================================
CREATE TABLE alert_rule_channels (
    rule_id     UUID NOT NULL REFERENCES alert_rules(id) ON DELETE CASCADE,
    channel_id  UUID NOT NULL REFERENCES notification_channels(id) ON DELETE CASCADE,
    PRIMARY KEY (rule_id, channel_id)
);

-- ============================================================
-- ALERT EVENTS  (fired alert instances)
-- ============================================================
CREATE TABLE alert_events (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID        NOT NULL REFERENCES tenants(id),
    rule_id                 UUID        NOT NULL REFERENCES alert_rules(id),
    host_id                 UUID        REFERENCES monitored_hosts(id),
    check_id                UUID        REFERENCES check_definitions(id),
    severity                VARCHAR(20)  NOT NULL
                                CHECK (severity IN ('info','warning','critical','emergency')),
    status                  VARCHAR(20)  NOT NULL DEFAULT 'firing'
                                CHECK (status IN ('firing','resolved','acknowledged','silenced')),
    title                   VARCHAR(500) NOT NULL,
    message                 TEXT,
    metric_name             VARCHAR(100),
    metric_value            DOUBLE PRECISION,
    threshold_value         DOUBLE PRECISION,
    labels                  JSONB        NOT NULL DEFAULT '{}',
    fired_at                TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    resolved_at             TIMESTAMPTZ,
    acknowledged_at         TIMESTAMPTZ,
    acknowledged_by         UUID        REFERENCES users(id),
    acknowledgement_note    TEXT,
    silenced_until          TIMESTAMPTZ,
    silenced_by             UUID        REFERENCES users(id),
    escalation_policy_id    UUID        REFERENCES escalation_policies(id),
    current_escalation_step INTEGER      DEFAULT 0,
    last_escalated_at       TIMESTAMPTZ
);

-- ============================================================
-- ALERT NOTIFICATIONS  (delivery log per channel)
-- ============================================================
CREATE TABLE alert_notifications (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL,
    event_id        UUID        NOT NULL REFERENCES alert_events(id) ON DELETE CASCADE,
    channel_id      UUID        NOT NULL REFERENCES notification_channels(id),
    status          VARCHAR(20)  NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending','sent','failed','skipped')),
    sent_at         TIMESTAMPTZ,
    error_message   TEXT,
    retry_count     INTEGER      NOT NULL DEFAULT 0,
    payload         JSONB        NOT NULL DEFAULT '{}',  -- the rendered notification payload
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ============================================================
-- SILENCES  (suppress alerts matching criteria for a period)
-- ============================================================
CREATE TABLE silences (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id),
    created_by  UUID        NOT NULL REFERENCES users(id),
    matchers    JSONB        NOT NULL DEFAULT '{}',  -- {host_id, rule_id, labels}
    starts_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    ends_at     TIMESTAMPTZ  NOT NULL,
    comment     TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
