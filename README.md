# blog aggregator

- DATABASE : gator (postgresql)

  - install

  ```
  sudo apt update
  sudo apt install postgresql postgresql-contrib
  ```

  - set a pwd if you linux/wsl

  ```
  sudo passwd postgres
  ```

- LANGUAGE : golang

  - `https://go.dev/doc/instal1`

- Migration : goose

  - install

  ```
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```

  - how to connect

  ```
  goose postgres <connection_string> up
  ```

- sql generate (sqlc)

  ```
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  ```

  - add a config (sqlc.yaml)

  ```
  version: "2"
  sql:
    - schema: "sql/schema"
      queries: "sql/queries"
      engine: "postgresql"
      gen:
        go:
          out: "internal/database"import _ "github.com/lib/pq"
  ```

  - how to connect postgresql

  ```
  go get github.com/lib/pq
  ```

  - and then write this in main.go

  ```
  import _ "github.com/lib/pq"
  ```
