import { CardData, CardType } from "@/lib/types/generic_card_types";
import { clsx } from "clsx";
import React from "react";
import { AuthorPill } from "../Pills/TextPills";

const getTimestampLabel = (type: CardType): string => {
  switch (type) {
    case CardType.Author:
      return "last active";
    case CardType.Docket:
      return "last update";
    case CardType.Document:
      return "filing date";
    default:
      return "";
  }
};

const getTypeColor = (type: CardType): string => {
  switch (type) {
    case CardType.Author:
      return "bg-purple-500";
    case CardType.Docket:
      return "bg-green-500";
    case CardType.Document:
      return "bg-red-500";
    default:
      return "";
  }
};

const getIcon = (type: CardType) => {
  switch (type) {
    case CardType.Author:
      return (
        <div
          className={clsx(
            getTypeColor(type),
            "w-4 h-4 rounded-sm flex justify-center items-center",
          )}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-3 w-3 text-white"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fillRule="evenodd"
              d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
              clipRule="evenodd"
            />
          </svg>
        </div>
      );
    case CardType.Docket:
      return (
        <div
          className={clsx(
            getTypeColor(type),
            "w-4 h-4 rounded-sm flex justify-center items-center",
          )}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-3 w-3 text-white"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fillRule="evenodd"
              d="M4 4a2 2 0 012-2h8a2 2 0 012 2v12a1 1 0 110 2h-3a1 1 0 01-1-1v-2a1 1 0 00-1-1H9a1 1 0 00-1 1v2a1 1 0 01-1 1H4a1 1 0 110-2V4z"
              clipRule="evenodd"
            />
          </svg>
        </div>
      );
    case CardType.Document:
      return (
        <div
          className={clsx(
            getTypeColor(type),
            "w-4 h-4 rounded-sm flex justify-center items-center",
          )}
        >
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

export enum CardSize {
  Medium = "medium",
  Large = "large",
  Small = "small",
}
// For this component add support for displaying the card in additional sizes.
// For medium keep everything the same.
// For small collapse it so that all the card info can be displayed on a single line of text. Designed for instances where you need to display hundreds of results in a compact list.
// For large design try to render it instead like its the main header on the page, make everything as big as possible and label each section with what it means. For example a page for a docket might have a large card at the top of the page, with a search interface for its documents below it.
const Card: React.FC<{ data: CardData; size?: CardSize }> = ({
  data,
  size,
}) => {
  if (!size) {
    size = CardSize.Small;
  }
  if (size == CardSize.Medium) {
    return <MediumCard data={data} />;
  }
  if (size == CardSize.Small) {
    return <MediumCard data={data} />;
  }
  if (size == CardSize.Large) {
    return <MediumCard data={data} />;
  }
};

const MediumCard: React.FC<{ data: CardData }> = ({ data }) => {
  return (
    <div className="card bg-base-200 shadow-xl p-4 mb-4">
      <div className="flex justify-between items-center mb-2">
        <div className="flex items-center">
          {getIcon(data.type)}
          <h2 className="card-title ml-2">{data.name}</h2>
        </div>
        <div
          className={`${getTypeColor(data.type)} text-white px-2 py-1 rounded-sm text-xs capitalize`}
        >
          {data.type}
        </div>
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
        <span className="text-xs">{`${data.timestamp} - index:${data.index}`}</span>
      </div>
      {data.authors && (
        <div className="mt-4">
          <h3 className="text-sm font-bold mb-2">Authors</h3>
          <div className="bg-pink-100 p-2 rounded-sm">
            {data.authors.map((author, index) => (
              <AuthorPill author={author} />
            ))}
          </div>
        </div>
      )}
    </div>
  );
};


export default Card;

// Export types for adapter layer
export type { CardData };
