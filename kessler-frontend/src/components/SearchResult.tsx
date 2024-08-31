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
    <div className="standard-box">
      <h1>{data.name}</h1>
      <span />
      <div dangerouslySetInnerHTML={{ __html: data.text }} />
      <span />
      <p>{data.docketID}</p>
    </div>
  );
};

export default SearchResult;
