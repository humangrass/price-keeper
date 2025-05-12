# Price Keeper

## Getting price from Jupiter

```bash
go run ./cmd/jupiter/main.go --jupiter-token-id So11111111111111111111111111111111111111112
```

## Keeper

```bash
go run ./cmd/keeper/main.go
```

## Migrations

```bash
go run ./cmd/migrate/main.go --config-file keeper.example.yaml --help
```

```bash
go run ./cmd/migrate/main.go --config-file keeper.example.yaml up-to --help
```

## Swagger

```bash
swag init -g cmd/keeper/main.go --output docs --parseDependency --parseInternal
```
