# NimbusNotify

## 1. Вимоги

### Функціональні

- Користувач може підписатись на оновлення погоди для конкретного міста з певною періодичністю
- Користувач повинен підтвердити підписку та мати змогу відписатись
- Користувач має отримувати регулярні оновлення погоди по кожній підписці
- Дані про погоду отримуємо з WeatherAPI.com

### Нефункціональні

- Система повинна бути доступною 99% часу
- Має витримувати 10000 тис. активних користувачів
- Листи мають відправлятись із затримкою менше ніж 1 год.
- Система має валідувати дані користувача
- Дані мають зберігатись без спотворень, без часткового збереження

### Обмеження

- Бюджет: мінімальна інфраструктура

## 2. Архітектура
![System Architecture](./arch-bg.png)

## 3. Оцінка навантаження

### Припущення

- 10000 активних підписок
- 80% користувачів отримують щоденні оновлення (8000), 20% — щогодинні (2000)
- Пік надсилання листів: 09:00–10:00 UTC
- 1000 унікальних міст
- Середній розмір листа: ~2 КБ

### Estimated Load

| Операція              | Частота (у пік)              | Коментарі                                  |
|-----------------------|------------------------------|--------------------------------------------|
| Запити на підписку    | 50/хв (~0.83/с)              | У періоди зростання кількості користувачів |
| Надсилання email      | 10000/год  (~2.8/с)          | Масова розсилка в межах однієї години      |
| Запити до Weather API | 1000/год  (~0.28/с)          | Для кожного унікального міста за годину    |
| Записи в БД           | (1000 + 10000)/год  (~3.1/с) | Оновлення погоди, статуси відправлення        |

### Коментарі
- використати Redis для кешу погоди
- потрібно реалізувати batch-відправлення листів
- необхідно коректно налаштувати пули з'єднань до БД та HTTP з'єднань

## 4. Структура БД
### 4.1 Оригінальна структура БД
```mermaid
erDiagram
    LOCATION ||--o{ SUBSCRIPTION : "has"
    LOCATION ||--o{ WEATHER : "has"
    SUBSCRIBER ||--o{ SUBSCRIPTION : "has"
    SUBSCRIPTION ||--o{ TOKEN : "has"

    LOCATION {
        id serial PK
        name varchar(60) UK
    }

    SUBSCRIBER {
        id serial PK
        email varchar(60) UK
        created_at timestamp
    }

    SUBSCRIPTION {
        id serial PK
        subscriber_id int FK
        location_id int FK
        location_name varchar(60)
        frequency frequency
        status subscription_status
        created_at timestamp
        updated_at timestamp
    }

    TOKEN {
        token char(36) PK
        subscription_id int FK
        type token_type
        created_at timestamp
        expires_at timestamp
        used_at timestamp
    }

    WEATHER {
        location_id int PK,FK
        last_updated timestamp PK
        fetched_at timestamp
        temperature numeric
        humidity numeric
        description varchar(400)
    }
```

### 4.2 Структура БД після видалення таблиці LOCATION
```mermaid
erDiagram
    SUBSCRIBER ||--o{ SUBSCRIPTION : "has"
    SUBSCRIPTION ||--o{ TOKEN : "has"

    SUBSCRIBER {
        id serial PK
        email varchar(60) UK
        created_at timestamp
    }

    SUBSCRIPTION {
        id serial PK
        subscriber_id int FK
        location_name varchar(60)
        frequency frequency
        status subscription_status
        created_at timestamp
        updated_at timestamp
    }

    TOKEN {
        token char(36) PK
        subscription_id int FK
        type token_type
        created_at timestamp
        expires_at timestamp
        used_at timestamp
    }

    WEATHER {
        location_name varchar(60) PK
        last_updated timestamp PK
        fetched_at timestamp
        temperature numeric
        humidity numeric
        description varchar(400)
    }
```


## 5. Sequence Diagrams

### 5.1 Підписка на оновлення погоди
```mermaid
sequenceDiagram
    participant U as User
    participant API as API Server
    participant DB as PostgreSQL
    participant MG as Mailgun

    U->>API: POST /subscribe
    API->>DB: Insert subscriber & subscription (pending)
    API->>DB: Create confirmation token
    API->>MG: Send "Confirmation" email
    API-->>U: Subscription is pending

    Note over U,API: User receives email and clicks confirm link

    U->>API: GET /confirm/:token
    API->>DB: Validate and mark confirmation token as used
    API->>DB: Update subscription status to "confirmed"
    API->>DB: Create unsubscribe token
    API->>MG: Send "Confirmation successful" email
    API-->>U: Subscription confirmed
```

### 5.2 Відписка від оновлення
```mermaid
sequenceDiagram
    participant U as User
    participant API as API Server
    participant DB as PostgreSQL
    participant MG as Mailgun

    U->>API: GET /unsubscribe/:token
    API->>DB: Validate unsubscribe token
    API->>DB: Delete subscription
    API->>MG: Send "End of subscription" email
    API-->>U: Unsubscribed successfully
```