# groonga-manage

# usage

```main.go
package main

import(
    "log"

    "github.com/rssh-jp/gogroongahttpdmanage"
)

func main(){
    gogroongahttpdmanage.Initialize("http", "local-groonga.com", "10041")
    res, err := gogroongahttpdmanage.Select("table=Users")
    if err != nil{
        log.Fatal(err)
    }

    log.Println(res)
}
```


