# go-steps

Учебный/боевой проект!

Адрес `/delivery/calculate` - рассчитает доставку для транспортной компании с двух складов отправки. 

Для локальной проверки вызвать:
```
curl http://localhost:8080/delivery/calculate?Type=1&GateId=656008&Weight=2.350&CountryIso=643&OrderAmount=1202.30
```

Пакеты:
- echo
- go-playground

### CI

- добавлена MultiStage сборка

### Кодегенерация OpenAPI

```
 oapi-codegen -generate="types" -package delivery ./openapi/delivery.yml > internal/generated/Types.gen.go
```

### Источники

- https://github.com/deepmap/oapi-codegen
- https://gobyexample.com.ru/