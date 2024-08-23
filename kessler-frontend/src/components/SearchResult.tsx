import { Card } from "@mui/joy";

type SearchFields = {
  id: string;
  name: string;
  text: string;
  docketID: string;
};

type SearchResultProps = {
  data: SearchFields;
};

const SearchResult = ({ data }: SearchResultProps) => {
  return (
    <Card
      style={{
        padding: "15px",
        border: "1px solid grey",
        borderRadius: "10px",
        backgroundColor: "inherit",
        width: "90%",
        maxHeight: "15em",
      }}
    >
      <h1>{data.name}</h1>
      <span />
      <div dangerouslySetInnerHTML={{ __html: data.text }} />
      <span />
      <p>{data.docketID}</p>
    </Card>
  );
};

export default SearchResult;
