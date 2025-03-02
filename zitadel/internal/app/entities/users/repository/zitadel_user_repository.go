package repository

import (
	"context"

	"app/internal/app/entities/users/domain"

	"github.com/muonsoft/errors"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	object "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object/v2"
	user "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"
)

type ZitadelUserRepository struct {
	client *client.Client
}

func NewZitadelUserRepository(client *client.Client) *ZitadelUserRepository {
	return &ZitadelUserRepository{client: client}
}

func (r *ZitadelUserRepository) FindByIDs(ctx context.Context, ids []string) (map[string]*domain.User, error) {
	response, err := r.client.UserServiceV2().ListUsers(ctx, &user.ListUsersRequest{
		Query: &object.ListQuery{Limit: uint32(len(ids))},
		Queries: []*user.SearchQuery{
			{
				Query: &user.SearchQuery_InUserIdsQuery{
					InUserIdsQuery: &user.InUserIDQuery{UserIds: ids},
				},
			},
		},
	})
	if err != nil {
		return nil, errors.Errorf("find users: %w", err)
	}

	users := make(map[string]*domain.User, len(response.Result))
	for _, in := range response.Result {
		human := in.GetHuman()
		if human == nil {
			continue
		}
		out := &domain.User{
			ID:        in.UserId,
			Email:     human.Email.Email,
			FirstName: human.Profile.GivenName,
			LastName:  human.Profile.FamilyName,
		}

		users[out.ID] = out
	}

	return users, nil
}
