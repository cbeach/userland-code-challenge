package main

import (
	"errors"
	"fmt"
  "bytes"
  "strings"
	"io/ioutil"
	"net/http"
  "encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)


type Person struct {
  Name string
  FavoriteColor string
  FavoriteAnimal string
}

type People struct {
  People []Person `json:"data"`
}

var (
	// DefaultHTTPPostAddress Default Address
	DefaultHTTPPostAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func getNames(body map[string][]map[string]string)string {
  for key, value := range body {
    fmt.Print("%v: %v", key, value);
  }
  return "all done"
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  var people People
  // Parse the json body
  json.Unmarshal([]byte(request.Body), &people)
  var nameBuffer bytes.Buffer
  // iterate through the request data and pull out all of the names
  for _, person := range people.People {
    nameBuffer.WriteString(fmt.Sprintf("%v\n", person.Name))
  }

  // Start: a bunch of error handling that I gor for free from `sam init`
  fmt.Println(nameBuffer.String())
	resp, err := http.Post(DefaultHTTPPostAddress, "application/json", strings.NewReader(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}
  // End: a bunch of error handling that I gor for free from `sam init`

	return events.APIGatewayProxyResponse{
		Body:       nameBuffer.String(),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
