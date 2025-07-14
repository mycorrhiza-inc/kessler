import ErrorMessage from "@/components/style/messages/ErrorMessage";
import { fetchDocketCardData } from "../RenderedObjectCards/RednderedObjectCard";
import AllInOneServerSearch from "../SearchBar/AllInOneServerSearch";
import { TypedUrlParams } from "@/lib/types/url_params";
import { CardSize } from "@/components/style/cards/SizedCards";
import Card from "../Card/LinkedCard";
import { toFilterRecords } from "@/lib/types/typed_filters";

export default async function ConvoPage({
  convo_id,
  urlParams,
}: {
  convo_id: string;
  urlParams: TypedUrlParams;
}) {
  try {
    const card_data = await fetchDocketCardData(convo_id);
    return (
      <>
        <Card data={card_data} size={CardSize.Large} disableHref />
        <AllInOneServerSearch
          aboveSearchElement={
            <h1 className="text-2xl font-bold mb-4">
              Search {card_data.name}'s Filings
            </h1>
          }
          urlParams={urlParams}
          inherentRouteFilters={toFilterRecords({ convo_id: convo_id })}
          baseUrl={`/dockets/${convo_id}`}
        />
      </>
    );
  } catch (err) {
    return (
      <ErrorMessage
        error={err}
        message="Could not get conversation data from server :("
      />
    );
  }
}
