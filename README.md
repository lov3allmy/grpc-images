## Задание

Необходимо написать сервис на Golang работающий по gRPC.

Сервис должен:

1. Принимать бинарные файлы (изображения) от клиента и сохранять их на жесткий
   диск.
2. Иметь возможность просмотра списка всех загруженных файлов в формате:
    ```
    Имя файла | Дата создания | Дата обновления
    ```
3. Отдавать файлы клиенту.
4. Ограничивать количество одновременных подключений с клиента:
    - на загрузку/скачивание файлов - 10 конкурентных запросов;
    - на просмотр списка файлов - 100 конкурентных запросов.

## Запуск сервера

```shell
make run && make build
```

## RPCs

### Загрузка нового изображения на сервер
```protobuftext
UploadImage(UploadImageRequest)
message UploadImageRequest {
  string name = 1;
  bytes data = 2;
}
```
### Загрузка и замена изображения на сервере
```protobuftext
UpdateImage(UpdateImageRequest)

message UpdateImageRequest {
  string id = 1;
  string name = 2;
  bytes data = 3;
}
```
### Скачивание изображения с сервера
```protobuftext
DownloadImage(DownladImageRequest)

message UpdateImageRequest {
  string id = 1;
}
```
### Просмотр списка загруженных файлов
```protobuftext
GetImagesList(GetImagesListRequest)

message GetImagesListResponse {
  message ImageInfo {
    string name = 1;
    google.protobuf.Timestamp created_at = 2;
    google.protobuf.Timestamp modified_at = 3;
  }

  repeated ImageInfo imageInfo = 1;
}
```
