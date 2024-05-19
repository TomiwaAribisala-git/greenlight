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
curl localhost:4000/v1/healthcheck
```

```sh
curl -X POST localhost:4000/v1/movies
```

```sh
curl localhost:4000/v1/movies/123
```

```sh
go install github.com/rakyll/hey@latest
```

```sh
BODY='{"title":"moana","year":2000,"runtime":123,"genres":["sci-fi","sci-f"]}'
```

```sh
curl -i -d "$BODY" localhost:4000/v1/movies
```

### Supported Go types to JSON type
- bool ⇒ JSON boolean
- string ⇒ JSON string
- int*, uint*, float*, rune ⇒ JSON number
- array, slice ⇒ JSON array
- struct, map ⇒ JSON object
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