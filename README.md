# ZBOT
Телеграм-бот, обеспечивающий отправку сообщений о состоянии хостов, контролируемых  системой мониторинга Zabbix.

## Установка и запуск 
1. Установить docker и docker compose.
Инструкция по установке доступна в официальной документации
2. В папке с проектом создать файл `.env`, содержащий токен вашего телеграм-бота и учетные данные пользователя в системе Zabbix.<br>
_Пример:_
```
BOT_TOKEN= your_bot's_token
ZABBIX_URLAPI= https://example.com/zabbix/api_jsonrpc.php # ваш веб-интерфейс Zabbix 
ZABBIX_USERNAME=  your_username
ZABBIX_PASSWORD= your_password
ZABBIX_WILDCARDSHOSTS=*hostanme* #маска названий хостов
ZABBIX_SLEEP = 5 # интервал (в минутах) между опросами Zabbix API
```
3. В папке с проектом выполнить команду
```commandline
docker compose up
```
## Применение
### Регистрация пользователя
```commandline
docker exec zbt /app/zbot users -add username_from_the_telegram-bot
```
### Просмотр списка зарегистрированных пользователей
```commandline
docker exec zbt /app/zbot users -list
```
### Удаление зарегистрированного пользователя 
```commandline
docker exec zbt /app/zbot users -del  username_from_the_telegram-bot
```
### Взаимодействие с ботом посредством мессенджера.
1. После перехода в окно диалога с ботом, для аутентификации пользователя, выполните команду
```
/start
```
2. Просмотр списка доступных команд
```
/help
```
3. Просмотр списка проблемных хостов
```
/listp
```
4.  Просмотр списка восстановленных хостов
```
/listr
```
> Примечание!    
Обновление сведений о состоянии хостов производиться с заданным интервалом (ZABBIX_SLEEP). По умолчанию - пять минут.  
> Изменение состояния хостов сопровождается сообщением.  
> Пример:  
> The number of problematic hosts is 40 (1 new, 0 changed).  
> The number of restored hosts is 1.
```
Пример выполнения команды /listp:
Host name: Хост_1, problems:[Нет пинга]  
ch_Host name: Хост_22, problems:[Disk space is critically low]  
Host name: Хост_67:[Нет пинга]  
new_Host name:  Хост_34, problems:[Нет пинга Disk space is critically low]
The number of problematic hosts is 4 (1 new, 1 changed).
```
>Примечание!  
>Значение префиксов:  
> ch_ - изменился перечень проблем (относительно предыдущего обновления сведений);  
> new_ - новый проблемный хост (относительно предыдущего обновления сведений);  
