# MCP Notifications System

## Overview

WhatsApp MCP Server поддерживает встроенные MCP нотификации для real-time уведомлений о новых сообщениях и изменениях их статуса.

## Автоматические подписки

### Как это работает

Когда вы отправляете сообщение контакту через tool `send_message`, ваша MCP сессия **автоматически подписывается** на получение нотификаций от этого контакта. Это означает:

- ✅ Вы будете получать реальные уведомления о новых сообщениях от этого контакта
- ✅ Подписки уникальны для каждой MCP сессии
- ✅ Дублирующие подписки автоматически предотвращаются
- ✅ При отправке сообщения новому контакту автоматически создается подписка

### Пример использования

```json
// 1. Отправляем сообщение контакту
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "send_message",
    "arguments": {
      "to": "1234567890@s.whatsapp.net",
      "text": "Привет!"
    }
  }
}

// Ответ подтверждает отправку и подписку
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [{
      "type": "text",
      "text": "Message sent successfully to 1234567890@s.whatsapp.net. You are now subscribed to notifications from this chat."
    }]
  }
}

// 2. Когда контакт отправляет вам сообщение, вы получаете нотификацию:
{
  "jsonrpc": "2.0",
  "method": "notifications/message",
  "params": {
    "chat": "1234567890@s.whatsapp.net",
    "message_id": "3EB0ABCD1234",
    "from": "1234567890@s.whatsapp.net",
    "text": "Привет! Как дела?",
    "timestamp": 1234567890
  }
}
```

## Типы нотификаций

### Новое сообщение (`notifications/message`)

Отправляется когда приходит новое сообщение от контакта, на которого вы подписаны.

**Структура:**
```json
{
  "jsonrpc": "2.0",
  "method": "notifications/message",
  "params": {
    "chat": "string",          // JID чата
    "message_id": "string",    // ID сообщения
    "from": "string",          // JID отправителя
    "text": "string",          // Текст сообщения
    "timestamp": number        // Unix timestamp
  }
}
```

> **Примечание:** В настоящее время нотификации о статусе доставки/прочтения не отправляются. Эта информация сохраняется только в базе данных и доступна через API истории сообщений.

## Управление подписками

### Автоматическое управление

Система автоматически управляет подписками:

- **При отправке сообщения**: автоматически создается подписка на этот чат
- **При повторной отправке**: подписка не дублируется
- **Изоляция сессий**: каждая MCP сессия имеет свои подписки

### Архитектура

```
┌─────────────────────────────────────────┐
│       MCP Client (Your App)             │
│  ┌──────────┐         ┌──────────┐     │
│  │Session A │         │Session B │     │
│  └────┬─────┘         └────┬─────┘     │
└───────┼──────────────────┼────────────┘
        │                  │
        │  send_message    │  send_message
        ▼                  ▼
┌─────────────────────────────────────────┐
│         WhatsApp MCP Server              │
│                                          │
│  ┌────────────────────────────────────┐ │
│  │   Subscription Manager              │ │
│  │                                     │ │
│  │  Session A: [chat1, chat2]         │ │
│  │  Session B: [chat3]                │ │
│  └────────────────────────────────────┘ │
│                                          │
│         ▼ WhatsApp Events                │
│  ┌────────────────────────────────────┐ │
│  │   New Message from chat1           │ │
│  │   → Notify only Session A          │ │
│  └────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

## Преимущества

1. **Real-time уведомления** - мгновенные уведомления о новых сообщениях
2. **Автоматическое управление** - не нужно вручную подписываться
3. **Изоляция сессий** - каждый клиент получает только свои уведомления
4. **Нет дублирования** - система предотвращает повторные подписки
5. **Встроенный MCP протокол** - использует стандартные MCP нотификации

## Технические детали

### Session ID и контекст

Session ID автоматически передается в MCP протоколе:
- **HTTP транспорт**: через заголовок `Mcp-Session-Id`
- **stdio транспорт**: через внутренний контекст MCP

Session ID извлекается из контекста в методе `SendMessage` используя `server.ClientSessionFromContext(ctx)`, поэтому не нужно передавать его явно как параметр.

### SubscriptionManager

Центральный компонент для управления подписками:

```go
type SubscriptionManager struct {
    // sessionID -> chatJID -> subscribed
    subscriptions map[string]map[string]bool
    mutex         sync.RWMutex
    mcpServer     *server.MCPServer
}
```

**Основные методы:**
- `Subscribe(sessionID, chatJID)` - добавляет подписку
- `Unsubscribe(sessionID, chatJID)` - удаляет подписку
- `IsSubscribed(sessionID, chatJID)` - проверяет подписку
- `NotifyNewMessage(...)` - отправляет нотификацию о новом сообщении

### Интеграция с WhatsApp событиями

При получении события от WhatsApp:

1. **handleMessage** - обрабатывает входящие сообщения
   - Сохраняет в БД
   - Вызывает `NotifyNewMessage` для подписанных сессий

2. **handleReceipt** - обрабатывает статусы доставки/прочтения
   - Обновляет статусы в БД (без отправки нотификаций)

## Примеры интеграции

### Python Client

```python
import asyncio
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

async def main():
    server_params = StdioServerParameters(
        command="./whatsmeow-mcp",
        args=["stdio"]
    )
    
    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            await session.initialize()
            
            # Отправляем сообщение (автоподписка)
            result = await session.call_tool(
                "send_message",
                arguments={
                    "to": "1234567890@s.whatsapp.net",
                    "text": "Hello!"
                }
            )
            
            # Слушаем нотификации
            async for notification in session.notifications():
                if notification.method == "notifications/message":
                    params = notification.params
                    print(f"New message from {params['from']}: {params['text']}")

asyncio.run(main())
```

### TypeScript Client

```typescript
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";

const transport = new StdioClientTransport({
  command: "./whatsmeow-mcp",
  args: ["stdio"]
});

const client = new Client({
  name: "whatsapp-client",
  version: "1.0.0"
}, {
  capabilities: {}
});

await client.connect(transport);

// Отправляем сообщение
const result = await client.callTool({
  name: "send_message",
  arguments: {
    to: "1234567890@s.whatsapp.net",
    text: "Hello!"
  }
});

// Слушаем нотификации
client.setNotificationHandler({
  async handleNotification(notification) {
    if (notification.method === "notifications/message") {
      const { chat, from, text } = notification.params;
      console.log(`New message from ${from}: ${text}`);
    }
  }
});
```

## Ограничения и особенности

1. **Сессионность** - подписки привязаны к MCP сессии и удаляются при завершении сессии
2. **Только текстовые сообщения** - пока поддерживаются только текстовые сообщения
3. **Нет групповых чатов** - система работает с личными чатами и группами одинаково
4. **Автоматическая подписка** - отписаться нельзя, подписка создается автоматически

## Безопасность

- ✅ Каждая сессия изолирована
- ✅ Нотификации отправляются только подписанным сессиям
- ✅ Thread-safe операции с подписками (mutex)
- ✅ Автоматическая очистка при завершении сессии

