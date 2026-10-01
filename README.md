# nagorneva_reviews

Сервис отзывов докторов на Go. Архитектура повторяет подход `medblogers_base`: composition root находится в `internal/app.go` и `internal/init.go`, HTTP-маршруты тонкие, бизнес-логика лежит в `internal/module/<module>/action/<action>`, а запросы PostgreSQL — рядом в `dal.go`.

## Модули

- `reviews` — публичные методы карты и фильтра отзывов без авторизации.
- `admin` — CRUD сущностей, просмотр отзывов врача, загрузка фото и JWT-вход. Все методы `/api/v1/admin/*`, кроме `/login`, требуют `Authorization: Bearer <access_token>`.

`personal_data_consent` не удаляет фактические данные врача из БД. Он управляет их выдачей в публичном `GetDoctorsOnMap`: без согласия API возвращает пустые `full_name` и `photo`.

## Быстрый старт

```bash
cp .env.example .env
docker compose up -d
goose -dir migrations postgres "$DATABASE_URL" up
go run ./cmd/nagorneva_reviews
```

Перед первым запуском задайте `JWT_SECRET` длиной не менее 32 символов. При указании `ADMIN_EMAIL` и `ADMIN_PASSWORD` приложение создаст единственного администратора при запуске — пароль сохраняется только в bcrypt-хеше.

## API

Контракты находятся в `api/reviews/reviews.proto` и `api/admin/admin.proto`. Для генерации клиента, gRPC-кода и gateway-кода:

```bash
buf dep update
buf generate
```

Основные HTTP методы:

| Метод | Назначение |
|---|---|
| `POST /api/v1/doctors_on_map` | Доктора по координатам, с городом, доступным фото и активными отзывами/курсами |
| `POST /api/v1/reviews/filter` | Фильтр отзывов по `course_id`, `lat`, `lon`, `radius_km` |
| `POST /api/v1/admin/login` | Вход администратора |
| `/api/v1/admin/{users,cities,specialties,doctors,courses,reviews}` | `GET`, `POST`, `GET/{id}`, `PATCH/{id}`, `DELETE/{id}` |
| `GET /api/v1/admin/doctors/{id}/reviews` | Все отзывы выбранного врача, включая неактивные |
| `POST /api/v1/admin/doctors/{id}/photo` | `multipart/form-data`, поле `photo`; загрузка в S3 |

Для S3 совместимого MinIO используются `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, `S3_BUCKET` и `S3_PUBLIC_URL`. Имя объекта — `doctors/<doctor_id>/<timestamp>.<ext>`.

## Важное решение по данным

В исходном ТЗ `speciallity` содержит опечатку. В БД и API использовано согласованное имя `specialties` / `specialty_id`. Это предотвращает постоянную неоднозначность в SQL и фронтенде. Координаты врача сохранены отдельно от города: город задаёт фильтр и справочник, а координаты врача — точку на карте.
