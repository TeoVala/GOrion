# GOrion v0.2.523

**Quick API Creation Toolkit**

### Please note that this project is currently a work in progress, with many features and improvements actively being developed.

GOrion is a streamlined Go project thats aims rapid API development. It provides a simple and efficient setup for programmers to quickly create APIs, leveraging the power of `chi` for routing, `gorm` for the ORM and `gorm gen` to generate relations (also some sql magic to get the db relationships and feed them to `gorm`). Gorion focuses on generating a functional API based on your database schema, including handling relationships and providing a simple intuitive way to generate code templates for routes and handlers (more to be middleware, models etc).

**Version:** v0.2.523

## Features

*   **Rapid API Generation:** Quickly generate API endpoints based on your database structure.
*   **`chi` Router:** Utilizes the lightweight and flexible `chi` router for HTTP request handling.
*   **`gorm` ORM:** Seamless integration with `gorm` for powerful database interactions.
*   **`gorm gen` Integration:** Automates code generation from your database schema, including model and repository code.
*   **Automatic Relationship Handling:** Fetches and generates code for database relationships using `gorm`.
*   **Template Autogeneration:** Uses templates to streamline the generation of API code.
*   **Simple Project Setup:** Designed for ease of use and quick project initialization.

**Note:** The project doesn't provide any way to generate database relations and tables it focuses solely on generating APIs from existing database schemas.

## Technologies Used

*   **`net/http`:** Go's standard HTTP library.
*   **`github.com/go-chi/chi/v5`:** A lightweight, idiomatic, and composable router for building HTTP services in Go.
*   **`gorm.io/gorm`:** A fantastic ORM library for Go, making database interactions easy.
*   **`gorm.io/gen`:** A powerful tool for generating Go code from GORM models and database schemas.

## Getting Started

### Prerequisites

*   Go (Tried on version 1.24)


### Installation

```bash
go get github.com/your_github_username/gorion # Replace with your actual repository path
```

### Usage

1.  **Configure your database connection:** Update your configuration file (or environment variables) with your database connection details.
2.  **Define your routes:** GOrion uses chi to define the routes following the [Chi documentations](https://go-chi.io/#/pages/routing)
3.  **Define your GORM models:** Create your GORM models that reflect your database tables. Under the hood it uses `gorm gen` and some sql call magic to generate the structs for each table.
4.  **Run the generation command:** Execute the GOrion command to generate the API code based on your models and database schema. (`go run manage.go genmodels`).
5.  **Start the server:** Run your Go application to start the generated API server.



## Acknowledgements
GOrion Uses many open source libraries give them some love they are amazing:

 - [Chi](https://go-chi.io/)
 - [The amazing Gorm](https://gorm.io/)

## Run Commands
Currently GOrion uses "go run" to run all the commands, this is bound to change in the following versions to a CLI which is better for now you can run commands this way.

Generate route from template

```bash
  go run manage.go make:route -name:test

```

Generate handler from template

```bash
  go run manage.go make:handler -name:test

```

Auto generate models

```bash
  go run manage.go genmodels
```

Start the server

```bash
  go run manage.go runserver
```
Optionally you can add port (you can also define a default one inside the .env file): 
`
--port: portNum
`

#### Feel free to contact me with any suggestions or feedback: [contact@teovala.com](mailto:contact@teovala.com)
