import OrganizationTableInfiniteScroll from "@/components/Organizations/OrganizationTable";
import SearchBox from "@/components/Search/SearchBox";

const OrgLookupPage = () => {
  return (
    <>
      <h1 className="text-3xl font-bold">Organizations</h1>
      <div className="overflow-x-hidden border-r pr-4">
        <SearchBox />
        <OrganizationTableInfiniteScroll />
      </div>
    </>
  );
};

export default OrgLookupPage;
