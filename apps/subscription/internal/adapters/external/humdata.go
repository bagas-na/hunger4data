package external

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	pb "hunger4data/pb/subcription"
	"io"
	"log"
	"net/http"
	"net/url"
	"subscription/internal/adapters/model"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

/*
	curl -X 'GET' \
	  'https://hapi.humdata.org/api/v2/food-security-nutrition-poverty/food-security?app_identifier=&ipc_phase=3%2B&ipc_type=current&start_date=2025-01-01&admin_level=0&output_format=json&limit=20&offset=0' \
	  -H 'accept: application/json'
*/
type HAPIResponse struct {
	Data []model.Country `json:"data"`
}

// Redis sync
func GetHumData() HAPIResponse {
	authInfo := "hunger4data:hunger4data@email.com"
	appIdentifier := base64.StdEncoding.EncodeToString([]byte(authInfo))

	apiURL := "https://hapi.humdata.org/api/v2/food-security-nutrition-poverty/food-security"
	limit := 100
	offset := 0
	latestByCountry := make(map[string]model.Country)

	for {
		params := url.Values{}
		params.Add("sort_by", "reference_period_end")
		params.Add("sort_order", "desc")
		params.Add("app_identifier", appIdentifier)
		params.Add("ipc_phase", "3+")
		params.Add("ipc_type", "current")
		params.Add("admin_level", "0")
		params.Add("start_date", "2025-10-01")
		params.Add("limit", fmt.Sprintf("%d", limit))
		params.Add("offset", fmt.Sprintf("%d", offset))

		fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

		resp, err := http.Get(fullURL)
		if err != nil {
			log.Fatal(err)
			break
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var currentdata HAPIResponse
		json.Unmarshal(body, &currentdata)
		if len(currentdata.Data) == 0 {
			break
		}

		for _, record := range currentdata.Data {
			if _, exists := latestByCountry[record.Name]; !exists {
				latestByCountry[record.Name] = record
			}
		}
		offset += limit
	}

	var finalResponse HAPIResponse
	for _, country := range latestByCountry {
		finalResponse.Data = append(finalResponse.Data, country)
	}
	return finalResponse
}

// GRPC -> Rest
func GetHumDataRedis(rdb *redis.Client) (*pb.Get_Countries_Response, error) {
	ctx := context.Background()
	cacheKey := "countries:latest"
	cachedBytes, err := rdb.Get(ctx, cacheKey).Bytes()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("cache empty")
		}
		return nil, err
	}

	response := &pb.Get_Countries_Response{
		Message: "Success",
	}

	if err := proto.Unmarshal(cachedBytes, response); err != nil {
		return nil, err
	}

	return response, nil
}
