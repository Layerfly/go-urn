# go-urn

A powerful, extensible utility for working with Uniform Resource Names (URNs) in Go.

Go port of [@jescrich/urn](https://www.npmjs.com/package/@jescrich/urn).

## URN Format

```
urn:<entity>:<id>[:<key>:<value>]*
```

Examples:
- `urn:customer:jescrich@sampledomain.com`
- `urn:order:12345:vendor:amazon:status:shipped`
- `urn:uuid:6e8bc430-9c3a-11d9-9669-0800200c9a66`

Maximum length: 255 characters.

## Installation

```bash
go get github.com/layerfly/go-urn
```

## Usage

```go
import "github.com/layerfly/go-urn"
```

### Parse

```go
u, err := urn.Parse("urn:order:12345:vendor:amazon:status:shipped")
// u.Entity    → "order"
// u.ID        → "12345"
// u.Attributes() → map[string]string{"vendor": "amazon", "status": "shipped"}
```

### Compose

```go
result, err := urn.Compose("order", "12345", map[string]string{
    "vendor": "amazon",
    "status": "shipped",
})
// → "urn:order:12345:vendor:amazon:status:shipped"
```

### Create UUID

```go
result := urn.CreateUUID("session")
// → "urn:session:550e8400-e29b-41d4-a716-446655440000"
```

### Extract Entity / ID

```go
entity, err := urn.Entity("urn:orders:1234") // → "orders"
id, err := urn.ID("urn:orders:1234")          // → "1234"
```

### Get Attribute Value

```go
val, found, err := urn.Value("urn:order:123:status:shipped", "status")
// val → "shipped", found → true
```

### Validate

```go
urn.IsValid("urn:orders:1234")        // → true
urn.IsValid("invalid:orders:1234")    // → false
```

### Add / Remove Attributes

```go
updated, err := urn.AddAttribute("urn:orders:1234", "customer", "john-doe")
// → "urn:orders:1234:customer:john-doe"

updated, err = urn.RemoveAttribute("urn:orders:1234:customer:john-doe", "customer")
// → "urn:orders:1234"
```

### Get All Attributes

```go
attrs, err := urn.GetAllAttributes("urn:orders:1234:customer:john-doe:status:pending")
// → map[string]string{"customer": "john-doe", "status": "pending"}
```

### Normalize

```go
normalized, err := urn.Normalize("URN:EXAMPLE:Animal:Ferret:Nose")
// → "urn:example:Animal:Ferret:Nose"
```

### Vendor (convenience)

```go
vendor, found, err := urn.Vendor("urn:orders:1234:vendor:amazon")
// → "amazon", true
```

## License

MIT
