# Sismos-api

The project designed for educational purposes to learn development technologies and practices. The application is packaged with Docker for easy deployment.


## Tech Stack

**Server:**   
- **Lenguage:** Golang
- **ORM:** Gorm
- **Framework:** Gin
- **Database:** MySQL
- **Containerization:** Docker


## API Reference

#### Get all Feature

```http
  GET /api/features
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `mag_type[]`   | `string` |  ["ml", "md"]         |
| `per_page`     | `string` |  10                   |
| `page`         | `string` |  1                    |




#### Get Feature

```http
  GET /api/${id_feature}/feature
```

#### UpDate Feature

```http
  PUT /api/${id_feature}/feature
```

| Body      | Type     | Value                    |
| :--------    | :------- | :------------------------------|
| `event_id`        | `string`  | **Required**. nc74032466       |
| `magnitude`       | `int`     | **Required** 1           | 
| `place`           | `string`  | **Required**. 1 km N of The Geysers, CA      |
| `time`            | `string`  | **Required**. 2024-04-10T22:59:40Z           |
| `tsunami`         | `bool`    | **Required**. false           |
| `mag_type`        | `string`  | **Required**. md           |
| `external_url`    | `string`  | **Required**. https://earthquake.usgs.gov/earthquakes/eventpage/nc74032466           |
| `title`           | `string`  | **Required**. M 0.8 - 0 km N of The Geysers, CA           |
| `longitude`       | `int`     | **Required**. -123           |
| `latitude`        | `int`     | **Required**. 39           |


#### Get Comment

```http
  GET /api/${id_feature}/comments
```

#### Create taskComments

```http
  POST /api/${id_feature}/comment
```

| Body          | Type     | Description                    |
| :--------     | :------- | :------------------------------|
| `body`        | `string` | **Required**. comment          |




