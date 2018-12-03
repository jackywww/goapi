package main

import (
	"context"
	"gopkg.in/olivere/elastic.v5"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"encoding/json"
)

var client *elastic.Client
var host = "http://192.168.0.8:9200/"

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
	router.GET("/product-list", func(c *gin.Context) {
		termsQuery := elastic.NewTermsQuery("category.category_id", "2").Boost(1)
		boolQ := elastic.NewBoolQuery()
		boolQ = boolQ.Must(termsQuery)
		nestedQ := elastic.NewNestedQuery("category", boolQ)
		nBoolQ := elastic.NewBoolQuery()
		nBoolQ = nBoolQ.Must(nestedQ)
		constantScoreQ := elastic.NewConstantScoreQuery(nBoolQ)
		res, err = client.Search().Index("magento2_default_catalog_product").Type("product").Query(constantScoreQ).From(0).Size(20).Do(context.Background())
		b, _ := json.Marshal(res)
		c.String(200, string(b))
	})
	router.Run(":8080")
}

