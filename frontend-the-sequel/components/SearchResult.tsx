type SearchFields = {
  name: string;
  text: string;
  docketID: string;
};

type SearchResultProps = {
  data: SearchFields;
};

const SearchResult = ({ data }: SearchResultProps) => {
  return (
    <div
      style={{
        padding: "15px",
        border: "1px solid white",
        borderRadius: "10px",
        backgroundColor: "inherit",
      }}
    >
      <h1>{data.name}</h1>
      <span />
      <div dangerouslySetInnerHTML={{ __html: data.text }} />
      <span />
      <p>{data.docketID}</p>
    </div>
  );
};

export default SearchResult;
