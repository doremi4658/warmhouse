# Project_template, просьба открыть readme и открыть в code тогда код схем копируется коректно и работает в plantuml


# Задание 1. Анализ и планирование

### 1. Описание функциональности монолитного приложения

**Управление отоплением:**

- Пользователи могут: Включать/выключать систему отопления, устанавливать температуру.
- Система поддерживает: Одностороннюю синхронную коммуникацию "сервер -> устройство". Сервер отправляет команды управления на модуль отопления по запросу пользователя.

**Мониторинг температуры:**

- Пользователи могут: Просматривать текущую температуру в доме.
Система поддерживает: Одностороннюю синхронную коммуникацию "сервер -> устройство". Сервер периодически опрашивает датчик температуры для получения актуальных данных.

### 2. Анализ архитектуры монолитного приложения

Технологический стек: Бэкенд-приложение написано на Go, в качестве системы управления базами данных используется PostgreSQL. Архитектурный стиль: Монолит. Весь функционал (бизнес-логика, аутентификация, работа с БД, API) упакован в единую, неразделимую кодовую базу и процесс.

### 3. Определение доменов и границы контекстов

Отопление, включить/выключить натсроить температуру. 

### **4. Проблемы монолитного решения**

На данный момент функциона, включение/выключение и уставновка температуры, пробелм в архитектуре нет, но если приложение хотят расширить исходя из целевого решения видны следующие сложности:
Невозможность независимого масштабирования
Низкая отказоустойчивость

### 5. Визуализация контекста системы — диаграмма С4

@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Context.puml

title Диаграмма контекста системы "Теплый дом"

Person(пользователь, "Пользователь", "Хозяин дома, управляет отоплением, мониторит температуру")
Person(специалист, "Технический специалист", "Подключает систему отопления")

System(монолит, "Монолитное приложение", "Приложение на Go, управляет отоплением и температурой")
  System_Ext(датчик_температуры, "Датчик температуры", "Устройство меряет температуру в доме")
  System_Ext(модуль_отопления, "Модуль отопления", "Устройство включает/выелючает отопление")
  System_Ext(бд, "PostgreSQL", "База данных приложения")

Rel(пользователь, монолит, "Удаленно включает/выключает отопление в доме, просматривает температуру")
Rel(специалист, монолит, "Подключает систему")

Rel(монолит, датчик_температуры, "Запрашивает текущую температуру")
Rel(монолит, модуль_отопления, "Включает/выключает отопление",)
Rel(монолит, бд, "Сохраняет и читает данные")

@enduml

```markdown
[Код необходимо вставить в PlantUML, как поделиться ссылкой на прямую с моей схемой не нашел, буду благодарян за подсказку как это сделать](https://www.planttext.com/)
```

# Задание 2. Проектирование микросервисной архитектуры


**Диаграмма контейнеров (Containers)**

@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml
title Диаграмма контейнеров системы "Умный дом"

Person(user, "Пользователь", "Управляет умным домом через веб-приложение")

System_Boundary(ecosystem, "Экосистема Умный Дом") {
    Container(web_app, "Веб-приложение", "Интерфейс управления умным домом\nи агрегатор API микросервисов")
    
    Container(heating_s, "Heating Service", "Go", "Управление отоплением\n- Температура\n- Режимы обогрева\n- График работы")
    Container(lighting_s, "Lighting Service", "Go", "Управление освещением\n- Включение/выключение\n- Яркость\n- Таймеры") 
    Container(gate_s, "Gate Service", "Go", "Управление воротами\n- Открытие/закрытие\n- Статус\n- История доступа")
    Container(monitoring_s, "Monitoring Service", "Go", "Наблюдение за домом\n- Камеры\n- Датчики движения\n- История событий")
    Container(device_s, "Device Service", "Go", "Общее управление устройствами\n- Регистрация\n- Статус\n- Обновления")
    
    ContainerDb(heating_db, "Heating DB", "PostgreSQL", "Данные отопления\n- Температуры\n- Графики\n- Настройки")
    ContainerDb(lighting_db, "Lighting DB", "PostgreSQL", "Данные освещения\n- Состояния\n- Таймеры\n- Сцены")
    ContainerDb(gate_db, "Gate DB", "PostgreSQL", "Данные ворот\n- Статусы\n- История\n- Доступы")
    ContainerDb(monitoring_db, "Monitoring DB", "TimescaleDB", "Данные наблюдения\n- Видео\n- Датчики\n- События")
    ContainerDb(device_db, "Device DB", "PostgreSQL", "Реестр устройств\n- Метаданные\n- Состояния\n- Конфигурации")
}

System_Ext(device_gateway, "Device Gateway", "Шлюз для подключения устройств\n- MQTT/HTTP/WebSocket\n- Аутентификация устройств\n- Маршрутизация данных")

System_Ext(heating_device, "Устройство отопления", "Датчик температуры, термостат")
System_Ext(lighting_device, "Устройство освещения", "Умная лампа, выключатель")
System_Ext(gate_device, "Устройство ворот", "Датчик положения")
System_Ext(monitoring_device, "Устройство наблюдения", "Камера, датчик движения")

' Взаимодействия пользователя
Rel(user, web_app, "Использует", "HTTPS")

' Прямые взаимодействия веб-приложения с микросервисами
Rel(web_app, heating_s, "API вызовы", "HTTP/REST\n/api/heating/*")
Rel(web_app, lighting_s, "API вызовы", "HTTP/REST\n/api/lighting/*") 
Rel(web_app, gate_s, "API вызовы", "HTTP/REST\n/api/gate/*")
Rel(web_app, monitoring_s, "API вызовы", "HTTP/REST\n/api/monitoring/*")
Rel(web_app, device_s, "API вызовы", "HTTP/REST\n/api/device/*")

' Взаимодействия с отдельными БД
Rel(heating_s, heating_db, "Чтение/запись", "SQL")
Rel(lighting_s, lighting_db, "Чтение/запись", "SQL")
Rel(gate_s, gate_db, "Чтение/запись", "SQL")
Rel(monitoring_s, monitoring_db, "Чтение/запись", "SQL")
Rel(device_s, device_db, "Чтение/запись", "SQL")

' Взаимодействия между сервисами
Rel(heating_s, device_s, "Получает информацию об устройствах", "HTTP/REST")
Rel(lighting_s, device_s, "Получает информацию об устройствах", "HTTP/REST")
Rel(gate_s, device_s, "Получает информацию об устройствах", "HTTP/REST")
Rel(monitoring_s, device_s, "Получает информацию об устройствах", "HTTP/REST")

' Подключение устройств к Device Gateway
Rel(heating_device, device_gateway, "Подключается", "MQTT/HTTP")
Rel(lighting_device, device_gateway, "Подключается", "MQTT/HTTP")
Rel(gate_device, device_gateway, "Подключается", "MQTT/HTTP")
Rel(monitoring_device, device_gateway, "Подключается", "MQTT/HTTP")

' Device Gateway подключается к Device Service
Rel(device_gateway, device_s, "Регистрация и управление\nустройствами", "HTTP/REST/gRPC")

@enduml

**Диаграмма компонентов (Components)**

@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

Package "Heating Service" {
    
    Component(heating_controller, "HeatingController", "REST API", "Обрабатывает HTTP запросы\n- Установка температуры\n- Получение статуса\n- Управление режимами")
    
    Component(temp_monitor, "TemperatureMonitor", "Компонент мониторинга", "Мониторинг температуры\n- Сбор данных с датчиков\n- Оповещения")
    
    Component(schedule_manager, "ScheduleManager", "Компонент расписаний", "Управление расписаниями\n- Графики обогрева\n- Сезонные настройки")
    
    Component(device_manager, "DeviceManager", "Управление устройствами", "Объединяет:\n- Взаимодействие с устройствами\n- Команды термостатам\n- Клиент Device Service")
    
    Component(heating_repository, "HeatingRepository", "Компонент объединение работы в БД", "Необходим для работы с Бд не напрямую через каждый сервис")
    
    Component(notification_client, "NotificationClient", "Клиент уведомлений", "Отправка уведомлений\n- Температурные предупреждения\n- Оповещения о неисправностях")
}

ContainerDb(heating_db, "Heating DB", "PostgreSQL")

' Прямые взаимодействия контроллера с компонентами
Rel(heating_controller, temp_monitor, "Запрос данных температуры", "Direct call")
Rel(heating_controller, schedule_manager, "Управление расписанием", "Direct call")
Rel(heating_controller, device_manager, "Управление устройствами", "Direct call")
Rel(heating_controller, heating_repository, "Сохранение/чтение данных", "Direct call")
Rel(heating_controller, notification_client, "Отправка уведомлений", "Direct call")

' Взаимодействия между компонентами
Rel(temp_monitor, heating_controller, "Отправка оповещений", "Event/callback")
Rel(temp_monitor, heating_repository, "Сохраняет метрики", "Repository call")
Rel(device_manager, heating_repository, "Сохраняет статус устройств", "Repository call")
Rel(schedule_manager, device_manager, "Отправка команд по расписанию", "Direct call")

' Взаимодействие с БД
Rel(heating_repository, heating_db, "Чтение/запись", "JDBC")

' Внешние системы
Component(device_service, "Device Service", "Внешний сервис")
Rel(device_manager, device_service, "API вызовы", "HTTP/REST")

@enduml


**Диаграмма кода (Code)**

Диаграмма кода, в виде диаграммы последовательности компонента расписаний 

@startuml
title Диаграмма последовательности ScheduleManager 

actor "Пользователь" as User
participant "HeatingController" as HC
participant "ScheduleManager" as SM
participant "HeatingRepository" as HR
participant "DeviceManager" as DM
database "Heating DB" as DB
participant "Device Service" as DS

== Сценарий: Создание нового расписания обогрева ==

User -> HC: POST /api/heating/schedule\n{temp: 22°C, days: [Mon-Fri], time: "08:00"}
activate HC

HC -> SM: createSchedule(scheduleData)
activate SM

SM -> HR: validateSchedule(scheduleData)
activate HR
HR --> SM: validationResult
deactivate HR

alt Валидация не пройдена
    SM --> HC: Error: Invalid schedule
    HC --> User: 400 Bad Request
else Валидация успешна
    SM -> HR: saveSchedule(scheduleData)
    activate HR
    HR -> DB: INSERT INTO schedules (...)
    activate DB
    DB --> HR: success
    deactivate DB
    HR --> SM: scheduleId
    deactivate HR
    
    SM --> HC: scheduleCreated(scheduleId)
    deactivate SM
    
    HC -> DM: applySchedule(scheduleId)
    activate DM
    DM -> DS: GET /api/devices?type=thermostat
    activate DS
    DS --> DM: devicesList
    deactivate DS
    DM -> DM: calculateDeviceCommands(devicesList, scheduleData)
    DM --> HC: scheduleApplied
    deactivate DM
    
    HC --> User: 201 Created\n{scheduleId: "123"}
end
deactivate HC

== Сценарий: Обновление существующего расписания ==

User -> HC: PUT /api/heating/schedule/123\n{temp: 21°C, time: "07:30"}
activate HC

HC -> SM: updateSchedule(scheduleId, newData)
activate SM

SM -> HR: getSchedule(scheduleId)
activate HR
HR -> DB: SELECT * FROM schedules WHERE id=scheduleId
activate DB
DB --> HR: scheduleData
deactivate DB
HR --> SM: schedule
deactivate HR

SM -> HR: updateSchedule(scheduleId, newData)
activate HR
HR -> DB: UPDATE schedules SET ... WHERE id=scheduleId
activate DB
DB --> HR: success
deactivate DB
HR --> SM: updated
deactivate HR

SM -> DM: updateActiveSchedule(scheduleId, newData)
activate DM
DM -> DS: POST /api/devices/commands/batch\n{commands: [...]}
activate DS
DS --> DM: commandsProcessed
deactivate DS
DM --> SM: updated
deactivate DM

SM --> HC: scheduleUpdated
deactivate SM

HC --> User: 200 OK\n{message: "Schedule updated"}
deactivate HC

@enduml

# Задание 3. Разработка ER-диаграммы

@startuml
title ER-диаграмма для экосистемы "Умный дом"

entity "User" as user {
  *id : UUID <<PK>>
  --
  *email : VARCHAR(255)
  *name : VARCHAR(100)
  created_at : TIMESTAMP
}

entity "House" as house {
  *id : UUID <<PK>>
  --
  *user_id : UUID <<FK>>
  *name : VARCHAR(100)
  address : TEXT
}

entity "Device" as device {
  *id : UUID <<PK>>
  --
  *house_id : UUID <<FK>>
  *type : VARCHAR(50)
  *name : VARCHAR(100)
  serial_number : VARCHAR(100)
  status : VARCHAR(20)
  configuration : JSON
}

entity "Telemetry" as telemetry {
  *id : UUID <<PK>>
  --
  *device_id : UUID <<FK>>
  *metric : VARCHAR(50)
  *value : DECIMAL(10,4)
  *timestamp : TIMESTAMP
}

entity "Schedule" as schedule {
  *id : UUID <<PK>>
  --
  *device_id : UUID <<FK>>
  *action : VARCHAR(50)
  *time : TIME
  days : VARCHAR(20)
  is_active : BOOLEAN
}

entity "Scenario" as scenario {
  *id : UUID <<PK>>
  --
  *user_id : UUID <<FK>>
  *name : VARCHAR(100)
  conditions : JSON
  actions : JSON
  is_active : BOOLEAN
}

' Связи между сущностями
user ||--o{ house : "owns"
house ||--o{ device : "has"
device ||--o{ telemetry : "generates"
device ||--o{ schedule : "has"
user ||--o{ scenario : "creates"

note right of device
  type: 'heating', 'lighting', 
  'gate', 'monitoring'
end note

note right of telemetry
  metric: 'temperature', 'power', 
  'state', 'motion', etc.
end note

@enduml


# Задание 4. Создание и документирование API

### 1. Тип API

В моих endpoint используется только REST API.

### 2. Документация API

openapi: 3.0.0
info:
  title: Умный дом - API микросервисов
  description: |
    REST API для взаимодействия между микросервисами системы "Умный дом"
    Соответствует ER-диаграмме системы
  version: 1.0.0

tags:
  - name: Devices
    description: Управление устройствами (Device Service)
  - name: Telemetry
    description: Телеметрия устройств
  - name: Schedules
    description: Расписания устройств
  - name: Scenarios
    description: Сценарии автоматизации

paths:
  # ========== DEVICE SERVICE ENDPOINTS ==========
  /devices:
    get:
      tags: [Devices]
      summary: Получение списка устройств дома
      parameters:
        - name: houseId
          in: query
          required: true
          schema:
            type: string
            format: uuid
        - name: type
          in: query
          schema:
            type: string
            enum: [heating, lighting, gate, monitoring]
      responses:
        '200':
          description: Список устройств
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Device'
    
    post:
      tags: [Devices]
      summary: Регистрация нового устройства
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/DeviceCreate'
      responses:
        '201':
          description: Устройство зарегистрировано
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Device'

  /devices/{deviceId}:
    get:
      tags: [Devices]
      summary: Получение информации об устройстве
      parameters:
        - name: deviceId
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Информация об устройстве
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Device'
    
    patch:
      tags: [Devices]
      summary: Обновление устройства
      parameters:
        - name: deviceId
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/DeviceUpdate'
      responses:
        '200':
          description: Устройство обновлено

  # ========== TELEMETRY ENDPOINTS ==========
  /telemetry:
    post:
      tags: [Telemetry]
      summary: Отправка телеметрии устройства
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/TelemetryCreate'
      responses:
        '201':
          description: Телеметрия сохранена

  /telemetry/{deviceId}:
    get:
      tags: [Telemetry]
      summary: Получение телеметрии устройства
      parameters:
        - name: deviceId
          in: path
          required: true
          schema:
            type: string
            format: uuid
        - name: metric
          in: query
          schema:
            type: string
        - name: from
          in: query
          schema:
            type: string
            format: date-time
        - name: to
          in: query
          schema:
            type: string
            format: date-time
      responses:
        '200':
          description: Список телеметрии
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Telemetry'

  # ========== SCHEDULE ENDPOINTS ==========
  /schedules:
    post:
      tags: [Schedules]
      summary: Создание расписания для устройства
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ScheduleCreate'
      responses:
        '201':
          description: Расписание создано
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Schedule'

  /schedules/{deviceId}:
    get:
      tags: [Schedules]
      summary: Получение расписаний устройства
      parameters:
        - name: deviceId
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Список расписаний
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Schedule'

  /schedules/{scheduleId}:
    delete:
      tags: [Schedules]
      summary: Удаление расписания
      parameters:
        - name: scheduleId
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: Расписание удалено

  # ========== SCENARIO ENDPOINTS ==========
  /scenarios:
    get:
      tags: [Scenarios]
      summary: Получение сценариев пользователя
      parameters:
        - name: userId
          in: query
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Список сценариев
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Scenario'
    
    post:
      tags: [Scenarios]
      summary: Создание сценария
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ScenarioCreate'
      responses:
        '201':
          description: Сценарий создан
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Scenario'

  /scenarios/{scenarioId}:
    put:
      tags: [Scenarios]
      summary: Обновление сценария
      parameters:
        - name: scenarioId
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ScenarioUpdate'
      responses:
        '200':
          description: Сценарий обновлен

components:
  schemas:
    # ========== DEVICE SCHEMAS ==========
    Device:
      type: object
      properties:
        id:
          type: string
          format: uuid
        house_id:
          type: string
          format: uuid
        type:
          type: string
          enum: [heating, lighting, gate, monitoring]
        name:
          type: string
          maxLength: 100
        serial_number:
          type: string
          maxLength: 100
        status:
          type: string
          maxLength: 20
        configuration:
          type: object
          additionalProperties: true
      example:
        id: "550e8400-e29b-41d4-a716-446655440000"
        house_id: "550e8400-e29b-41d4-a716-446655440001"
        type: "heating"
        name: "Термостат гостиная"
        serial_number: "TH-001-2024"
        status: "online"
        configuration: {"target_temp": 23.5, "mode": "auto"}

    DeviceCreate:
      type: object
      required:
        - house_id
        - type
        - name
      properties:
        house_id:
          type: string
          format: uuid
        type:
          type: string
          enum: [heating, lighting, gate, monitoring]
        name:
          type: string
          maxLength: 100
        serial_number:
          type: string
          maxLength: 100
        configuration:
          type: object
          additionalProperties: true

    DeviceUpdate:
      type: object
      properties:
        name:
          type: string
          maxLength: 100
        status:
          type: string
          maxLength: 20
        configuration:
          type: object
          additionalProperties: true

    # ========== TELEMETRY SCHEMAS ==========
    Telemetry:
      type: object
      properties:
        id:
          type: string
          format: uuid
        device_id:
          type: string
          format: uuid
        metric:
          type: string
          maxLength: 50
        value:
          type: number
          format: float
        timestamp:
          type: string
          format: date-time
      example:
        id: "550e8400-e29b-41d4-a716-446655440002"
        device_id: "550e8400-e29b-41d4-a716-446655440000"
        metric: "temperature"
        value: 22.5
        timestamp: "2024-01-15T14:30:00Z"

    TelemetryCreate:
      type: object
      required:
        - device_id
        - metric
        - value
      properties:
        device_id:
          type: string
          format: uuid
        metric:
          type: string
          maxLength: 50
        value:
          type: number
          format: float
        timestamp:
          type: string
          format: date-time

    # ========== SCHEDULE SCHEMAS ==========
    Schedule:
      type: object
      properties:
        id:
          type: string
          format: uuid
        device_id:
          type: string
          format: uuid
        action:
          type: string
          maxLength: 50
        time:
          type: string
          format: time
        days:
          type: string
          maxLength: 20
        is_active:
          type: boolean
      example:
        id: "550e8400-e29b-41d4-a716-446655440003"
        device_id: "550e8400-e29b-41d4-a716-446655440000"
        action: "set_temperature"
        time: "18:00:00"
        days: "mon,tue,wed,thu,fri"
        is_active: true

    ScheduleCreate:
      type: object
      required:
        - device_id
        - action
        - time
      properties:
        device_id:
          type: string
          format: uuid
        action:
          type: string
          maxLength: 50
        time:
          type: string
          format: time
        days:
          type: string
          maxLength: 20
        is_active:
          type: boolean
          default: true

    # ========== SCENARIO SCHEMAS ==========
    Scenario:
      type: object
      properties:
        id:
          type: string
          format: uuid
        user_id:
          type: string
          format: uuid
        name:
          type: string
          maxLength: 100
        conditions:
          type: object
        actions:
          type: object
        is_active:
          type: boolean
      example:
        id: "550e8400-e29b-41d4-a716-446655440004"
        user_id: "550e8400-e29b-41d4-a716-446655440005"
        name: "Вечерний режим"
        conditions: {"time": "18:00", "presence": true}
        actions: {"lights": "dim", "temperature": 21}
        is_active: true

    ScenarioCreate:
      type: object
      required:
        - user_id
        - name
      properties:
        user_id:
          type: string
          format: uuid
        name:
          type: string
          maxLength: 100
        conditions:
          type: object
        actions:
          type: object
        is_active:
          type: boolean
          default: true

    ScenarioUpdate:
      type: object
      properties:
        name:
          type: string
          maxLength: 100
        conditions:
          type: object
        actions:
          type: object
        is_active:
          type: boolean

    # ========== ERROR SCHEMA ==========
    Error:
      type: object
      properties:
        errorCode:
          type: string
        message:
          type: string
        timestamp:
          type: string
          format: date-time

  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

security:
  - BearerAuth: []

# Задание 5. Работа с docker и docker-compose

Перейдите в apps.

Там находится приложение-монолит для работы с датчиками температуры. В README.md описано как запустить решение.

Вам нужно:

1) сделать простое приложение temperature-api на любом удобном для вас языке программирования, которое при запросе /temperature?location= будет отдавать рандомное значение температуры.

Locations - название комнаты, sensorId - идентификатор названия комнаты

```
	// If no location is provided, use a default based on sensor ID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// If no sensor ID is provided, generate one based on location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}
```

2) Приложение следует упаковать в Docker и добавить в docker-compose. Порт по умолчанию должен быть 8081

3) Кроме того для smart_home приложения требуется база данных - добавьте в docker-compose файл настройки для запуска postgres с указанием скрипта инициализации ./smart_home/init.sql

Для проверки можно использовать Postman коллекцию smarthome-api.postman_collection.json и вызвать:

- Create Sensor
- Get All Sensors

Должно при каждом вызове отображаться разное значение температуры

Ревьюер будет проверять точно так же.


# **Задание 6. Разработка MVP**

Необходимо создать новые микросервисы и обеспечить их интеграции с существующим монолитом для плавного перехода к микросервисной архитектуре. 

### **Что нужно сделать**

1. Создайте новые микросервисы для управления телеметрией и устройствами (с простейшей логикой), которые будут интегрированы с существующим монолитным приложением. Каждый микросервис на своем ООП языке.
2. Обеспечьте взаимодействие между микросервисами и монолитом (при желании с помощью брокера сообщений), чтобы постепенно перенести функциональность из монолита в микросервисы. 

В результате у вас должны быть созданы Dockerfiles и docker-compose для запуска микросервисов. 
