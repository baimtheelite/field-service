package clients

import (
	"context"
	"field-service/clients/config"
	"field-service/common/util"
	"field-service/config"
	config2 "field-service/config"
	"field-service/constants"
	"fmt"
	"net/http"
	"time"
)

type UserClient struct {
	client config.IClientConfig
}

type IUserClient interface {
	GetUserByToken(ctx context.Context, token string) (*UserData, error)
}

func NewUserClient(client config.IClientConfig) IUserClient {
	return &UserClient{
		client: client,
	}
}

func (u *UserClient) GetUserByToken(ctx context.Context) (*UserData, error) {
	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		config2.Config.AppName, u.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var response UserResponse
	request := u.client.Client().Clone().
		Set(constants.Authorization, bearerToken).
		Set(constants.XApiKey, apiKey).
		Set(constants.XServiceName, config.Config.AppName).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL()))

	resp, _, errs := request.EndStruct(&response)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user response: %s", response.Message)
	}

	return &response.Data, nil
}
