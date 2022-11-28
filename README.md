# Image Search and Manipulation Lib

This lib contians 2 packages
- db : a package which consists of an image metadata DB and the APIs to find relavent images by title or set of tags using
       lithammer/fuzzysearch package. Internally, to improve Fuzzy match efficiency, API will perform preprossing
       1) tokenizing the string  2) eliminate the stop words and redundant words  
       3) find stemming word  4) lowercasing the words 5) sort the list of words, 
       before feeding to Fuzzy match function. 
       For example, “Andy is running home” would match “Andy runs hom”.
       
- imageop: a package which consists of the functions manipulating the image, such a resizing, croping etc (in progress)

For image metadata DB, it is currently using in-memory search by loading a JSON DB file during package init.

## Install

```
go get github.com//c-taoguan/tao-imagelib
```

## Usage

```go
package main

import (
    "fmt"
    imagedb "github.com/c-taoguan/tao-imagelib/db"
)

func main() {

   images := imagedb.SearchImagesByTitle("devrev launch")
   for _, url := range images {
       fmt.Println(url)
   }

   tags := []string{"birthday"}
   images = imagedb.SearchImagesByTags(tags)
   for _, url := range images {
       fmt.Println(url)
   }  
}
```

## License

MIT
