import { clsx } from "clsx";
import React from "react";

enum CardType {
  Author = "author",
  Docket = "docket",
  Document = "document",
}

interface BaseCardData {
  name: string;
  description: string;
  timestamp: string;
  authors?: Array<string>;
  extraInfo?: string;
}

interface AuthorCardData extends BaseCardData {
  type: CardType.Author;
}

interface DocketCardData extends BaseCardData {
  type: CardType.Docket;
}

interface DocumentCardData extends BaseCardData {
  type: CardType.Document;
}

type CardData = AuthorCardData | DocketCardData | DocumentCardData;

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
            "w-4 h-4 rounded flex justify-center items-center",
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
            "w-4 h-4 rounded flex justify-center items-center",
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
            "w-4 h-4 rounded flex justify-center items-center",
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

enum CardSize {
  Medium = "medium",
  Large = "large",
  Small = "small",
}
// For this component add support for displaying the card in additional sizes.
// For medium keep everything the same.
// For small collapse it so that all the card info can be displayed on a single line of text. Designed for instances where you need to display hundreds of results in a compact list.
// For large design try to render it instead like its the main header on the page, make everything as big as possible and label each section with what it means. For example a page for a docket might have a large card at the top of the page, with a search interface for its documents below it.
//
const Card: React.FC<{ data: CardData; size?: CardSize }> = ({
  data,
  size,
}) => {
  if (!size) {
    size = CardSize.Medium;
  }
  return (
    <div
      className={clsx("card shadow-xl mb-4", {
        "bg-base-200 p-4": size === CardSize.Medium,
        "bg-base-100 p-2 text-sm": size === CardSize.Small,
        "bg-base-100 p-6 rounded-lg": size === CardSize.Large,
      })}
    >
      <div
        className={clsx("flex justify-between items-center", {
          "mb-2": size !== CardSize.Small,
          "gap-2": size === CardSize.Small,
        })}
      >
        <div className="flex items-center gap-2">
          {size !== CardSize.Small && getIcon(data.type)}
          <h2
            className={clsx({
              "card-title": true,
              "text-lg": size === CardSize.Large,
              "text-base": size === CardSize.Medium,
              "text-sm": size === CardSize.Small,
            })}
          >
            {size === CardSize.Large && (
              <span className="text-gray-500 mr-2">Name:</span>
            )}
            {data.name}
          </h2>
        </div>
        <div
          className={clsx(
            `${getTypeColor(data.type)} text-white px-2 py-1 rounded capitalize`,
            {
              "text-xs": size !== CardSize.Large,
              "text-sm": size === CardSize.Large,
            },
          )}
        >
          {data.type}
        </div>
      </div>

      {size !== CardSize.Small && (
        <div
          className={clsx({
            "mb-4": size === CardSize.Medium,
            "mb-6": size === CardSize.Large,
          })}
        >
          <p
            className={clsx({
              "text-sm": size === CardSize.Medium,
              "text-base": size === CardSize.Large,
              "text-gray-600": size === CardSize.Large,
            })}
          >
            {size === CardSize.Large && (
              <span className="text-gray-500 mr-2">Description:</span>
            )}
            {data.description}
          </p>
          {data.extraInfo && (
            <p
              className={clsx({
                "text-xs text-gray-500": size === CardSize.Medium,
                "text-sm text-gray-600": size === CardSize.Large,
              })}
            >
              {size === CardSize.Large && (
                <span className="text-gray-500 mr-2">Details:</span>
              )}
              {data.extraInfo}
            </p>
          )}
        </div>
      )}

      <div
        className={clsx("flex justify-between items-center", {
          "text-xs": size !== CardSize.Large,
          "text-sm": size === CardSize.Large,
        })}
      >
        <span className="text-gray-500">
          {size === CardSize.Large && (
            <span className="mr-2">Timestamp Type:</span>
          )}
          {getTimestampLabel(data.type)}
        </span>
        <span>
          {size === CardSize.Large && (
            <span className="text-gray-500 mr-2">Date:</span>
          )}
          {data.timestamp}
        </span>
      </div>

      {data.authors && size !== CardSize.Small && (
        <div
          className={clsx({
            "mt-4": size === CardSize.Medium,
            "mt-6": size === CardSize.Large,
          })}
        >
          <h3
            className={clsx("font-bold", {
              "text-sm mb-2": size === CardSize.Medium,
              "text-base mb-4": size === CardSize.Large,
            })}
          >
            Authors
          </h3>
          <div
            className={clsx("bg-pink-100 rounded", {
              "p-2": size === CardSize.Medium,
              "p-4": size === CardSize.Large,
            })}
          >
            {data.authors.map((author, index) => (
              <span
                key={index}
                className={clsx({
                  "text-sm": size === CardSize.Medium,
                  "text-base": size === CardSize.Large,
                })}
              >
                {author}
                {index !== data.authors!.length - 1 && ", "}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default Card;
