package api

import (
	"errors"
	"fmt"
	"reflect"

	elastic "gopkg.in/olivere/elastic.v3"

	"github.com/cafebazaar/booker-resources/common"
)

const (
	esIndex = "resources"
)

var (
	_elasticClient *elastic.Client
	debugMode      = common.ConfigString("LOG_LEVEL") == "DEBUG"
)

func elasticClient() (*elastic.Client, error) {
	if _elasticClient == nil {
		elasticURL := common.ConfigString("ELASTIC_URL")
		if elasticURL == "" {
			return nil, errors.New("No ELASTIC_URL was given")
		}

		var err error
		var client *elastic.Client
		if debugMode {
			client, err = elastic.NewClient(
				elastic.SetSniff(false),
				elastic.SetURL(elasticURL),
				elastic.SetTraceLog(common.LogrusInfoLogger),
				elastic.SetInfoLog(common.LogrusInfoLogger),
				elastic.SetErrorLog(common.LogrusErrorLogger),
			)
		} else {
			client, err = elastic.NewClient(
				elastic.SetSniff(false),
				elastic.SetURL(elasticURL),
				elastic.SetErrorLog(common.LogrusErrorLogger),
			)
		}
		if err != nil {
			return nil, fmt.Errorf("Error while elastic.NewClient: %s", err)
		}

		_elasticClient = client
	}

	return _elasticClient, nil
}

type item struct {
	Category string
	Name     string
	Spec     map[string]string
}

//func getReservation(objectURI string, timestamp uint64) (*reservation, error) {
//	client, err := elasticClient()
//	if err != nil {
//		return nil, fmt.Errorf("Error while getting elastic client: %s", err)
//	}
//
//	exists, err := client.IndexExists(esIndex).Do()
//	if err != nil {
//		return nil, err
//	}
//	if !exists {
//		return nil, nil
//	}
//
//	query := elastic.NewBoolQuery().Must(
//		elastic.NewTermQuery("URI", objectURI),
//		elastic.NewRangeQuery("startTimestamp").Lt(timestamp),
//		elastic.NewRangeQuery("endTimestamp").Gt(timestamp),
//	)
//	searchResult, err := client.Search().
//		Index(esIndex).
//		Query(query).
//		From(0).Size(1).
//		Do()
//	if err != nil {
//		return nil, fmt.Errorf("Error while querying elastic for reservation at the given time: %s", err)
//	}
//
//	var rsv reservation
//	for _, item := range searchResult.Each(reflect.TypeOf(rsv)) {
//		rsv, ok := item.(reservation)
//		if !ok {
//			return nil, fmt.Errorf("Failed to convert item to reservation. item=%v", item)
//		}
//
//		return &rsv, nil
//	}
//
//	return nil, nil
//}

func getItem(categoryName, itemName string) (*item, error) {
	client, err := elasticClient()
	if err != nil {
		return nil, fmt.Errorf("Error while getting elastic client: %s", err)
	}

	exists, err := client.IndexExists(esIndex).Do()
	if err != nil {
		return nil, fmt.Errorf("Error while checking elastic index: %s", err)
	}
	if !exists {
		return nil, nil
	}

	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("Category", categoryName),
		elastic.NewTermQuery("Name", itemName),
	)

	result, err := client.Search().
		Index(esIndex).
		Query(query).
		From(0).Size(1).
		Do()

	if err != nil {
		return nil, fmt.Errorf("Error in getting resource: %s", err)
	}

	var it item
	for _, res := range result.Each(reflect.TypeOf(it)) {
		it, ok := res.(item)
		if !ok {
			return nil, fmt.Errorf("failed to convert result to resource object: %v", res)
		}
		return &it, nil
	}
	return nil, nil
}

func createItem(categoryName, itemName string, specMap map[string]string) (*item, error) {
	client, err := elasticClient()
	if err != nil {
		return nil, fmt.Errorf("Error while getting elastic client: %s", err)
	}

	exists, err := client.IndexExists(esIndex).Do()
	if err != nil {
		return nil, fmt.Errorf("Error while checking elastic index: %s", err)
	}
	if !exists {
		createIndex, err := client.CreateIndex(esIndex).Do()
		if err != nil {
			return nil, fmt.Errorf("Error while creating elastic index: %s", err)
		}
		if !createIndex.Acknowledged {
			return nil, errors.New("Error while creating elastic index: Not Acknowledged")
		}
	}

	it := item{
		Category: categoryName,
		Name:     itemName,
		Spec:     specMap,
	}

	_, err = client.Index().
		Index(esIndex).
		Type("item").
		BodyJson(&it).
		Do()
	if err != nil {
		return nil, fmt.Errorf("Error while creating the reservation: %s", err)
	}

	return &it, nil
}

func getCategoryItems(categoryName string) ([]*item, error) {
	client, err := elasticClient()
	if err != nil {
		return nil, fmt.Errorf("Error while getting elastic client: %s", err)
	}

	exists, err := client.IndexExists(esIndex).Do()
	if err != nil {
		return nil, fmt.Errorf("Error while checking elastic index: %s", err)
	}
	if !exists {
		return nil, nil
	}

	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("Category", categoryName),
	)

	results, err := client.Search().
		Index(esIndex).
		Query(query).
		Do()

	if err != nil {
		return nil, fmt.Errorf("Error in getting resource: %s", err)
	}

	ret := make([]*item, 0)
	for _, res := range results.Each(reflect.TypeOf(item{})) {
		it, ok := res.(item)
		if ok {
			ret = append(ret, &it)
		}
	}
	return ret, nil
}

func getCategories() ([]string, error) {
	client, err := elasticClient()
	if err != nil {
		return nil, fmt.Errorf("Error while getting elastic client: %s", err)
	}

	exists, err := client.IndexExists(esIndex).Do()
	if err != nil {
		return nil, fmt.Errorf("Error while checking elastic index: %s", err)
	}
	if !exists {
		return nil, nil
	}

	agg := elastic.NewTermsAggregation().Field("Category")

	results, err := client.Search().
		Index(esIndex).
		Query(elastic.NewMatchAllQuery()).
		Aggregation("categories", agg).
		Do()

	if err != nil {
		return nil, fmt.Errorf("Error in getting resource: %s", err)
	}

	aggs, found := results.Aggregations.Terms("categories")

	if !found {
		return nil, fmt.Errorf("error in terms aggregation: %s", err)
	}

	ret := make([]string, 0)
	for _, categoryBucket := range aggs.Buckets {
		categoryName, ok := categoryBucket.Key.(string)
		if ok {
			ret = append(ret, categoryName)
		}
	}

	return ret, nil
}
