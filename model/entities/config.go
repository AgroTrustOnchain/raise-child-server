package entities

import "time"

type Config struct {
	ID             string    `json:"id"`
	Key            string    `json:"key"`
	Description    string    `json:"description"`
	ActorProfileID *string   `json:"actor_profile_id"`
	ActorAddress   *string   `json:"actor_address"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type NumericConfig struct {
	Config
	Value int64 `json:"value"`
}

type StringConfig struct {
	Config
	Value string `json:"value"`
}

type PlatformConfig struct {
	Config
	Value interface{} `json:"value"`
}

func (c *PlatformConfig) GetChildAgeMaxLimitConfigKey() string {
	return "child_age_max_limit"
}

func (c *PlatformConfig) GetChildAgeMinLimitConfigKey() string {
	return "child_age_min_limit"
}

func (c *PlatformConfig) GetMinRegionStaffsConfigKey() string {
	return "min_region_staffs"
}

func (c *NumericConfig) GetTable() string {
	return "numeric_configs"
}

func (c *StringConfig) GetTable() string {
	return "string_configs"
}
