# products-parser-grpc
from course

### Running:
```make
make run # runs docker container
``` 
```make
make grpc-client # runs grpc client to test request to server
``` 

### Methods:
```go
// parses csv file and adds data to database
// returns result(how much added etc)

// args in FetchRequest:
// url(string) - string to server where csv is located
rpc Fetch(FetchRequest) returns (FetchResponse) {}; 
```
```go
// returns list of all products

// supported parameters in ListRequest:
// reversed(bool) - should returned list be reversed
// mode(SortingMode(string)) - type of sorting the returned list
// pageSize(int64) - how much products should be returned
// pageOffset(int64) - how much products has to be skipped
rpc List(ListRequest) returns (ListResponse) {};
```