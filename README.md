### Commands
```sh
go version
```

```sh
go run ./cmd/api
```

```sh
go get github.com/julienschmidt/httprouter@v1
```

```sh
go install github.com/rakyll/hey@latest
```

```sh
go get github.com/lib/pq@v1
```

```sh
go get golang.org/x/time/rate@latest
```

```sh
go get golang.org/x/crypto/bcrypt@latest
```

```sh
go get github.com/go-mail/mail/v2@v2
```

```sh
export GREENLIGHT_DB_DSN='postgres://greenlight:pa55word@localhost:5432/greenlight?sslmode=disable'
```

```sh
echo $GREENLIGHT_DB_DSN
```

```sh
curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
```

```sh
$ apt-get update
```

```sh
$ apt-get install -y migrate
```

```sh
migrate -version
```

```sh
migrate create -seq -ext=.sql -dir=./migrations create_movies_table
```

```sh
migrate create -seq -ext=.sql -dir=./migrations add_movies_check_constraints
```

```sh
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
```

```sh
$ migrate -path=./migrations -database=$EXAMPLE_DSN goto 1
```

```sh
migrate -path=./migrations -database=$EXAMPLE_DSN up
```

```sh
migrate -path=./migrations -database=$EXAMPLE_DSN down
```

```sh
migrate create -seq -ext=.sql -dir=./migrations create_users_table
```

```sh
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
```

```sh
migrate create -seq -ext .sql -dir ./migrations create_tokens_table
```

```sh
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
```

// Creating database GIN indexes for tables on certain fields to avoid full table scans
```sh
migrate create -seq -ext .sql -dir ./migrations add_movies_indexes'
```
```sh
migrate -path ./migrations -database $GREENLIGHT_DB_DSN up
```

### Fixing SQL Migration Errors
- Investigate the original error and figure out if the migration file which failed was partially applied.
- Then you need to manually roll-back the partially applied migration.
- Once that’s done, then you must also ‘force’ the version number in the schema_migrations table to the correct value.
```sh
migrate -path=./migrations -database=$EXAMPLE_DSN force 1
```
- Once you force the version, the database is considered ‘clean’ and you should be able to run migrations again without any problem.

### Commands
```sh
curl localhost:4000/v1/healthcheck
```

```sh
curl -X POST localhost:4000/v1/movies
```

```sh
curl localhost:4000/v1/movies/123
```

```sh
BODY='{"title":"Moana","year":2016,"runtime":107, "genres":["animation","adventure"]}'
```

```sh
BODY='{"title":"Black Panther","year":2018,"runtime":134,"genres":["action","adventure"]}'
```

```sh
BODY='{"title":"Deadpool","year":2016, "runtime":108,"genres":["action","comedy"]}'
```

```sh
BODY='{"title":"The Breakfast Club","year":1986, "runtime":96,"genres":["drama"]}'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/movies
```

```sh
BODY='{"title":"Black Panther","year":2018,"runtime":134,"genres":["sci-fi","action","adventure"]}'
```

```sh
curl -X PUT -d "$BODY" localhost:4000/v1/movies/2
```

```sh
curl -X PATCH -d '{"year": 1985}' localhost:4000/v1/movies/4
```

```sh
xargs -I % -P8 curl -X PATCH -d '{"runtime":97}' "localhost:4000/v1/movies/4" < <(printf '%s\n' {1..8})
```

```sh
curl localhost:4000/v1/movies
```

// When using curl to send a request containing more than one query string
// parameter, you must wrap the URL in quotes for it to work correctly.
```sh
curl "localhost:4000/v1/movies?title=moana&genres=animation,adventure&page=1&page_size=5&sort=year"
```

```sh
curl localhost:4000/v1/movies
```

```sh
curl "localhost:4000/v1/movies?title=black+panther"
```

```sh
curl "localhost:4000/v1/movies?genres=adventure"
```

```sh
curl "localhost:4000/v1/movies?title=moana&genres=animation,adventure"
```

```sh
curl "localhost:4000/v1/movies?genres=western"
```

// Return all movies where the title includes the case-insensitive word 'panther'0.
```sh
curl "localhost:4000/v1/movies?title=panther"
```

// Return all movies where the title includes the case-insensitive words 'the' and'club'
```sh
curl "localhost:4000/v1/movies?title=the+club"
```

```sh
curl "localhost:4000/v1/movies?sort=-title"
```

```sh
curl "localhost:4000/v1/movies?sort=-runtime"
```

```sh

```sh
for i in {1..6}; do curl http://localhost:4000/v1/healthcheck; done
```

```sh
go run ./cmd/api/ -limiter-enabled=false
```

```sh
curl "localhost:4000/v1/movies?page_size=2"
```

```sh
curl "localhost:4000/v1/movies?page_size=2&page=2"
```

```sh
curl "localhost:4000/v1/movies?page_size=2&page=2"
```

```
BODY='{"name": "Alice Smith", "email": "alice@example.com", "password": "pa55word"}'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/users
```

```sh
BODY='{"name": "Godae Hill", "email": "godae@example.com", "password": "pa55word"}'
```

```sh
curl -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users
```

```sh
BODY='{"name": "Carol Smith", "email": "carol@example.com", "password": "pa55word"}'
```

```sh
curl -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users
```

```sh
BODY='{"name": "Dave Smith", "email": "dave@example.com", "password": "pa55word"}'
```

```sh
curl -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users
```

```sh
BODY='{"name": "Faith Smith", "email": "faith@example.com", "password": "pa55word"}'
```

```sh
curl -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users
```

```sh
curl -X PUT -d '{"token": "ZYGQTPU5PKKJRY7SFOAMKXPGQY"}' localhost:4000/v1/users/activated
```

```sh
BODY='{"email": "alice@example.com", "password": "pa55word"}'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/tokens/authentication
```

```sh
curl -i -H "Authorization: Bearer XXXXXXXXXXXXXXXXXXXXXXXXXX" localhost:4000/v1/healthcheck
```

```sh
curl -i -H "Authorization: Bearer XXXXXXXXXXXXXXXXXXXXXXXXXX" localhost:4000/v1/movies/1
```

```sql
SELECT email FROM users WHERE activated = true;
```

```sh
BODY='{"email": "faith@example.com", "password": "pa55word"}'   // user is already activated from commands above
```

```sh
curl -d "$BODY" localhost:4000/v1/tokens/authentication
```

```sh
curl -H "Authorization: Bearer XXXXXXXXXXXXXXXXXXXXXXXXXX" localhost:4000/v1/movies/1
```

```sh
curl -X DELETE -H "Authorization: Bearer XXXXXXXXXXXXXXXXXXXXXXXXXX" localhost:4000/v1/movies/1
```

```sh
migrate create -seq -ext .sql -dir ./migrations add_permissions
```

```sh
migrate -path ./migrations -database $GREENLIGHT_DB_DSN up
```

```sh
BODY='{"name": "Grace Smith", "email": "grace@example.com", "password": "pa55word"}'
```

```sh
curl -d "$BODY" localhost:4000/v1/users
```

### Supported Go types to JSON type
- bool ⇒ JSON boolean
- string ⇒ JSON string
- int*, uint*, float*, rune ⇒ JSON number
- array, slice ⇒ JSON array
- struct, map ⇒ JSON object
- slice ofstucts ⇒ array of objects
- nil pointers, interface values, slices, maps, etc. ⇒ JSON null
- chan, func, complex* ⇒ Not supported
- time.Time ⇒ RFC3339-format JSON string
- []byte ⇒ Base64-encoded JSON string

### JSON type to Supported Go types 
- JSON boolean ⇒ bool
- JSON string ⇒ string
- JSON number ⇒ int*, uint*, float*, rune
- JSON array ⇒ array, slice
- JSON object ⇒ struct, map

### PostgreSQL and Go Types
- It’s generally sensible to align your Go and database integer types to avoid overflows or other compatibility
problems
smallint, smallserial ⇒ int16 (-32768 to 32767)
integer, serial ⇒ int32 (-2147483648 to 2147483647)
bigint, bigserial ⇒ int64 (-9223372036854775808 to 9223372036854775807)

### SMTP Server
- [MailTrap](https://mailtrap.io/)
- SMTP Settings: Credentials
