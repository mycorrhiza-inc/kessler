package admin

import (
	"context"
	"kessler/gen/dbstore"
	"kessler/quickwit"
)

func IndexAllConversations(q dbstore.Queries, ctx context.Context) error {
	conversations, err := q.DocketConversationList(ctx)
	if err != nil {
		return err
	}
	quickwit.IndexConversations(conversations)
	return nil
}
