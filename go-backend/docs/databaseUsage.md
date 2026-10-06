- 1 Set up the database with creds

- 2.1 Create the database schema and queries required (database/query.sql, database/schema.sql)

- 2.2 Create migration scripts using goose (database/migrations/...)

- 3 In main.go connect to the database 

```
ConnectDB(ctx context.Context) (*pgxpool.Pool, error)
```

- 4 Create repository

```
repo interface ->contains Methods that will be used internally

-- THIS IS MOSTLY A TEMPLATE
type sqlcRepository struct {
	queries *database.Queries
}

-- constructor for the interface
func newRepository(queries *database.Queries) repository {
	return sqlcRepository{queries: queries}
}

func (r sqlcRepository) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return r.queries.CreateUser(ctx, arg)
}

type svc struct {
	repo repository
}
```

- 5.1 Consume the repo in service

```go
func newService(repo repository) *svc {
	return &svc{repo: repo}
}
```
- 5.2 Use the repo Method in the service

```go
func (s *svc) SignUp(ctx context.Context, request signUpRequest) (database.User, error) {
return s.repo.CreateUser(ctx, database.CreateUserParams{
		Firstname: pgtype.Text{String: request.FirstName, Valid: true},
		Lastname:  pgtype.Text{String: request.LastName, Valid: true},
		Username:  request.Username,
		Password:  pgtype.Text{String: string(passwordHash), Valid: true},
	})
}
```

- 6 Create the endpoint and use the service method