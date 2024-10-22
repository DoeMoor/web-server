tools:

  postgres:
    dbName: chirpy.

  goose:
   for migration
   dir to migration files
   ```
   -dir ./sql/schema/
   ```
  sqlc:
   for db connection
   ```yaml
   sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
        emit_json_tags: true
   ```


.env
```
DB_URL="postgres://user:passwd@localhost:port/chirpy?sslmode=disable"
PLATFORM="dev" or "prod"
SECRET="generated"
POLKA_KEY="also generated"
```

build:
go build -o ./_build/out && ./_build/out  //build and run