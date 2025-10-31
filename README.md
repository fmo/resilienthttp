## Installation 🚀

`go get github.com/fmo/resilienthttp`

## Usage

```
func main() {
  res, err := resilienthttp.Get("https://google.com")
  if err != nil {
    panic(err)
  }
  body, _ := io.ReadAll(res.Body)
  fmt.Println(string(body))
}
```

## Usage With Request Context

```
func main() {
    client := resilienthttp.NewClient()

    ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 2)
    defer cancel()

    req, err := resilienthttp.NewRequestWithContext(ctx, "GET", "http://localhost", nil)

    if err != nil {
        panic(err)
    }

    res, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    body, _ := io.ReadAll(res.Body)

    fmt.Println(string(body))
}
```
