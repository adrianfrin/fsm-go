package mysql

const Schema = `
CREATE TABLE IF NOT EXISTS fsm_entity (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    machine         VARCHAR(64)     NOT NULL COMMENT '状态机名称',
    machine_version VARCHAR(32)     NOT NULL COMMENT '状态机版本',
    entity_id       VARCHAR(128)    NOT NULL COMMENT '业务对象ID',
    state           VARCHAR(64)     NOT NULL COMMENT '当前状态',
    revision        BIGINT          NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
    data            JSON            NULL COMMENT '状态机运行所需扩展数据',
    create_at       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    update_at       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    deleted         TINYINT         NOT NULL DEFAULT 0 COMMENT '是否删除：0否，1是',
    PRIMARY KEY (id),
    UNIQUE KEY uk_fsm_entity_machine_entity (machine, entity_id),
    KEY idx_fsm_entity_machine_state (machine, state)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='状态机实体表';

CREATE TABLE IF NOT EXISTS fsm_state_log (
    id                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    machine           VARCHAR(64)     NOT NULL COMMENT '状态机名称',
    machine_version   VARCHAR(32)     NOT NULL COMMENT '状态机版本',
    entity_id         VARCHAR(128)    NOT NULL COMMENT '业务对象ID',
    event             VARCHAR(64)     NOT NULL COMMENT '触发事件',
    from_state        VARCHAR(64)     NOT NULL COMMENT '原状态',
    to_state          VARCHAR(64)     NOT NULL COMMENT '目标状态',
    transition_name   VARCHAR(128)    NOT NULL COMMENT '迁移名称',
    actor_id          VARCHAR(128)    NULL COMMENT '操作者ID',
    request_id        VARCHAR(128)    NULL COMMENT '请求ID',
    idempotency_key   VARCHAR(256)    NULL COMMENT '幂等键',
    payload           JSON            NULL COMMENT '请求载荷',
    create_at         DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    update_at         DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    deleted           TINYINT         NOT NULL DEFAULT 0 COMMENT '是否删除：0否，1是',
    PRIMARY KEY (id),
    KEY idx_fsm_state_log_entity (machine, entity_id, id),
    KEY idx_fsm_state_log_request (request_id),
    KEY idx_fsm_state_log_idempotency (machine, idempotency_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='状态机迁移日志表';

CREATE TABLE IF NOT EXISTS fsm_idempotency (
    id                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    machine           VARCHAR(64)     NOT NULL COMMENT '状态机名称',
    idempotency_key   VARCHAR(256)    NOT NULL COMMENT '幂等键',
    entity_id         VARCHAR(128)    NOT NULL COMMENT '业务对象ID',
    event             VARCHAR(64)     NOT NULL COMMENT '触发事件',
    result            JSON            NULL COMMENT '历史执行结果',
    status            VARCHAR(32)     NOT NULL COMMENT '处理状态',
    create_at         DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    update_at         DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    deleted           TINYINT         NOT NULL DEFAULT 0 COMMENT '是否删除：0否，1是',
    PRIMARY KEY (id),
    UNIQUE KEY uk_fsm_idempotency_machine_key (machine, idempotency_key),
    KEY idx_fsm_idempotency_entity (machine, entity_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='状态机幂等表';

CREATE TABLE IF NOT EXISTS fsm_outbox (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    topic           VARCHAR(128)    NOT NULL COMMENT '消息主题',
    msg_key         VARCHAR(256)    NOT NULL COMMENT '消息键',
    payload         JSON            NOT NULL COMMENT '消息内容',
    status          VARCHAR(32)     NOT NULL DEFAULT 'PENDING' COMMENT '消息状态',
    retry_count     INT             NOT NULL DEFAULT 0 COMMENT '重试次数',
    next_retry_at   DATETIME(3)     NULL COMMENT '下次重试时间',
    published_at    DATETIME(3)     NULL COMMENT '发布时间',
    create_at       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    update_at       DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    deleted         TINYINT         NOT NULL DEFAULT 0 COMMENT '是否删除：0否，1是',
    PRIMARY KEY (id),
    KEY idx_fsm_outbox_publish (status, next_retry_at, id),
    KEY idx_fsm_outbox_topic_key (topic, msg_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='状态机Outbox消息表';
`
