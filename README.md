# go-steps

Учебный проект!

Адрес `/delivery/calculate` - рассчитает доставку для транспортной компании с двух складов отправки. 

Для локальной проверки вызвать:
```
curl http://localhost:8080/delivery/calculate?type=1&GateId=656008&weight=2.350&CountryIso=643&OrderAmount=1202.30
```

Пакеты:
- echo
- go-playground

### CI

- добавлена MultiStage сборка