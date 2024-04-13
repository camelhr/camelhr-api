## CamelHR API

This is the api service of camelhr application built with golang.

It provides ReST apis for managing the hrms data.

### Local Development Setup

#### Prerequisites

* Go version `1.22` or above
* Docker Desktop

#### Start the application using docker

This application is configured to run with docker compose. Follow the steps below to start, restart or stop the application.

* Start or restart the application

  ```shell
  make up
  ```

  This will automatically build the project from source code and start the application & database containers in docker.
  If you make any changes to the code you can just rerun this command so that it will rebuild the project and update the respective container.

* Shutdown application using

  ```shell
  make down
  ```

* Shutdown application and clear database

  ```shell
  make nuke
  ```

#### Run tests & linters

* Run tests

  ```shell
  make test
  ```

* Run lint

  ```shell
  make lint
  ```

### Contribution Guidelines

> Every Contribution Makes a Difference

Read the [Contribution Guidelines](CONTRIBUTING.md) before you contribute.
