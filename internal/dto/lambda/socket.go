package lambda

type WebSocketEvent struct {
	RequestContext struct {
		RouteKey          string `json:"routeKey"`
		MessageID         string `json:"messageId"`
		EventType         string `json:"eventType"`
		ExtendedRequestID string `json:"extendedRequestId"`
		RequestTime       string `json:"requestTime"`
		MessageDirection  string `json:"messageDirection"`
		Stage             string `json:"stage"`
		ConnectedAt       int64  `json:"connectedAt"`
		RequestTimeEpoch  int64  `json:"requestTimeEpoch"`
		Identity          struct {
			UserAgent string `json:"userAgent"`
			SourceIP  string `json:"sourceIp"`
		} `json:"identity"`
		RequestID    string `json:"requestId"`
		DomainName   string `json:"domainName"`
		ConnectionID string `json:"connectionId"`
		ApiID        string `json:"apiId"`
	} `json:"requestContext"`
	Body            string `json:"body"`
	IsBase64Encoded bool   `json:"isBase64Encoded"`
}
