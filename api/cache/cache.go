package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/arkag/cope-and-hope-weather/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type WeatherFetcher interface {
	FetchWeather(city string) (models.WeatherData, error)
}

type CachedClient struct {
	RealFetcher WeatherFetcher
	TableName   string
	Dynamo      *dynamodb.Client
}

// CachedItem maps exactly to how the row looks in DynamoDB
type CachedItem struct {
	City      string             `dynamodbav:"city"`
	Weather   models.WeatherData `dynamodbav:"weather"`
	ExpiresAt int64              `dynamodbav:"expires_at"` // We'll configure AWS to auto-delete when this time passes
}

func (c *CachedClient) FetchWeather(city string) (models.WeatherData, error) {
	ctx := context.Background()

	// Circuit breaker: If we haven't configured Dynamo (e.g. local dev), bypass cache
	if c.Dynamo == nil {
		slog.Debug("cache bypassed (no dynamo client)", "city", city)
		return c.RealFetcher.FetchWeather(city)
	}

	// 1. Check DynamoDB for the city
	out, err := c.Dynamo.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(c.TableName),
		Key: map[string]types.AttributeValue{
			"city": &types.AttributeValueMemberS{Value: city},
		},
	})

	if err == nil && out.Item != nil {
		var item CachedItem
		if err := attributevalue.UnmarshalMap(out.Item, &item); err == nil {
			// Ensure it hasn't expired
			if time.Now().Unix() < item.ExpiresAt {
				slog.Debug("cache hit", "city", city)
				return item.Weather, nil
			}
		}
	}

	// 2. Cache miss. Fetch from API
	slog.Debug("cache miss, fetching from API", "city", city)
	weather, err := c.RealFetcher.FetchWeather(city)
	if err != nil {
		return models.WeatherData{}, err
	}

	// 3. Save back to DynamoDB (Expires in 2 hours)
	item := CachedItem{
		City:      city,
		Weather:   weather,
		ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
	}
	av, _ := attributevalue.MarshalMap(item)

	// We ignore errors on PutItem; if caching fails, the user still gets their weather
	_, _ = c.Dynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.TableName),
		Item:      av,
	})

	return weather, nil
}
