# json:api server

This example shows usage of [Neuron framework](https://github.com/neuronlabs/neuron) by providing the http 
server with `json:api` access and multiple `postgres` databases access.
 
## Installation

Requirements:

- PostgreSQL Database - with some accounts and databases
- [Optional] [Neuron Generator](https://github.com/neuronlabs/neuron-generator) - if additional fields are provided for the models

In order to run this application set two environmental variables:

- "NEURON_DEFAULT_POSTGRES" - t postgres url for blogs and posts eg.: `postgresql://user:neuronPassword@host:port/database1`
- "NEURON_COMMENTS_POSTGRES" - t postgres access url for comments eg.: `postgresql://user:neuronPassword@host:port/database2`

Then to execute the script within this directory: 
```shell script
go run ./
```
    
### http.Server

The module [github.com/neuronlabs/neuron-extensions/server/http](https://github.com/neuronlabs/neuron-extensions/tree/master/server/http)
provides t simple `http.Server` to which you can apply any other `APIs` or some custom routes. It implements neuron/server.Server interface.

### API

The server could use direct http handlers or apply previously prepared API's. 

One of the `APIs` that could be applied to it is [github.com/neuronlabs/neuron-extensions/server/http/api/jsonapi](https://github.com/neuronlabs/neuron-extensions/tree/master/server/http/api/jsonapi). 

It is an API for the models defined in t `json:api` specification. More details about this specification could be found here - [https://jsonapi.org/](https://jsonapi.org/).

It implements all the endpoint types, as well as all the queries defined in the specification. 

What's more it supports filtering the search results by providing it's filtering system. A query parameter - `filter[field][$operator]=value` 
on the `List` endpoint would filter the results in t user defined way. 

With the possibility of sorting, paginating, including relations and filtering it makes this API very powerful.

#### Example http queries:

- Get all posts with id in range `[10:40]`, sorted by the number of `likes` paginated by `10 per page`, with inclusion of related comments and without fetching their body:

```http request
GET http://localhost:8080/v1/api/posts?fields[posts]=title,likes,blog,comments&filter[id][$gte]=10&filter[id][$lte]=40&page[limit]=10&include=comments
Accept: application/vnd.api+json
```   

- Insert comment mapped to the post with `id = 4`:
```http request
POST http://localhost:8080/v1/api/comments
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

{
    "data" : {
        "type": "comments",
        "attributes": {
            "body": "This is t comment mapped to post '4'"    
        },
        "relationships": {
            "post": {
                "data": {"type":"posts","id":"4"}
            }
        }
    }
}
```    
 
- Patch the comment we have previously created and change it's body:
    This would change only selected attribute, leaving all fields non-changed.
```http request
PATCH http://localhost:8080/v1/api/comments/1
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

{
    "data" : {
        "type": "comments",
        "id": "1",
        "attributes": {
            "body": "This is changed body of this comment"    
        }        
    }
}
```
- Delete posts relationship comments - with selected `id = 4` and `id = 5` all other comments would be there.
```http request
DELETE http://localhost:8080/v1/api/posts/4/relationships/comments
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

{
    "data" : [
        {
            "type": "comments",
            "id": "4"            
        },
        {
            "type": "comments",
            "id": "5"            
        }
    ]
}
```
- Set posts relationship comments - to the comments with `id = 3` and `id = 10`
    This would clear all current relations with comments and set them to '3' and '10' 
```http request
PATCH http://localhost:8080/v1/api/posts/4/relationships/comments
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

{
    "data" : [
        {
            "type": "comments",
            "id": "3"            
        },
        {
            "type": "comments",
            "id": "10"            
        }
    ]
}
```
- Get comments relationships for post with id `10` and receive the body compressed using gzip.
```http request
GET http://localhost:8080/v1/api/posts/10/relationships/comments
Accept: application/vnd.api+json
Accept-Encoding: gzip
```

- Delete the comment with id `1`
```http request
DELETE http://localhost:8080/v1/api/comments/1
```

### Summary

By using [Neuron framework](https://github.com/neuronlabs/neuron) t developer is responsible mostly for the business logic, whereas 
neuron takes care of querying the database for specific models and showing them to users with an encoding defined in the json:api specification.

  
