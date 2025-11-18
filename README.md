# Domain Trust Go SDK

Official Go client for the **[Domain Trust API](https://domain-trust.globalcyberalliance.org)** — an initiative of
the [Global Cyber Alliance (GCA)](https://globalcyberalliance.org) dedicated to building trust in the domain name
ecosystem and reducing domain abuse worldwide.

**Domain Trust**, part of GCA’s
[Internet Integrity Program](https://globalcyberalliance.org/our-work/promoting-internet-integrity/), connects
registries, registrars, ISPs, CERTs, law enforcement, cyber responders, researchers, and financial institutions to share
data on malicious domains and take coordinated action.

Together, participants:

* Share data on malicious domains
* Identify and mitigate abuse faster
* Build objective, collaborative trust in the DNS ecosystem

Over **36 million+ domains** have been contributed by community members to date. Participation is free and open to
qualified organizations. [Learn more](https://globalcyberalliance.org/domain-trust).

---

## Installation

```bash
go get github.com/globalcyberalliance/domain-trust-go
```

---

## Quick Start

```go
package main

import (
	"context"
	"log"
	"time"

	dt "github.com/globalcyberalliance/domain-trust-go/v2"
	dtm "github.com/globalcyberalliance/domain-trust-go/v2/model"
)

func main() {
	ctx := context.Background()

	// Initialize the API client with your API key (optionally, call c.Login(ctx, "YOUR_EMAIL", "YOUR_PASSWORD")).
	c := dt.New("YOUR_API_KEY", dt.WithDebug(false), dt.WithTimeout(30*time.Second))

	// Find all domains submitted in the last 24 hours.
	filter := &dtm.DomainFilter{
		MetadataFilter: dtm.MetadataFilter{Limit: dtm.DefaultMetadataLimit},
		CreatedAfter:   time.Now().Add(-24 * time.Hour),
	}

	// Paginate over results. Alternatively, call c.FindDomains(ctx, filter) to only retrieve one page.
	domainIterator, err := c.FindDomainsPaged(ctx, filter)
	if err != nil {
		log.Fatalf("find domains: %v", err)
	}

	domains := make([]*dtm.Domain, 0, dtm.DefaultMetadataLimit)

	for domainIterator.Next() {
		domains = append(domains, domainIterator.Value())
	}

	if domainIterator.Err() != nil {
		log.Fatalf("find domains: %v", err)
	}
}
```

---

## Authentication

You should log in to the [Domain Trust UI](https://domain-trust.globalcyberalliance.org) and create an API key there.

Alternatively, you can log in with your registered Domain Trust credentials to obtain an API key:

```go
ctx := context.Background()
c := dt.New("") // Initialize without an API key for login.

key, err := c.Login(ctx, "YOUR_EMAIL", "YOUR_PASSWORD")
if err != nil {
    log.Fatalf("login: %v", err)
}

fmt.Printf("Your API key: %s\n", key.Value)
```

Then reuse that key when creating your main API client:

```go
c = dt.New(key.Value)
```

---

## Working with Domains

### Submit (create) new domains

```go
submissions := []*dtm.DomainSubmission{
    {
        Domain:         "malicious-example.com",
        AbuseType:      dtm.DomainAbuseTypePhishing,
        Activity:       dtm.DomainActivityActive,
        Classification: dtm.DomainClassificationPossiblyMalicious,
    },
}

errors, err := c.CreateDomains(ctx, submissions...)
if err != nil {
    log.Fatalf("create domains: %v", err)
}

if len(errors) > 0 {
    for _, e := range errors {
        fmt.Printf("Error for %s: %s\n", e.Domain, e.Message)
    }
}
```

---

### Query domains (basic)

```go
filter := &dtm.DomainFilter{
    AbuseType: "phishing",
    Limit:     dtm.DefaultMetadataLimit,
}

domains, err := c.FindDomains(ctx, filter)
if err != nil {
    log.Fatalf("find domains: %v", err)
}

for _, d := range domains {
    fmt.Println(d.Name)
}
```

---

### Query domains with pagination iterator

If your result set is large, you can stream through pages using the iterator interface:

```go
iter, err := c.FindDomainsPaged(ctx, &dtm.DomainFilter{Limit: 100})
if err != nil {
    log.Fatal(err)
}

for iter.Next() {
    d := iter.Value()
    fmt.Printf("Domain: %s (%s)\n", d.Name, d.AbuseType)
}

if err := iter.Err(); err != nil {
    log.Fatalf("pagination failed: %v", err)
}
```

This uses an efficient **lazy pagination** mechanism — only one page is kept in memory at a time, and new pages are
fetched automatically as you iterate.

---

## Configuration Options

You can customize the client using functional options:

```go
c := dt.New("YOUR_API_KEY",
    dt.WithContentType(dt.ContentTypeJSON),
    dt.WithDebug(true),
    dt.WithTimeout(10*time.Second),
)
```

| Option             | Description                                       |
|--------------------|---------------------------------------------------|
| `WithClient`       | Use a custom `*http.Client`                       |
| `WithContentType`  | Override default content type (`CBOR` by default) |
| `WithDebug`        | Enables verbose request/response logging          |
| `WithEncodingType` | Override encoding (`ZSTD` by default)             |
| `WithTimeout`      | Sets HTTP client timeout                          |

---

## API Documentation

Comprehensive endpoint documentation is available here: *
*[https://domain-trust.docs.globalcyberalliance.org](https://domain-trust.docs.globalcyberalliance.org)**

---

## Models

All object models are defined in a separate model package, generated automatically from Domain Trust’s internal schemas.
These files are part of our code generation workflow and should not be edited manually.

---

## 🪪 License

Apache 2.0 License.