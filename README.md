# nagorneva_reviews

Сервис отзывов докторов на Go. Сборка приложения находится в `internal/app.go` и `internal/init.go`, транспорт — в `internal/app/server`.

HTTP JSON → grpc-gateway / Clay descriptor → interceptor → API → action → DAL → PostgreSQL.

Нативный gRPC → Clay gRPC-сервер → interceptor → тот же API.

Как в `medblogers_base`, Clay регистрирует нативный gRPC и HTTP grpc-gateway по сгенерированным service descriptors. Gateway вызывает сервис внутри процесса; отдельное сетевое gRPC-соединение между HTTP и обработчиками не требуется. Прямые gRPC-запросы и JSON-запросы проходят одну проверку JWT через `WithGRPCUnaryMiddlewares`. Multipart-загрузка фото остаётся отдельным HTTP-адаптером с проверкой того же токена.

## Архитектура

```text
internal/
  app/api/                       # Преобразование protobuf ↔ DTO
  app/server/                    # HTTP, grpc-gateway, gRPC и shutdown
  app/middleware/                # JWT и CORS
  config/                        # Конфигурация из окружения
  modules/
    admin/
      module.go                  # Сборка зависимостей
      action/login/              # Вход и отдельный DAL
      client/s3/                 # Внешнее хранилище
      entities/
        aggregator.go
        doctors/                 # Аналогично: users, cities, specialties, courses, reviews
          aggregator.go
          action/create_doctor/  # Сценарий и небольшой интерфейс Storage
          action/update_doctor/
          action/upload_doctor_photo/
          dto/                   # Типизированные входы и результаты
          dal/                   # SQL конкретной сущности
        users/service/password/  # Проверка и хеширование пароля
    reviews/
      action/<operation>/        # Публичный сценарий
        dal/                     # Его SQL
      domain/                    # Модели публичного каталога
  pkg/postgres/                  # Интерфейс DB, чтение строк, перевод ошибок БД
  pkg/domain_error/              # Ошибки, не зависящие от PostgreSQL и транспорта
```

Actions не импортируют pgx и не строят SQL. Общего CRUD с именами таблиц и `map[string]any` нет: каждая сущность имеет собственный репозиторий. Пакет `dto` задаёт контракт с actions и API; PostgreSQL-типы остаются в DAL. Интерфейсы зависимостей определены рядом с actions. Слой `service` используется для самостоятельной бизнес-логики, а простые операции вызывают DAL напрямую.

Необязательные поля PATCH представлены указателями: отсутствие поля сохраняет значение, а `false`, `0` и пустая строка применяются явно. Создание и обновление врача вместе с чтением связанных справочников выполняются одним SQL-выражением. Репозитории принимают `postgres.DB`, который реализуют и пул, и транзакция pgx; отдельный механизм транзакций пока не требуется.

Создание начального администратора вынесено в `users/action/ensure_initial_user`. Повторный запуск не меняет существующий пароль, а конкурентная вставка защищена `ON CONFLICT`. Существующие ограничения и внешние ключи БД сохранены.

## Модули

- `reviews` — публичные методы карты и фильтра отзывов без авторизации.
- `admin` — CRUD сущностей, просмотр отзывов врача, загрузка фото и JWT-вход. Все методы `/api/v1/admin/*`, кроме `/login`, требуют `Authorization: Bearer <access_token>`.

`personal_data_consent` не удаляет фактические данные врача из БД. Он управляет их выдачей в публичном `GetDoctorsOnMap`: без согласия API возвращает пустые `full_name` и `photo`.

## Быстрый старт

```bash
cp .env.example .env
docker compose up -d
make migrate-up
go run ./cmd/nagorneva_reviews
```

Перед первым запуском задайте `JWT_SECRET` длиной не менее 32 символов. При указании `ADMIN_EMAIL` и `ADMIN_PASSWORD` приложение создаст администратора с этим email при запуске, если его ещё нет — пароль сохраняется только в bcrypt-хеше.

HTTP API работает через grpc-gateway на `HTTP_ADDR` (по умолчанию `:8080`), а нативный gRPC API — на `GRPC_ADDR` (по умолчанию `:7002`). Оба транспорта используют контракты из `api/`.

`DEBUG=false` по умолчанию. При `DEBUG=true` JWT-проверка отключается для admin HTTP и native gRPC вызовов — используйте это только локально.

Для браузерных запросов задайте `CORS_ALLOWED_ORIGINS` как список точных origin через запятую, например `https://reviews.example.ru,https://admin.example.ru`. По умолчанию разрешены только `http://localhost:3000` и `http://localhost:5173`; значения с путём, query-параметрами или фрагментом игнорируются.

## API

Контракты находятся в `api/reviews/reviews.proto` и `api/admin/admin.proto`. Для генерации Go-сообщений, gRPC-клиентов/сервера, gateway и OpenAPI:

```bash
make buf
```

Версии плагинов зафиксированы в `buf.gen.yaml`, зависимости proto — в `buf.lock`. Команда не обновляет их автоматически. Go-код, OpenAPI и Clay descriptors (`*.pb.goclay.go`) генерируются в `gen/go`. Для локального плагина `protoc-gen-goclay` выполните `make deps`. Clay объединяет OpenAPI-схемы сервисов при запуске.

Swagger UI доступен по `/docs`, объединённая схема — по `/swagger.json`. Имена RPC и существующие HTTP URL сохранены. JSON gateway использует стандартный protobuf-формат: camelCase в ответах, `int64` как строки; во входных запросах допустимы также имена полей из proto в snake_case.

Ошибки отсутствующей записи преобразуются в gRPC `NotFound` / HTTP 404, дубликата — в `AlreadyExists` / 409, внешнего ключа — в `FailedPrecondition` / 400, ограничений данных — в `InvalidArgument` / 400. Внутренние ошибки не раскрывают SQL клиенту.

Основные HTTP методы:

| Метод | Назначение |
|---|---|
| `POST /api/v1/doctors_on_map` | Доктора по координатам, с городом, доступным фото и активными отзывами/курсами |
| `POST /api/v1/doctors/filter` | Фильтр врачей по `city_ids`, `specialty_ids`, `course_ids`, `is_active`, `personal_data_consent`, `reviews_sort` |
| `POST /api/v1/reviews/filter` | Фильтр отзывов по `course_id`, `lat`, `lon`, `radius_km` |
| `POST /api/v1/admin/login` | Вход администратора |
| `/api/v1/admin/{users,cities,specialties,doctors,courses,reviews}` | `GET`, `POST`, `GET/{id}`, `PATCH/{id}`, `DELETE/{id}` |
| `GET /api/v1/admin/doctors/{id}/reviews` | Все отзывы выбранного врача, включая неактивные |
| `POST /api/v1/admin/doctors/{id}/photo` | `multipart/form-data` с полем `photo` либо JSON protobuf с фото в base64; загрузка в S3 |

Загрузка фотографий временно отключена: в `.env` и `.env.example` оставьте `S3_BUCKET=` пустым. В этом режиме S3-клиент не создаётся, ключи хранилища не требуются, а multipart-загрузка фото возвращает HTTP 503 (gRPC — `FailedPrecondition`). Остальные функции приложения работают без S3.

Чтобы вернуть загрузку фотографий через Yandex Object Storage по схеме из `medblogers_base`, укажите существующий бакет и его статические ключи в `S3_ENDPOINT`, `S3_REGION`, `S3_ACCESS_KEY`, `S3_SECRET_KEY` и `S3_BUCKET`. Приложение не создаёт бакет само. `S3_PUBLIC_URL` необязателен: по умолчанию ссылки имеют вид `https://storage.yandexcloud.net/<S3_BUCKET>/<object-key>`. Имя объекта — `images/user_<doctor_id>_<timestamp>.<ext>`.

Для этого проекта выбран отдельный бакет `nagorneva-reviews-photos`. Параметры подключения хранятся в локальном `.env` (не попадает в Git); в `.env.example` бакет и ключи оставлены пустыми. После изменения `.env` перезапустите backend. Если `S3_BUCKET` задан в окружении процесса, очистите его и там: переменные окружения имеют приоритет над `.env`. Для отображения фото на публичной карте бакет должен разрешать анонимный `s3:GetObject` для `images/*`; публичный список объектов и публичная запись не нужны. Без `S3_BUCKET` загрузка отключена, а при указанном бакете без ключей приложение вернёт ошибку при запуске.

При первичной настройке создайте бакет в Yandex Cloud и примените политику из `scripts/s3-public-photos-policy.json`. Сервисному аккаунту, которому принадлежат статические ключи, нужен доступ к загрузке объектов в новый бакет. Если создание через S3 API возвращает `AccessDenied`, используйте аккаунт с правом создания бакетов. Существующие бакеты `medblogers_base` для этого проекта не используются.

## Локальная проверка публичной карты

Локальная копия виджета находится в `dev/map/reviews-map-widget.js`. Тестовая страница использует `api-url="http://localhost:8080"` без атрибута `demo` и показывает журнал запросов к backend.

```bash
python3 -m http.server 5173 --bind 127.0.0.1 --directory dev/map
```

Откройте `http://localhost:5173`. Backend должен работать на порту `8080` и иметь доступ к своей базе PostgreSQL. Для Яндекс Карт и их дополнительных модулей требуется интернет. После обновления скачанного файла в корне проекта скопируйте его в `dev/map/reviews-map-widget.js`.

Пример подключения локального виджета:

```html
<nagorneva-reviews-map
  api-url="http://localhost:8080"
  yandex-api-key="YOUR_YANDEX_MAPS_API_KEY"
></nagorneva-reviews-map>
<script src="http://localhost:5173/reviews-map-widget.js"></script>
```

## Важное решение по данным

В исходном ТЗ `speciallity` содержит опечатку. В БД и API использовано согласованное имя `specialties` / `specialty_id`. Это предотвращает постоянную неоднозначность в SQL и фронтенде. Координаты врача сохранены отдельно от города: город задаёт фильтр и справочник, а координаты врача — точку на карте.

В `POST /api/v1/doctors/filter` значения внутри одного массива работают как «или», а непустые группы фильтров объединяются как «и». Например, врач из любого указанного города, с любой указанной специальностью и прошедший любой указанный курс. Поле `course_ids` сверяется с отдельной связью `doctor_courses`, а не с текстом отзывов. Необязательные `is_active` и `personal_data_consent` принимают `true` или `false`; если поле не передано, соответствующий признак не фильтруется. `reviews_sort` принимает `REVIEWS_SORT_DESC` или `REVIEWS_SORT_ASC`; без него и при одинаковом числе отзывов врачи отсортированы по ID. Ответ содержит `doctor_count` — общее число врачей, совпавших с фильтром.


## Проверки

```bash
make build
make test-race
make lint
TEST_DATABASE_URL='postgres://reviews_test:test@localhost:55493/reviews_test?sslmode=disable' make test-integration
```

`make deps` устанавливает Buf, Goose, protoc-gen-goclay и golangci-lint v2 в `bin/`. Интеграционные тесты требуют отдельную тестовую PostgreSQL: создают временную схему, применяют миграции и удаляют только эту схему после выполнения. Отсутствие `TEST_DATABASE_URL` считается ошибкой запуска интеграционных тестов.

Для локальной тестовой БД:

```bash
docker run --rm -d --name reviews-test-db \
  -p 127.0.0.1:55493:5432 \
  -e POSTGRES_USER=reviews_test -e POSTGRES_PASSWORD=test \
  -e POSTGRES_DB=reviews_test postgres:15
# После тестов:
docker stop reviews-test-db
```

CI выполняет сборку, тесты с детектором гонок, интеграционные тесты PostgreSQL и линтер. Транспортные тесты поднимают HTTP и gRPC на случайных локальных портах и проверяют авторизацию обоих транспортов, CORS, JSON/multipart-маршруты и большие фото. Интеграционные тесты проверяют типизированный CRUD, частичные обновления, связи и ограничения БД, публичные фильтры и скрытие персональных данных.
