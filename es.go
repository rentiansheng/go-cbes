package cbes

import (
    "gopkg.in/olivere/elastic.v2"
//    "fmt"
    "fmt"
)

// connect to elastic search and build the client
func connectEs (settings *Settings) (*elastic.Client, error) {
    client, err := elastic.NewClient(elastic.SetURL(settings.ElasticSearch.Urls...))
    if err != nil {
        return nil , err
    }

    return client, nil
}

// Check if the index exists
func checkIndex(settings *Settings, client *elastic.Client) (bool, error) {
    exists, err := client.IndexExists(settings.ElasticSearch.Index).Do()
    if err != nil {
        return false, err
    }

    if !exists {
        return false, nil
    }

    return true, nil
}

// Create Index
func createIndex(settings *Settings, client *elastic.Client) (bool, error) {
    builder, err := client.CreateIndex(settings.ElasticSearch.Index).Do()
    if err != nil {
        return false, err
    }

    return builder.Acknowledged, nil
}

// Open connection
func OpenEs (settings *Settings) (*elastic.Client, error) {
    client, err := connectEs(settings)
    if err != nil {
        return nil, err
    }

    indexExists, err := checkIndex(settings, client)
    if err != nil {
        return nil, err
    }

    if !indexExists {
        createIndex(settings, client)
    }

    return client, nil
}

// put model mapping
func addMapping(mapping string, modelName string) error {
    index := dbSettings.ElasticSearch.Index
    es := *Connection.es

    res, err := es.PutMapping().Index(index).Type(modelName).BodyString(mapping).Do()
    if err != nil {
        return err
    }
    if res == nil {
        return fmt.Errorf("expected put mapping response; got: %v", res)
    }
    if !res.Acknowledged {
        return fmt.Errorf("expected put mapping ack; got: %v", res.Acknowledged)
    }

    return nil
}