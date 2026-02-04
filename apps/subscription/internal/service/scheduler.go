package service

import (
	"context"
	pb "hunger4data/pb/subcription"
	"log"
	"subscription/internal/adapters/external"

	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"google.golang.org/protobuf/proto"
)

func fetchAndSync(rdb *redis.Client) error {
	rawHapiData := external.GetHumData()

	protoResponse := &pb.Get_Countries_Response{
		Message: "Success",
	}

	for _, c := range rawHapiData.Data {
		pCountry := &pb.Country{
			Id:                        c.Id,
			Name:                      c.Name,
			IpcPhase:                  c.IpcPhase,
			LocationCode:              c.LocationCode,
			PopulationInPhase:         c.PopulationInPhase,
			PopulationFractionInPhase: c.PopulationFractionInPhase,
		}
		protoResponse.Countries = append(protoResponse.Countries, pCountry)
	}

	binaryData, _ := proto.Marshal(protoResponse)
	return rdb.Set(context.Background(), "countries:latest", binaryData, 24*time.Hour).Err()
}

func StartSyncWithoutScheduler(rdb *redis.Client) {
	err := fetchAndSync(rdb)
	if err != nil {
		log.Printf("Cron Error: Failed to sync data: %v", err)
	} else {
		log.Println("Cron: Sync successful!")
	}
}

func StartSyncScheduler(rdb *redis.Client) {
	c := cron.New()
	_, err := c.AddFunc("0 */12 * * *", func() {
		log.Println("Cron: Starting HAPI Data Sync...")
		err := fetchAndSync(rdb)
		if err != nil {
			log.Printf("Cron Error: Failed to sync data: %v", err)
		} else {
			log.Println("Cron: Sync successful!")
		}
	})

	if err != nil {
		log.Fatal("Could not start cron:", err)
	}

	c.Start()
}
