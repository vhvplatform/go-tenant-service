package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vhvplatform/go-tenant-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UsageMetricsRepository handles usage metrics data access
type UsageMetricsRepository struct {
	collection *mongo.Collection
}

// NewUsageMetricsRepository creates a new usage metrics repository
func NewUsageMetricsRepository(db *mongo.Database) *UsageMetricsRepository {
	collection := db.Collection("usage_metrics")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tenant_id", Value: 1},
				{Key: "period", Value: 1},
			},
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetExpireAfterSeconds(2592000), // 30 days
		},
	}

	_, _ = collection.Indexes().CreateMany(ctx, indexes)

	return &UsageMetricsRepository{collection: collection}
}

// RecordMetrics records usage metrics for a tenant
func (r *UsageMetricsRepository) RecordMetrics(ctx context.Context, metrics *domain.UsageMetrics) error {
	metrics.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, metrics)
	if err != nil {
		return fmt.Errorf("failed to record metrics: %w", err)
	}

	return nil
}

// GetMetricsByTenant retrieves usage metrics for a tenant
func (r *UsageMetricsRepository) GetMetricsByTenant(ctx context.Context, tenantID string, period string) (*domain.UsageMetrics, error) {
	var metrics domain.UsageMetrics
	filter := bson.M{"tenant_id": tenantID}

	if period != "" {
		filter["period"] = period
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	err := r.collection.FindOne(ctx, filter, opts).Decode(&metrics)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return &metrics, nil
}

// GetMetricsHistory retrieves historical usage metrics for a tenant
func (r *UsageMetricsRepository) GetMetricsHistory(ctx context.Context, tenantID string, limit int) ([]*domain.UsageMetrics, error) {
	if limit <= 0 {
		limit = 30
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"tenant_id": tenantID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics history: %w", err)
	}
	defer cursor.Close(ctx)

	var metrics []*domain.UsageMetrics
	if err := cursor.All(ctx, &metrics); err != nil {
		return nil, fmt.Errorf("failed to decode metrics: %w", err)
	}

	return metrics, nil
}

// UpdateMetrics updates or creates usage metrics for a tenant in a specific period
func (r *UsageMetricsRepository) UpdateMetrics(ctx context.Context, tenantID string, period string, apiCalls, storage, bandwidth int64) error {
	filter := bson.M{
		"tenant_id": tenantID,
		"period":    period,
	}

	update := bson.M{
		"$inc": bson.M{
			"api_call_count": apiCalls,
			"storage_used":   storage,
			"bandwidth_used": bandwidth,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to update metrics: %w", err)
	}

	return nil
}

// GetAggregatedMetrics retrieves aggregated metrics across all tenants
func (r *UsageMetricsRepository) GetAggregatedMetrics(ctx context.Context, period string) (map[string]interface{}, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"period": period}}},
		{{Key: "$group", Value: bson.M{
			"_id":             nil,
			"total_api_calls": bson.M{"$sum": "$api_call_count"},
			"total_storage":   bson.M{"$sum": "$storage_used"},
			"total_bandwidth": bson.M{"$sum": "$bandwidth_used"},
			"tenant_count":    bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate metrics: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregated metrics: %w", err)
	}

	if len(results) == 0 {
		return map[string]interface{}{
			"total_api_calls": 0,
			"total_storage":   0,
			"total_bandwidth": 0,
			"tenant_count":    0,
		}, nil
	}

	// Convert bson.M to map[string]interface{}
	result := make(map[string]interface{})
	for k, v := range results[0] {
		if k != "_id" {
			result[k] = v
		}
	}

	return result, nil
}
