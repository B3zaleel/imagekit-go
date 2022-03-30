# imagekit-go
An ImageKit.io SDK for Golang.

## Installation

To install this SDK, you need to install Go and set your Go workspace first.
1. Run the command below in the terminal to install imagekit-go.
   ```powershell
   go get -u github.com/B3zaleel/imagekit-go
   ```
2. Import it in your code as shown below.
   ```go
   import imagekit "github.com/B3zaleel/imagekit-go"
   ```

## Examples

```go
package main

import (
    "fmt"
    "log"
    "encoding/json"

    imagekit "github.com/B3zaleel/imagekit-go"
)

func main() {
    imgKit := imagekit.ImageKit{
        PublicKey: "",
        PrivateKey: "",
        UrlEndpoint: "",
    }

    fileDetails, err := imgKit.Upload(
		"https://blogs.bing.com/getmedia/50f66486-7a9f-49db-8f44-44cde3ea955f/BingHomepage-KastellorizoIsland_Greece.aspx",
		"bing-image.jpg",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	} else {
		jsonBodyStr, err := json.Marshal(fileDetails)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("%s\n", string(jsonBodyStr))
		}
	}
}
```

## Related Projects

+ [ImageKit SDK for Python](https://github.com/imagekit-developer/imagekit-python)
