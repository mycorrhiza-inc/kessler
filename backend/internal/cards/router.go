package cards

import (
	"kessler/internal/dbstore"

	"github.com/gorilla/mux"
)

// Go ahead and implement some route handlers for these RegisterCardLookupRoutes that will go ahead and fetch data from postgres as detailed in these route handlers, and then return the data in a format that will be parsable by the frontend in this format from these zod validators:
//
//	export const BaseCardDataValidator = z.object({
//	  name: z.string(),
//	  object_uuid: z.string().uuid(),
//	  description: z.string(),
//	  timestamp: z.string(),
//	  extraInfo: z.string().optional(),
//	  index: z.number(),
//	});
//
//	export const AuthorCardDataValidator = BaseCardDataValidator.extend({
//	  type: z.literal(CardType.Author),
//	});
//
//	export const DocketCardDataValidator = BaseCardDataValidator.extend({
//	  type: z.literal(CardType.Docket),
//	});
//
// Write some endpoints that will go ahead and fetch data from the cache in the exact same mechanism that is used for the search hydration endpoints in
// /home/nicole/Documents/mycorrhizae/kessler/backend/internal/search/hydrate.go
// that could go ahead and fetch this data using
// Namely they will take in object UUID's and return the card values.
func RegisterCardLookupRoutes(r *mux.Router, db dbstore.DBTX) error {
	var nilErr error
	return nilErr
}
