# ТЗ Golang developer
Зарегистрироваться на GetBlock.io, используя API GetBlock создать сервис,
выводящий адрес, баланс которого изменился (в любую сторону) больше остальных
за последние сто блоков.  
Получить номер последнего блока можно с помощью следующего метода:  
https://getblock.io/docs/available-nodes-methods/ETH/JSON-RPC/eth_blockNumber/  
А данные блока вместе с транзакциями через
https://getblock.io/docs/available-nodes-methods/ETH/JSON-RPC/eth_getBlockByNumber/  
Важно: API ключи в репозитории не хранить  
## Start
+ Сервис работает через http сервер.  
+ Адрес порт заданы по умолчанию: localhost:8080  
+ Endpoint (apikey указывается в адресной строке):
find_address/YOUR_API_KEY (альтернотивное решение было передача ключа через аргумент os.Args[1])
+ example http://localhost:8080/find_address/...
