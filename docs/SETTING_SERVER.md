# Настройка сервера
Настройка сервера включает три основных этапа:
### 1. Создание пользователя, под которым будет осуществляется деплой (Можно использовать уже существующего пользователя)
### 2. Настройка ssh ключей и авторизация по ним
### 3. Раздача прав на управление различными сервисами, если необходимо
В данном примере настройка производилась на сервере Ubuntu

### 1. Создание пользователя для деплоя
```bash
adduser deployer
```
Добавляем в группу, если требуется 
```bash
usermod -a -G www-data deployer
```
Установка прав на директорию проекта
```bash
chown deployer:www-data /path/to/dir
```
### 2. Генерация ssh ключа
```bash
ssh-keygen -t rsa -b 4096 -C "mail@klepov.info"
```
Добавление авторизации через ключи
Необходимо в файл /home/deployer/.ssh/authorized_keys добавить публичный ключ с локальной машины. Если файл authorized_keys отсутствует, создать. После добавления с локальной машины при соединение по ssh пароль не требуется.
### 3. Права на управление различными сервисами, если необходимо
Для добавления прав необходимо отредактировать файл конфигурации
```bash
sudo visudo
```
Добавляем в конец файла
```bash
Cmnd_Alias YOUR_SERVICE_NAME_YOUR_COMMAND = /usr/sbin/service YourServiceName YourCommand
# No-Password Commands
YourUserDeployerCommand ALL=NOPASSWD: YOUR_SERVICE_NAME_YOUR_COMMAND
```
, где
YOUR_SERVICE_NAME_YOUR_COMMAND - название вашей константы
YourServiceName - ваш сервис
YourCommand - команда для управления, например рестарт
YourUserDeployerCommand - пользователь из под которого происходит деплой

Пример, релоад php-fpm
```bash
Cmnd_Alias PHP_RELOAD = /usr/sbin/service php7.2-fpm reload
# No-Password Commands
deployer ALL=NOPASSWD: PHP_RELOAD
```