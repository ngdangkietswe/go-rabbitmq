/**
 * Author : ngdangkietswe
 * Since  : 8/26/2025
 */

package models

type ExchangeConfig struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"auto_delete"`
}

type QueueConfig struct {
	Name       string `json:"name"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"auto_delete"`
}

type BindingConfig struct {
	Queue      string `json:"queue"`
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routing_key"`
}
