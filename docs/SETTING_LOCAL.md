# Настройка локальной машины
* ОС - Windows 10 (настройка и запуск производится в консоле git-bash), Ubuntu
### 1. Установка
Скачиваем файл dep.exe(для linux dep) и перемещаем его в директорию /usr/bin/
```bash
wget https://github.com/evgeny-klyopov/dep/releases/download/v1.0.6/dep.windows-amd64.exe.tar.gz
tar -xvf dep.windows-amd64.exe.tar.gz
```
Название файла не важно
```bash
mv dep.exe /usr/bin/dep.exe
```
Или установка через go
```bash
go get github.com/evgeny-klyopov/dep
```
### 2. Настройка файла конфигурации для деплоя
В папке проекта необходимо создать файл deploy.json, в котором описать свои настройки. Примеры конфигураций находятся в папке example. 