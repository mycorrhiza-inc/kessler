# Kessler-go


# Building docs for ingest 

run this command and it should update properly 

```bash
swag init -g cmd/ingest/main.go -o cmd/ingest/docs
```

# Datababase Migrations

We use [`goose`](https://pressly.github.io/goose) for our migrations.

Migrations can be found in the `migrations` directory.
