# Настройка локальной машины
* ОС - Windows 10 (настройка и запуск производится в консоле git-bash), Ubuntu
### 1. Установка
Скачиваем файл deploy.exe(для linux deploy) и перемещаем его в директорию /usr/bin/
```bash
wget https://github.com/evgeny-klyopov/dep/raw/master/deploy.exe
```
Название файла не важно
```bash
mv dep.exe /usr/bin/dep
```
Или установка через go
```bash
go get github.com/evgeny-klyopov/dep
```
### 2. Настройка файла конфигурации для деплоя
В папке проекта необходимо создать файл deploy.json, в котором описать свои настройки. Примеры конфигураций находятся в папке example. 