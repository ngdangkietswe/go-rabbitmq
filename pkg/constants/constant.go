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
	QueueNotification Queue = "queue_notification"
	QueueLog          Queue = "queue_log"
)

type Exchange string

const (
	ExchangeNotification Exchange = "exchange_notification"
	ExchangeLog          Exchange = "exchange_log"
)

type RoutingKey string

const (
	RoutingKeyNotification RoutingKey = "notification.created"
	RoutingKeyLog          RoutingKey = "log.created"
)
