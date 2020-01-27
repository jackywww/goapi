
package main

import (
        "context"
        "gopkg.in/olivere/elastic.v5"
        "github.com/gin-gonic/gin"
        "log"
        "os"
        "encoding/json"
        "strconv"
        "fmt"
        "github.com/tidwall/gjson"
        "reflect"
)

var client *elastic.Client
var host = "http://192.168.1.8:9200/"

func main() {
        router := gin.Default()
        errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
        var err error
        client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(host))

        if err != nil {
                panic(err)
        }

        var res *elastic.SearchResult
        // Query string parameters are parsed using the existing underlying request object.
        // The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
        router.GET("/products/:categoryId/:page/:size", func(c *gin.Context) {
                categoryId := c.Param("categoryId")
                page, _ := strconv.Atoi(c.Param("page"))
                size, _ := strconv.Atoi(c.Param("size"))

                termsQuery := elastic.NewTermsQuery("category.category_id", categoryId).Boost(1)
                boolQ := elastic.NewBoolQuery()
                boolQ = boolQ.Must(termsQuery)
                nestedQ := elastic.NewNestedQuery("category", boolQ)
                nBoolQ := elastic.NewBoolQuery()
                nBoolQ = nBoolQ.Must(nestedQ)
                constantScoreQ := elastic.NewConstantScoreQuery(nBoolQ)

                //gFields := [...]string{"entity_id","price","name","image","stock"}
                filter := elastic.NewFetchSourceContext(true).Include("entity_id","price","name","image","stock")
                res, err = client.Search().Index("magento2_default_catalog_product").Type("product").Query(constantScoreQ).FetchSourceContext(filter).From(page).Size(size).Do(context.Background())
                b, _ := json.Marshal(res)
                //fmt.Println(string(b))
                value := gjson.Get(string(b), "hits.hits.#._source")
                fmt.Println(reflect.TypeOf(client.Search()).String())
                c.String(200, value.String())
        })
        router.Run(":8080")
}
