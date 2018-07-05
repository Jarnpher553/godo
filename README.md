# Usage:
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

        g.Get("/", func(c *godo.Context) {
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
