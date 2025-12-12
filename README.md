# runninghub_tools
RunningHub go工具库

# Usage

## get
```bash
go get github.com/Friday-fighting/runninghub_tools@latest
```

## use
```go
package main

import (
    "github.com/Friday-fighting/runninghub_tools/runninghub_client"
	"context"
	"fmt"
)

func main() {
    ctx := context.Background()
    const apiKey = "your api key"
    client := runninghub_client.NewClient(&runninghub_client.RunningHubClientConfig{
        ApiKey: apiKey,
    })
    info, err := client.GetAccountInfo(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Println(info)
}

```


