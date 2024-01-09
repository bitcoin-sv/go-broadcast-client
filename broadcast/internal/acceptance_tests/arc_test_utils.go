package acceptancetests

import (
	"fmt"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
	"github.com/bitcoin-sv/go-broadcast-client/httpclient"
	"github.com/rs/zerolog"
)

const errorResponse409 = `
{
	"detail": "Transaction could not be processed",
	"extraInfo": "rpc error: code = Unknown desc = key could not be found",
	"instance": null,
	"status": 409,
	"title": "Generic error",
	"txid": null,
	"type": "https://bitcoin-sv.github.io/arc/#/errors?id=_409"
}
`

func getBroadcaster(clientsNum int) (broadcast.Client, []string) {
	urls := make([]string, 0, clientsNum)
	testLogger := zerolog.Nop()
	httpClient := httpclient.NewHttpClient()
	builder := broadcast_client.Builder().WithHttpClient(httpClient)

	for i := 0; i < clientsNum; i++ {
		config := broadcast_client.ArcClientConfig{APIUrl: fakeApiUrl(i), Token: fakeApiToken(i)}
		builder.WithArc(config, &testLogger)
		urls = append(urls, config.APIUrl)
	}

	broadcaster := builder.Build()
	return broadcaster, urls
}

func fakeApiUrl(i int) string {
	return fmt.Sprintf("http://arc%d-api-url", i)
}

func fakeApiToken(i int) string {
	return fmt.Sprintf("arc%d-token", i)
}

func requestUrl(base, suffix string) string {
	return fmt.Sprintf("%s%s", base, suffix)
}
