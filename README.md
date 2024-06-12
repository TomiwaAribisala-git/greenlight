### Building APIs and Web Apps with Golang

### General Commands
```sh
make help
```

```sh
make run
```

### Installing HTTP Router
```sh
go get github.com/julienschmidt/httprouter@v1
```

### Run API
```sh
go run ./cmd/api
```

```sh
go run ./cmd/api -help
```

### Handlers Requests
```sh
curl -i localhost:4000/v1/healthcheck 
```

```sh
curl -i -X OPTIONS localhost:4000/v1/healthcheck
```

### Show all Movies
```sh
curl -i -X POST localhost:4000/v1/movies
```

### Show a Movie
```sh
curl -i localhost:4000/v1/movies/123
```

### Supported Go types to JSON types
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
- Encoding of nested objects is supported. So, for example, if you have a slice of structs in Go that will encode to an array of objects in JSON.
- Any pointer values will encode as the value pointed to.
- Channels, functions and complex number types cannot be encoded. If you try to do so you’ll get a json.UnsupportedTypeError error at runtime.

### Creating a Movie
```sh
BODY='{"title":"Moana","year":2016,"runtime":107, "genres":["animation","adventure"]}'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/movies
```

### JSON type to Supported Go types 
- JSON boolean ⇒ bool
- JSON string ⇒ string
- JSON number ⇒ int*, uint*, float*, rune
- JSON array ⇒ array, slice
- JSON object ⇒ struct, map

### Installing Postgres
```sh
sudo apt install postgresql
```

```sh
psql --version
```

```sh
sudo -u postgres psql
```

```sh
psql --host=localhost --dbname=greenlight --username=greenlight
```

- [greenlight.sql](./greenlight.sql)

```sh
export GREENLIGHT_DB_DSN='postgres://greenlight:pa55word@localhost:5432/greenlight?sslmode=disable'
```

```sh
echo $GREENLIGHT_DB_DSN
```

```sql 
SELECT current_user
```

### Install Database Driver
```sh
go get github.com/lib/pq@v1
```

### SQL Migrations
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
migrate -path=./migrations -database=$EXAMPLE_DSN version
```

```sh
migrate -path=./migrations -database=$EXAMPLE_DSN goto 1
```

```sh
migrate -path=./migrations -database=$EXAMPLE_DSN up
```

```sh
migrate -path=./migrations -database=$EXAMPLE_DSN down
```

```sql
\dt
```

```sql
SELECT * FROM schema_migrations;
```

```sql
\d movies
```

### Database Insert A Movie
```sh
BODY='{"title":"Moana","year":2016,"runtime":107, "genres":["animation","adventure"]}'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/movies
```

```sh
BODY='{"title":"Black Panther","year":2018,"runtime":134,"genres":["action","adventure"]}'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/movies
```

```sh
BODY='{"title":"Deadpool","year":2016, "runtime":108,"genres":["action","comedy"]}'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/movies
```

```sh
BODY='{"title":"The Breakfast Club","year":1986, "runtime":96,"genres":["drama"]}'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/movies
```

```sql
SELECT * FROM movies;
```

### Database Fetch A Movie
```sh
curl -i localhost:4000/v1/movies/2
```

### Database PUT Update A Movie
```sh
BODY='{"title":"Black Panther","year":2018,"runtime":134,"genres":["sci-fi","action","adventure"]}'
```

```sh
curl -X PUT -d "$BODY" localhost:4000/v1/movies/2
```

### Database PATCH Update A Movie
```sh
curl -X PATCH -d '{"year": 1985}' localhost:4000/v1/movies/4
```

### Database Concurrent Update A Movie
```sh
xargs -I % -P8 curl -X PATCH -d '{"runtime":97}' "localhost:4000/v1/movies/4" < <(printf '%s\n' {1..8})
```

### Database Delete A Movie
```sh
curl -X DELETE localhost:4000/v1/movies/3
```

### Database Filter Request
- When using curl to send a request containing more than one query string parameter, you must wrap the URL in quotes for it to work correctly.
```sh
curl "localhost:4000/v1/movies?title=moana&genres=animation,adventure&page=1&page_size=5&sort=year"
```

### Database Listing Data
```sh
curl localhost:4000/v1/movies
```

### Database Filtering Lists
```sh
curl "localhost:4000/v1/movies?title=black+panther"
```

```sh
curl "localhost:4000/v1/movies?genres=adventure"
```

```sh
curl "localhost:4000/v1/movies?title=moana&genres=animation,adventure"
```

### Database Full Text Search
```sh
curl "localhost:4000/v1/movies?title=panther"
```

```sh
curl "localhost:4000/v1/movies?title=the+club"
```

### Database GIN Indexes
```sh
migrate create -seq -ext .sql -dir ./migrations add_movies_indexes'
```
```sh
migrate -path ./migrations -database $GREENLIGHT_DB_DSN up
```

### Database Sorting Data
```sh
curl "localhost:4000/v1/movies?sort=-title"
```

```sh
curl "localhost:4000/v1/movies?sort=-runtime"
```

### Database Paginating Lists
```sh
curl "localhost:4000/v1/movies?page_size=2"
```

```sh
curl "localhost:4000/v1/movies?page_size=2&page=2"
```

### Returning Pagination Metadata
```sh
curl "localhost:4000/v1/movies?page=1&page_size=2"
```

```sh
$ curl localhost:4000/v1/movies?genres=adventure
```

### Rate Limiter
```sh
go get golang.org/x/time/rate@latest
```

```sh
for i in {1..6}; do curl http://localhost:4000/v1/healthcheck; done
```

```sh
go run ./cmd/api/ -limiter-enabled=false
```

### User Model Setup And Regristration
```sh
migrate create -seq -ext=.sql -dir=./migrations create_users_table
```

```sh
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
```

### Encrypting User Passwords
```sh
go get golang.org/x/crypto/bcrypt@latest
```

### Regristering New Users
```sh
BODY='{"name": "Alice Smith", "email": "alice@example.com", "password": "pa55word"}'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/users
```

### Install Mail Package
```sh
go get github.com/go-mail/mail/v2@v2
```

### Send Users Emails
```sh
BODY='{"name": "Bob Jones", "email": "bob@example.com", "password": "pa55word"}'
```

```sh
curl -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users
```

### Send Background Emails
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

### Activate User: Set up Tokens Database Table
```sh
migrate create -seq -ext .sql -dir ./migrations create_tokens_table
```

```sh
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
```

### Register User: Send Activation Token
```sh
BODY='{"name": "Faith Smith", "email": "faith@example.com", "password": "pa55word"}'
```

```sh
curl -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users
```

### Activate User with Email Activation Token
```sh
curl -X PUT -d '{"token": "ZYGQTPU5PKKJRY7SFOAMKXPGQY"}' localhost:4000/v1/users/activated
```

### Generating Authentication Token
```sh
BODY='{"email": "alice@example.com", "password": "pa55word"'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/tokens/authentication
```

### Authenticate Authentication Token with Authorization Header
```sh
curl -d '{"email": "alice@example.com", "password": "pa55word"}' localhost:4000/v1/tokens/authentication
```

```sh
curl -i -H "Authorization: Bearer XXXXXXXXXXXXXXXXXXXXXXXXXX" localhost:4000/v1/healthcheck
```

### Authenticate Authentication Token with Authorization Header of Activated User
```sql
SELECT email FROM users WHERE activated = true;
```

```sh
BODY='{"email": "faith@example.com", "password": "pa55word"}'
```

```sh
curl -i -H "Authorization: Bearer XXXXXXXXXXXXXXXXXXXXXXXXXX" localhost:4000/v1/movies/1
```

### Permissions SQL Migrations
```sh
migrate create -seq -ext .sql -dir ./migrations add_permissions
```

```sh
migrate -path ./migrations -database $GREENLIGHT_DB_DSN up
```

### Set Read/Write Permissions for a User and give access
```sql
-- Set the activated field for alice@example.com to true.
UPDATE users SET activated = true WHERE email = 'alice@example.com';
```

```sql
-- Give all users the 'movies:read' permission
INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE code = 'movies:read') FROM users;
```

```sql
-- Give faith@example.com the 'movies:write' permission
INSERT INTO users_permissions
    VALUES (
    (SELECT id FROM users WHERE email = 'faith@example.com'),
    (SELECT id FROM permissions WHERE code = 'movies:write')
);
```

```sql
-- List all activated users and their permissions.
SELECT email, array_agg(permissions.code) as permissions
FROM permissions
INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
INNER JOIN users ON users_permissions.user_id = users.id
WHERE users.activated = true
GROUP BY email;
```

```sh
BODY='{"email": "alice@example.com", "password": "pa55word"}'
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
BODY='{"email": "faith@example.com", "password": "pa55word"}' 
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

### Grant Permissions to New User which registers an account
```sh
BODY='{"name": "Grace Smith", "email": "grace@example.com", "password": "pa55word"}'
```

```sh
curl -d "$BODY" localhost:4000/v1/users
```

```sql
SELECT email, code FROM users
INNER JOIN users_permissions ON users.id = users_permissions.user_id
INNER JOIN permissions ON users_permissions.permission_id = permissions.id
WHERE users.email = 'grace@example.com';
```

### Activating CORS
```sh
go run ./cmd/examples/cors/simple
```

### Makefile Quality Control
```sh
go install honnef.co/go/tools/cmd/staticcheck@latest
```

```sh
which staticcheck
```

### Display binary version number
```sh
./bin/api -version
```

### Basic Load Tests
```sh
go install github.com/rakyll/hey@latest
```

```sh
go run ./cmd/api/ -limiter-enabled=false
```

```sh
BODY='{"email": "alice@example.com", "password": "pa55word"}'
```

```sh
hey -d "$BODY" -m "POST" http://localhost:4000/v1/tokens/authentication
```

### Monitoring Metrics
```sh
curl http://localhost:4000/debug/vars
```

