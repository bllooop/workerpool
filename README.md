Реализован примитивный worker-pool c возможностью динамически добавлять и удалять воркеры. Входные данные (строки) поступают в канал, воркеры их обрабатывают - выводят на экран номер воркера и сами данные.

### Запуск кода
Необходимо ввести для запуска команду:

```
go run main.go
```
Далее, будет возможность ввести желаемое количество worker и задач, которые эти worker и обработают. Строки генерируются по формату "задача №".
После запуска кода будет выводится работа worker.
Для удаления воркеров можно прописать в main pool.removeWork() и будет вызван метод удаления последнего воркера в списке.
