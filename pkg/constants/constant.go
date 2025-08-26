/**
 * Author : ngdangkietswe
 * Since  : 8/15/2025
 */

package constants

type Env string

const (
	EnvLocal       Env = "local"
	EnvDevelopment Env = "development"
	EnvProduction  Env = "production"
)

type Queue string

const (
	QueueNotification Queue = "notification_queue"
)

type Exchange string

const (
	ExchangeNotification Exchange = "notification_exchange"
)

type RoutingKey string

const (
	RoutingKeyNotification RoutingKey = "notification"
)
