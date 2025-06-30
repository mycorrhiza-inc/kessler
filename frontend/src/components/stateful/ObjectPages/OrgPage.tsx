import ErrorMessage from "@/components/style/messages/ErrorMessage"
import { fetchAuthorCardData } from "../RenderedObjectCards/RednderedObjectCard"
import AllInOneServerSearch from "../SearchBar/AllInOneServerSearch"
import { TypedUrlParams } from "@/lib/types/url_params"
import { CardSize } from "@/components/style/cards/SizedCards"
import Card from "../Card/LinkedCard"

export default async function OrgPage({ org_id, urlParams }: { org_id: string, urlParams: TypedUrlParams }) {
  try {
    const card_data = await fetchAuthorCardData(org_id)
    return (
      <>
        <Card data={card_data} size={CardSize.Large} disableHref />
        <AllInOneServerSearch
          aboveSearchElement={<h1 className="text-2xl font-bold mb-4">Search {card_data.name}'s Filings</h1>
          }
          urlParams={urlParams}

          inherentRouteFilters={{ "conversation_id": org_id }}
          baseUrl={`/orgs/${org_id}`}
        />
      </>
    )
  } catch (err) {
    return <ErrorMessage error={err} message="Could not get organization data from server :(" />
  }
}
