# Notes Service
Небольшой сервис, дающий возможность создать заметку, получить список всех заметок и просмотреть конкретную заметку

### Запуск
Для запуска сервиса выполните один из скриптов в папке `deployment`:

**Linux/macOS:**
```bash
./deployment/deploy.sh
```

**Windows (PowerShell):**
```powershell
.\deployment\deploy.ps1
```

## API

### Базовый URL
`http://127.0.0.1:8000`

### Основные эндпоинты:
- **POST** `/note` - Создание новой заявки
- **GET** `/note/{id}` - Получение заявки по ID
- **GET** `/notes` - Получение списка всех заявок

Полная документация API доступна в файле [docs/openapi.yaml](docs/openapi.yaml)
