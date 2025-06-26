import { CardData, CardType } from "@/lib/types/generic_card_types";
import { clsx } from "clsx";
import React, { ReactNode } from "react";
import { AuthorPill } from "../Pills/TextPills";
import { AuthorInformation } from "@/lib/types/backend_schemas";

// Card sizes
export enum CardSize {
  Medium = "medium",
  Large = "large",
  Small = "small",
}

// Helper functions
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
  const baseClasses = clsx(
    getTypeColor(type),
    "w-4 h-4 rounded-sm flex justify-center items-center"
  );
  switch (type) {
    case CardType.Author:
      return (
        <div className={baseClasses}>
          <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 text-white" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
          </svg>
        </div>
      );
    case CardType.Docket:
      return (
        <div className={baseClasses}>
          <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 text-white" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M4 4a2 2 0 012-2h8a2 2 0 012 2v12a1 1 0 110 2h-3a1 1 0 01-1-1v-2a1 1 0 00-1-1H9a1 1 0 00-1 1v2a1 1 0 01-1 1H4a1 1 0 110-2V4z" clipRule="evenodd" />
          </svg>
        </div>
      );
    case CardType.Document:
      return (
        <div className={baseClasses}>
          <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 text-white" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clipRule="evenodd" />
          </svg>
        </div>
      );
    default:
      return null;
  }
};

// Reusable subcomponents
const CardBadge: React.FC<{ type: CardType; size: CardSize }> = ({ type, size }) => (
  <div
    className={clsx(
      getTypeColor(type),
      "text-white rounded-sm capitalize",
      size === CardSize.Large ? "px-4 py-2 text-lg" : size === CardSize.Small ? "px-1 py-0.5 text-xs" : "px-2 py-1 text-xs"
    )}
  >
    {type}
  </div>
);

const CardHeader: React.FC<{ data: CardData; size: CardSize }> = ({ data, size }) => (
  <div className={clsx(
    "flex items-center justify-between",
    size === CardSize.Large ? "mb-4" : "mb-2"
  )}>
    <div className="flex items-center">
      {getIcon(data.type)}
      <h2 className={clsx(
        size === CardSize.Large ? "text-4xl font-bold ml-4" : "ml-2 text-lg font-semibold"
      )}>
        {data.name}
      </h2>
    </div>
    <CardBadge type={data.type} size={size} />
  </div>
);

const CardDescription: React.FC<{ data: CardData; size: CardSize }> = ({ data, size }) => {
  const className = clsx(
    size === CardSize.Large
      ? "text-lg mb-4"
      : size === CardSize.Medium
        ? "text-sm mb-4"
        : "text-xs flex-1 truncate"
  );
  let content = data.description;
  if (size === CardSize.Small && content.length > 100) {
    content = content.slice(0, 100) + "…";
  }
  if (size === CardSize.Medium && content.length > 300) {
    content = content.slice(0, 300) + "…";
  }
  return <p className={className}>{content}</p>;
};

const CardFooter: React.FC<{ data: CardData; size: CardSize }> = ({ data, size }) => (
  <div className={clsx(
    "flex items-center justify-between text-gray-500",
    size === CardSize.Large ? "text-lg" : "text-xs"
  )}>
    <span>
      {size === CardSize.Large ? `${getTimestampLabel(data.type)}:` : getTimestampLabel(data.type)}
    </span>
    <span>
      {/* {size === CardSize.Large ? data.timestamp : `${data.timestamp} - idx:${data.index}`} */}
      {data.timestamp}
    </span>
  </div>
);

const CardAuthors: React.FC<{ authors: AuthorInformation[]; size: CardSize }> = ({ authors, size }) => (
  <div className={clsx(size === CardSize.Large ? "mt-6" : "mt-4")}>
    <h3 className={clsx(
      size === CardSize.Large ? "text-2xl font-bold mb-2" : "text-sm font-semibold mb-2"
    )}>
      Authors
    </h3>
    <div className={clsx(
      size === CardSize.Large ? "flex flex-wrap gap-4" : "flex flex-wrap gap-2 bg-pink-100 p-2 rounded-sm"
    )}>
      {authors.map((author, idx) => (
        <AuthorPill key={idx} author={author} />
      ))}
    </div>
  </div>
);

// Card variants
const SmallCard: React.FC<{ data: CardData }> = ({ data }: { data: CardData }) => (
  <div className="card bg-base-200 shadow p-2 mb-2">
    <div className="flex items-center space-x-2 text-xs">
      {getIcon(data.type)}
      <span className="font-semibold truncate" title={data.name}>{data.name}</span>
      <CardBadge type={data.type} size={CardSize.Small} />
      <CardDescription data={data} size={CardSize.Small} />
      <CardFooter data={data} size={CardSize.Small} />
    </div>
  </div>
);

const MediumCard: React.FC<{ data: CardData }> = ({ data }: { data: CardData }) => (
  <div className="card bg-base-200 shadow-xl p-4 mb-4">
    <CardHeader data={data} size={CardSize.Medium} />
    <CardDescription data={data} size={CardSize.Medium} />
    <CardFooter data={data} size={CardSize.Medium} />
    {data.type === CardType.Document && data.authors && (
      <CardAuthors authors={data.authors} size={CardSize.Medium} />
    )}
  </div>
);

const LargeCard: React.FC<{ data: CardData }> = ({ data }: { data: CardData }) => (
  <div className="card bg-base-200 shadow-2xl p-8 mb-8 top-0 z-10 w-full max-w-full">
    <CardHeader data={data} size={CardSize.Large} />
    <div className="mb-6">
      <h3 className="text-2xl font-bold mb-2">Description</h3>
      <p className="text-lg">{data.description}</p>
    </div>
    {data.extraInfo && (
      <div className="mb-6">
        <h3 className="text-2xl font-bold mb-2">Extra Info</h3>
        <p className="text-lg text-gray-700">{data.extraInfo}</p>
      </div>
    )}
    <CardFooter data={data} size={CardSize.Large} />
    {data.type === CardType.Document && data.authors && (
      <CardAuthors authors={data.authors} size={CardSize.Large} />
    )}
  </div>
);



const Card = ({ data, disableHref, size = CardSize.Medium }: { data: CardData, disableHref?: boolean, size?: CardSize }) => {
  const rawCard: ReactNode = (() => {
    switch (size) {
      case CardSize.Large:
        return <LargeCard data={data} />;
      case CardSize.Small:
        return <SmallCard data={data} />;
      case CardSize.Medium:
        return <MediumCard data={data} />;
      default:
        return <MediumCard data={data} />;
    }
  })()
  if (disableHref) {
    return rawCard
  }

  return rawCard

};

export default Card;
export { LargeCard, MediumCard, SmallCard };
export type { CardData };
