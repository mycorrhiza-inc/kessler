package filters

import (
	"fmt"
	"kessler/internal/dbstore"

	"go.uber.org/zap/zapcore"
)

/*
dbstore.filter is defined as :

	type Filter struct {
		ID          uuid.UUID
		Name        string
		State       string
		FilterType  string
		Description pgtype.Text
		IsActive    pgtype.Bool
		CreatedAt   pgtype.Timestamptz
		UpdatedAt   pgtype.Timestamptz
	}
*/
type Filter struct {
	dbstore.Filter
}

func (f Filter) CacheKey() string {
	return fmt.Sprintf("public:filter:%s:%s", f.Dataset, f.ID)
}

func (f Filter) String() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%v,%s,%s",
		f.ID,
		f.Name,
		f.Dataset,
		f.FilterType,
		f.Description.String,
		f.IsActive.Bool,
		f.CreatedAt.Time.String(),
		f.UpdatedAt.Time.String())
}

func (f Filter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", f.Name)
	enc.AddString("Dataset", f.Dataset)
	return nil
}
