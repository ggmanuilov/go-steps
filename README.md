# go-steps

Учебный/боевой проект!

Адрес `/delivery/calculate` - рассчитает доставку для транспортной компании с двух складов отправки. 

Для локальной проверки вызвать:
```
curl "http://localhost:8080/delivery/calculate?DeliveryType=1&PvzId=656008&PointId=1&CountryIso=643&Weight=2.350&OrderAmount=1202.30"
```

Пакеты:
- echo
- go-playground

### Тесты

Контракт API задают приёмочные тесты (`tests/acceptance`): они собирают приложение через реальную фабрику и реестр маршрутов, поэтому проверяют и проводку `/delivery/calculate`:

```
make test
```

### CI

- добавлена MultiStage сборка
