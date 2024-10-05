# Список изменений:

1. Перенесена точка входа в приложение и изменён Dockerfile (TODO)
2. Добавлено более подробное логгирование, улучшены сообщения об ошибках.
3. Изменена структура проекта (см. `README.md`).
4. Используются следующие статусы:
    - `User`
        - `Get by name`
            - `200 OK`: пользователь найден
            - `204 NoContent`: такого пользователя не существует
            - `401 Unauthorized`: ошибка авторизации
            - `500 InternalServerError`: ошибка сервера
        - `register`
            - `201 Created` успешная регистрация
            - `400 BadRequest`: переданы некорректные данные
            - `409 Conflict`: пользователь уже зарегистрирован
            - `500 InternalServerError`: ошибка сервера
        - `login`
            - `200 OK`: авторизация успешна
            - `400 BadRequest`: переданы некорректные данные
            - `401 Unauthorized`: неверные логин/пароль или пользователя не существует
            - `500 InternalServerError`: ошибка сервера
    - `Role`
        - `Create`
            - `201 Created` роль успешно создана
            - `400 BadRequest`: переданы некорректные данные
            - `409 Conflict`: роль уже существует
            - `500 InternalServerError`: ошибка сервера
        - `Get by ID`
            - `200 OK`: роль найдена
            - `204 NoContent`: такой роли не существует
            - `400 BadRequest`: некорректный ID или запрос без параметра id
            - `401 Unauthorized`: ошибка авторизации
            - `500 InternalServerError`: ошибка сервера
        - `Get by name`
            - `200 OK`: роль найдена
            - `204 NoContent`: такой роли не существует
            - `400 BadRequest`: запрос без параметра name
            - `401 Unauthorized`: ошибка авторизации
            - `500 InternalServerError`: ошибка сервера
    - `Site`
        - `Create`
            - `201 Created` роль успешно создана
            - `400 BadRequest`: переданы некорректные данные
            - `500 InternalServerError`: ошибка сервера
        - `Get`
            - `200 OK`: site найден
            - `204 NoContent`: такого site не существует
            - `400 BadRequest`: некорректный ID или запрос без параметра id
            - `401 Unauthorized`: ошибка авторизации
            - `500 InternalServerError`: ошибка сервера
        - `Get list`
            - `200 OK`: sites найдены
            - `204 NoContent`: таких sites не существует
            - `400 BadRequest`: невозможно получить ID пользователя из JWT
            - `401 Unauthorized`: ошибка авторизации
            - `500 InternalServerError`: ошибка сервера
        - `Patch update`
            - `200 OK`: site успешно обновлён
            - `400 BadRequest`: переданы некорректные данные
            - `401 Unauthorized`: ошибка авторизации
            - `404 NotFound`: такого sites не существует
            - `500 InternalServerError`: ошибка сервера
        - `Soft delete`
            - `200 OK`: site успешно удалён
            - `400 BadRequest`: некорректный ID или запрос без параметра id
            - `401 Unauthorized`: ошибка авторизации
            - `404 NotFound`: такого site не существует
            - `409 Conflict`: site уже удалён
            - `500 InternalServerError`: ошибка сервера
        - `Soft delete`
            - `200 OK`: site успешно восстановлен
            - `400 BadRequest`: некорректный ID или запрос без параметра id
            - `401 Unauthorized`: ошибка авторизации
            - `404 NotFound`: такого sites не существует
            - `409 Conflict`: site уже восстановлен
            - `500 InternalServerError`: ошибка сервера
    - `Building`
        - `Create`
            - `201 Created`: здание успешно создано
            - `400 BadRequest`: переданы некорректные данные
            - `401 Unauthorized`: ошибка авторизации
            - `500 InternalServerError`: ошибка сервера
        - `Get`
            - `200 OK`: здание найдено
            - `204 NoContent`: такого здания не существует
            - `400 BadRequest`: некорректный ID или запрос без параметра id
            - `401 Unauthorized`: ошибка авторизации
            - `500 InternalServerError`: ошибка сервера
        - `Get list`
            - `200 OK`: здания найдены
            - `204 NoContent`: таких зданий не существует
            - `400 BadRequest`: некорректный ID или запрос без параметра id
            - `401 Unauthorized`: ошибка авторизации
            - `500 InternalServerError`: ошибка сервера
        - `Patch update`
            - `200 OK`: здание успешно обновлено
            - `400 BadRequest`: переданы некорректные данные
            - `401 Unauthorized`: ошибка авторизации
            - `404 NotFound`: такого здания не существует
            - `500 InternalServerError`: ошибка сервера
        - `Soft delete`
            - `200 OK`: здание успешно удалено
            - `400 BadRequest`: некорректный ID или запрос без параметра id
            - `401 Unauthorized`: ошибка авторизации
            - `404 NotFound`: такого здания не существует
            - `409 Conflict`: building уже удалено
            - `500 InternalServerError`: ошибка сервера
        - `Soft delete`
            - `200 OK`: здание успешно восстановлено
            - `400 BadRequest`: некорректный ID или запрос без параметра id
            - `401 Unauthorized`: ошибка авторизации
            - `404 NotFound`: такого здания не существует
            - `409 Conflict`: здение уже восстановлено
            - `500 InternalServerError`: ошибка сервера
    - `Floor`: аналогично Building
    - `WallType`: аналогично Building
    - `Wall`: аналогично Building
    - `AccessPointType`: аналогично Building
    - `AccessPointRadioTemplate`: аналогично Building
    - `AccessPoint`: аналогично Building
    - `AccessPointRadio`: аналогично Building
    - `SensorType`: аналогично Building
    - `Sensor`: аналогично Building
        
5. Добавлены/обновлены примеры ответов сервера в Postman для слудющих сценариев:
    - `Health`
        - `health`
            - `success`: статус 200 (добавлен)
            - `fail`: статус 500 (добавлен)
    - `Public`
        - `get picture`
            - `success`: статус 200 (добавлен)
            - `fail`: статус 404 (добавлен)
    - `User`
        - `Get by name`:
            - `success`: статус 200 (добавлен)
            - `not found`: статус 204 (добавлен)
            - `unauthorized`: статус 401 (добавлен)
        - `register`
            - `success`: статус 201 (добавлен)
            - `fail bad request`: статус 400 (добавлен)
            - `fail already exists`: статус 409 (добавлен)
            - `fail server error`: статус 500 (добавлен)
        - `login`
            - `success`: статус 200 (добавлен)
            - `fail bad request`: статус 400 (добавлен)
            - `fail bad login`: статус 401 (добавлен)
    - `Role`
        - `Create`
            - `success`: статус 201 (добавлен)
            - `fail already exists`: статус 409 (добавлен)
        - `Get by ID`
            - `success`: статус 200 (добавлен)
            - `not found`: статус 204 (добавлен)
            - `fail wrong id`: статус 400 (добавлен)
            - `fail bad request`: статус 404 (добавлен)
        - `Get by name`
            - `success`: статус 200 (добавлен)
            - `not found`: статус 204 (добавлен)
            - `fail bad request`: статус 404 (добавлен)
    - `Site`
        - `Create`
            - `success 1`: статус 201 (обновлён)
            - `success 2`: статус 201 (добавлен)
            - `fail bad request`: статус 404 (добавлен)
        - `Get`
            - `success`: статус 200 (обновлён)
            - `not found`: статус 204 (добавлен)
            - `fail wrong id`: статус 400 (добавлен)
        - `Get list`
            - `success`: статус 200 (обновлён)
        - `Patch update`
            - `success`: статус 200 (обновлён)
            - `fail not updated`: статус 400 (добавлен)
            - `fail not found`: статус 404 (добавлен)
            - `fail bad request`: статус 404 (добавлен)
        - `Soft delete`
            - `success`: статус 200 (обновлён)
            - `fail not found`: статус 404 (добавлен)
            - `fail already deleted`: статус 409 (обновлён)
            - `fail server error`: статус 500 (добавлен)
        - `Restore`
            - `success`: статус 200 (обновлён)
            - `fail wrong id`: статус 400 (добавлен)
            - `fail not found`: статус 404 (добавлен)
            - `fail already restored`: статус 409 (обновлён)
    - `Building`
        - `Create`
            - `success`: статус 201 (обновлён)
            - `fail unexisting site id`: статус 400 (добавлен)
            - `fail bad request`: статус 400 (добавлен)
        - `Get`
            - `success`: статус 200 (обновлён)
        - `Get list`
            - `success`: статус 200 (обновлён)
            - `not found`: статус 204 (добавлен)
        - `Soft delete`
            - `success`: статус 200 (обновлён)
            - `fail already deleted`: статус 409 (обновлён)
        - `Restore`
            - `success`: статус 200 (обновлён)
            - `fail already restored`: статус 409 (обновлён)
            - `fail not found`: статус 404 (добавлен)
        - `Patch update`
            - `success`: статус 200 (обновлён)
            - `fail not found`: статус 404 (добавлен)
    - `Floor`
        - `Create`
            - `success`: статус 201 (обновлён)
        - `Get`
            - `success`: статус 200 (обновлён)
        - `Get list`
            - `success`: статус 200 (обновлён)
        - `Soft delete`
            - `success`: статус 200 (обновлён)
            - `fail not found`: статус 404 (добавлен)
            - `fail already deleted`: статус 409 (обновлён)
        - `Restore`
            - `success`: статус 200 (обновлён)
            - `fail already restored`: статус 409 (обновлён)
        - `Patch update`
            - `success`: статус 200 (обновлён)
    - `WallType`
        - `Create`
            - `success`: статус 201 (обновлён)
            - `fail unexisting site id`: статус 400 (добавлен)
        - `Get`
            - `success`: статус 200 (обновлён)
            - `not found`: статус 204 (добавлен)
        - `Get list`
            - `success`: статус 200 (обновлён)
        - `Soft delete`
            - `success`: статус 200 (обновлён)
            - `fail already deleted`: статус 409 (обновлён)
        - `Restore`
            - `success`: статус 200 (обновлён)
            - `fail already restored`: статус 409 (обновлён)
        - `Patch update`
            - `success`: статус 200 (обновлён)
    - `Wall`
        - `Create`
            - `success`: статус 201 (обновлён)
        - `Get`
            - `success`: статус 200 (обновлён)
        - `Get detailed`
            - `success`: статус 200 (обновлён)
        - `Get list`
            - `success`: статус 200 (обновлён)
        - `Soft delete`
            - `success`: статус 200 (обновлён)
        - `Restore`
            - `success`: статус 200 (обновлён)
        - `Patch update`
            - `success`: статус 200 (обновлён)
    `AccessPointType`: аналогично Wall
    `AccessPointRadioTemplate`: аналогично Wall (без detailed)
    `AccessPoint`: аналогично Wall (+ get list detailed)
    `AccessPointRadio`: аналогично Wall (без detailed)
    `SensorType`: аналогично Wall (без detailed)
    `Sensor`: аналогично Wall (radios всегда null)
6. Пинг базы данных вынесен в репозиторий.
7. Добавлен слой middleware (проверка авторизации)
8. Создана единая точка маршрутизации.
9. Добавлены маршруты для `Role`: создание, получение по ID, получение по названию.
10. Восстановлены маршруты для `Site`: создание, получение по ID, получение списка, мягкое удаление, восстановление, patch update.
11. Восстановлены маршруты для `Floor`: создание, получение по ID, получение списка, мягкое удаление, восстановление, patch update.
12. Для `Floor` добавлена возможность **создавать** следующие поля: `cell_size_meter`, `north_area_indent_meter`, `south_area_indent_meter`, `west_area_indent_meter`, `east_area_indent_meter`. Указанные поля имеют значение по умолчанию или могут быть NULL, т.е. их указание необязательно при создании таблицы.
13. Для `Floor` добавлена возможность **обновлять** следующие поля: `cell_size_meter`, `north_area_indent_meter`, `south_area_indent_meter`, `west_area_indent_meter`, `east_area_indent_meter`.
14. Для таблицы `floors` полю `scale` добавлено **значение по умолчанию – 0.1**.
15. Восстановлены маршруты для `WallType`: создание, получение по ID, получение списка, мягкое удаление, восстановление, patch update.
16. Восстановлены маршруты для `Wall`: создание, получение по ID, получение детализированного объекта по ID, получение списка, мягкое удаление, восстановление, patch update.
17. Для `Wall` добавлена возможность **обновлять** следующие поля: `wall_type_id`.
18. Восстановлены маршруты для `AccessPointType`: создание, получение по ID, получение детализированного объекта по ID, получение списка, мягкое удаление, восстановление, patch update.
19. Восстановлены маршруты для `AccessPointRadioTemplate`: создание, получение по ID, получение списка, мягкое удаление, восстановление, patch update.
20. Для `AccessPointRadioTemplate` добавлена возможность **создавать/обновлять** следующее поле: `channel2`. Данное поле не является обязательным.
21. Маршрут `/radioTeamplate` **заменён** на `/ap-radio-template`.
22. Восстановлены маршруты для `AccessPoint`: создание, получение по ID, получение списка, получение детализированного объекта по ID, получения списка детализированных объектов,  мягкое удаление, восстановление, patch update.
23. Для `AccessPointRadio` добавлена возможность **создавать/обновлять** следующие поля: `channel2`, `is_active`.
24. Маршрут `/radio` **заменён** на `/ap-radio`.
25. Добавлена необязательная для заполнения колонка `color` для таблиц `access_points` и `sensors`.
26. Восстановлены маршруты для `SensorType`: создание, получение по ID, получение списка,  мягкое удаление, восстановление, patch update.
27. Восстановлены маршруты для `Sensor`: создание, получение по ID, получение списка, получение детализированного объекта по ID, получения списка детализированных объектов,  мягкое удаление, восстановление, patch update. **На текущий момент поле `radios` всегда `null`!**
28. Добавлены новые баги :3

## Миграции:

Обновление колонки `scale` тыблицы `floors`:
```sql
ALTER TABLE floors
ALTER COLUMN scale SET DEFAULT 0.1;
```

Обновление таблицы `access_point_radio_templates`:
```sql
ALTER TABLE access_point_radio_templates
ADD COLUMN channel2 SMALLINT CHECK (channel2 > 0),
ALTER COLUMN channel_width TYPE VARCHAR(32),
ALTER COLUMN channel_width SET NOT NULL,
DROP CONSTRAINT IF EXISTS access_point_radio_templates_channel_width_check;
```
Обновление таблицы `access_point_radios`:
```sql
ALTER TABLE access_point_radios
ADD COLUMN channel2 SMALLINT CHECK (channel2 > 0),
ALTER COLUMN channel_width TYPE VARCHAR(32),
ALTER COLUMN channel_width SET NOT NULL,
DROP CONSTRAINT IF EXISTS access_point_radios_channel_width_check;
```
Обновление таблицы `sensor_radio_templates`:
```sql
ALTER TABLE sensor_radio_templates
ADD COLUMN channel2 SMALLINT CHECK (channel2 > 0),
ALTER COLUMN channel_width TYPE VARCHAR(32),
ALTER COLUMN channel_width SET NOT NULL,
DROP CONSTRAINT IF EXISTS sensor_radio_templates_channel_width_check;
```
Обновление таблицы `sensor_radios`:
```sql
ALTER TABLE sensor_radios
ADD COLUMN channel2 SMALLINT CHECK (channel2 > 0),
ALTER COLUMN channel_width TYPE VARCHAR(32),
ALTER COLUMN channel_width SET NOT NULL,
DROP CONSTRAINT IF EXISTS sensor_radios_channel_width_check;
```

Добавление колонки `is_active` в `access_point_radio_templates` и `sensor_radio_templates`
```sql
-- Для таблицы access_point_radio_templates
ALTER TABLE access_point_radio_templates
ADD COLUMN is_active BOOLEAN;

UPDATE access_point_radio_templates
SET is_active = FALSE;

ALTER TABLE access_point_radio_templates
ALTER COLUMN is_active SET NOT NULL;


-- Для таблицы sensor_radio_templates
ALTER TABLE sensor_radio_templates
ADD COLUMN is_active BOOLEAN;

UPDATE sensor_radio_templates
SET is_active = FALSE;

ALTER TABLE sensor_radio_templates
ALTER COLUMN is_active SET NOT NULL;

```

Добавление колонки `color` в `access_points` и `sensors`:
```sql
ALTER TABLE access_points
ADD COLUMN color CHAR(6);

ALTER TABLE sensors
ADD COLUMN color CHAR(6);
```

Изменение таблицы `sensor_types`:
```sql
ALTER TABLE sensor_types
DROP COLUMN interface_0;
ALTER TABLE sensor_types
DROP COLUMN interface_1;
ALTER TABLE sensor_types
DROP COLUMN interface_2;

ALTER TABLE sensor_types
DROP COLUMN diagram;

ALTER TABLE sensor_types
DROP COLUMN correction_factor_24;
ALTER TABLE sensor_types
DROP COLUMN correction_factor_5;
ALTER TABLE sensor_types
DROP COLUMN correction_factor_6;

ALTER TABLE sensor_types
DROP COLUMN rx_ant_gain;

ALTER TABLE sensor_types
DROP COLUMN hor_rotation_offset;
ALTER TABLE sensor_types
DROP COLUMN vert_rotation_offset;

ALTER TABLE sensor_types
DROP COLUMN alias;

ALTER TABLE sensor_types
ADD COLUMN model VARCHAR(128);
UPDATE sensor_types
SET model = 'model';
ALTER TABLE sensor_types
ALTER COLUMN model SET NOT NULL;

ALTER TABLE sensor_types
ADD COLUMN z FLOAT;
UPDATE sensor_types
SET z = 1;
ALTER TABLE sensor_types
ALTER COLUMN z SET NOT NULL;
```

Изменение таблицы `sensors`:
```sql
ALTER TABLE sensors
DROP COLUMN interface_0,
DROP COLUMN interface_1,
DROP COLUMN interface_2;

ALTER TABLE sensors
DROP COLUMN alias;
```

