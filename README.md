# Usage:
```go
package main

import (
	"fmt"
	"jarnpher/godo"
	"log"
	"strconv"
)

func main() {
	g := godo.New()

	g.Use(func(h godo.HandleFunc) godo.HandleFunc {
		return func(c *godo.Context) {
			fmt.Println("lalalala")

			h(c)
		}
	})

	g.Use(func(h godo.HandleFunc) godo.HandleFunc {
		return func(c *godo.Context) {
			fmt.Println("nananana")

			h(c)
		}
	})

	g.Get("/users/:id", func(c *godo.Context) {
		id, _ := strconv.Atoi(c.Src())
		c.JSON(300, struct {
			Id   int
			Name string
			Age  int
		}{
			Id:   id,
			Name: "aaa",
			Age:  10,
		})
	})

	g.Get("/users/do", func(c *godo.Context) {
		c.JSON(300, struct {
			Name string
			Age  int
		}{
			Name: "aaa",
			Age:  10,
		})
	})

	err := g.Run(":9999")
	if err != nil {
		log.Fatalln("Server Internal Error")
	}
}

```