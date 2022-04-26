# Сервер

В папке Server представлен код реализации gRPC прокси-сервиса с двумя обработчиками.

Первый обработчик пополняет баланс указанного пользователя на указанную сумму, принимает на вход 2 значения:
1. ид юзера (или баланса), с которого списание;
2. сумма средств для перевода.

Второй переводит указанную сумму со счета первого пользователя на счет другого, при этом не может уйти в минус, принимает на вход 3 значения:
1. ид юзера (или баланса), с которого списание;
2. ид юзера, кому на счет поступают средства;
3. сумма средств для перевода.

СУБД - PostgreSQL.

Предусмотрен останов веб-сервера без потери обрабатываемых запросов.

Сервер принемает флаги при запуске:

* **addr**              по умолчанию: *localhost:50051*   - адрес сервера в формате host:port

* **port**              по умолчанию: *50051*             - порт сервера.

* **databaseHost**      по умолчанию: *localhost*         - адрес сервера базы данных.

* **databasePort**      по умолчанию: *5432*              - порт сервера базы данных.

* **databaseUser**      по умолчанию: *postgres*          - имя пользователя базы данных.

* **databasePassword**  по умолчанию  *(env) $pgxPass*    - пароль базы данных.

* **databaseName**     по умолчанию:  *postgres*          - имя базы данных.


# Интерфейс взаимодействия

Интерфейс взаимодействия клиента и сервера описан в файле proto/payment.proto. 

# СУБД

PostgreSQL.

# Автор

Никита Шеремета
guamoko95@gmail.com