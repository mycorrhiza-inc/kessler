import React from "react";

type CardType = "author" | "docket" | "document";

interface BaseCardData {
  name: string;
  description: string;
  timestamp: string;
  authors?: Array<string>;
  extraInfo?: string;
}

interface AuthorCardData extends BaseCardData {
  type: "author";
}

interface DocketCardData extends BaseCardData {
  type: "docket";
}

interface DocumentCardData extends BaseCardData {
  type: "document";
}

type CardData = AuthorCardData | DocketCardData | DocumentCardData;

const getTimestampLabel = (type: string): string => {
  switch (type) {
    case "author":
      return "last active";
    case "docket":
      return "last update";
    case "document":
      return "filing date";
    default:
      return "";
  }
};

const Card: React.FC<{ data: CardData }> = ({ data }) => {
  const getIcon = (type: string) => {
    switch (type) {
      case "author":
        return <div className="w-4 h-4 bg-purple-500 rounded-full"></div>;
      case "docket":
        return <div className="w-4 h-4 bg-green-500 rounded-full"></div>;
      case "document":
        return (
          <div className="w-4 h-4 bg-red-500 rounded flex justify-center items-center">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-3 w-3 text-white"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fillRule="evenodd"
                d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z"
                clipRule="evenodd"
              />
            </svg>
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <div className="card bg-base-100 shadow-md p-4 mb-4">
      <div className="flex items-center mb-2">
        {getIcon(data.type)}
        <h2 className="card-title ml-2">{data.name}</h2>
      </div>
      <div className="mb-4">
        <p className="text-sm">{data.description}</p>
        {data.extraInfo && (
          <p className="text-xs text-gray-500">{data.extraInfo}</p>
        )}
      </div>
      <div className="flex justify-between items-center">
        <span className="text-xs text-gray-500">
          {getTimestampLabel(data.type)}
        </span>
        <span className="text-xs">{data.timestamp}</span>
      </div>
      {data.authors && (
        <div className="mt-4">
          <h3 className="text-sm font-bold mb-2">Authors</h3>
          <div className="bg-pink-100 p-2 rounded">
            {data.authors.map((author, index) => (
              <span key={index} className="text-sm">
                {author}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default Card;
