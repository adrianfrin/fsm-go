# 消息

核心包只定义发布接口，不绑定具体消息系统。

## 发布接口

`MessagePublisher` 负责发布消息。

## Kafka

`messaging/kafka` 提供 Kafka publisher。

## Outbox

Outbox 用于把本地流程推进和对外消息发送衔接起来。流程运行时先在本地事务里写 Outbox，后台发布器再把消息发到 Kafka。
