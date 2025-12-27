package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tenant represents an organization/tenant in the system
type Tenant struct {
	ID               primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name             string                 `bson:"name" json:"name"`
	Domain           string                 `bson:"domain,omitempty" json:"domain,omitempty"`
	SubscriptionTier string                 `bson:"subscription_tier" json:"subscription_tier"`
	IsActive         bool                   `bson:"is_active" json:"is_active"`
	Settings         map[string]interface{} `bson:"settings,omitempty" json:"settings,omitempty"`
	IsolationConfig  *IsolationConfig       `bson:"isolation_config,omitempty" json:"isolation_config,omitempty"`
	CreatedAt        time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time              `bson:"updated_at" json:"updated_at"`
}

// IsolationConfig defines multi-tenant isolation settings
type IsolationConfig struct {
	DatabaseIsolation bool   `bson:"database_isolation" json:"database_isolation"`
	Namespace         string `bson:"namespace" json:"namespace"`
	DataEncryption    bool   `bson:"data_encryption" json:"data_encryption"`
}

// TenantUser represents a user's relationship with a tenant
type TenantUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID  string             `bson:"tenant_id" json:"tenant_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Role      string             `bson:"role" json:"role"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Subscription tiers
const (
	SubscriptionFree         = "free"
	SubscriptionBasic        = "basic"
	SubscriptionProfessional = "professional"
	SubscriptionEnterprise   = "enterprise"
)

// UsageMetrics tracks tenant resource usage
type UsageMetrics struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID      string             `bson:"tenant_id" json:"tenant_id"`
	APICallCount  int64              `bson:"api_call_count" json:"api_call_count"`
	StorageUsed   int64              `bson:"storage_used" json:"storage_used"`
	BandwidthUsed int64              `bson:"bandwidth_used" json:"bandwidth_used"`
	Period        string             `bson:"period" json:"period"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

// TenantConfiguration stores custom tenant configurations
type TenantConfiguration struct {
	Key       string      `bson:"key" json:"key"`
	Value     interface{} `bson:"value" json:"value"`
	Type      string      `bson:"type" json:"type"`
	Locked    bool        `bson:"locked" json:"locked"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
}
