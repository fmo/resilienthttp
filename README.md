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
