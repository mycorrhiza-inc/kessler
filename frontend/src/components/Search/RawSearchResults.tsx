import React from "react";
import { SearchResult } from "@/lib/types/new_search_types";
import Card, { CardSize } from "../NewSearch/GenericResultCard";
import { isSearchOffsetsValid } from "@/lib/adapters/genericSearchCallback";
import ErrorMessage from "../messages/ErrorMessage";

export default function RawSearchResults({ data }: { data: SearchResult[] }) {
  const isValidOrder = isSearchOffsetsValid(data);
  if (!isValidOrder) {
    return (
      <ErrorMessage message="Search results returned in invalid order, not displaying out of an abbundance of caution" />
    );
  }
  return (
    <div className="flex w-full">
      <div className="grid grid-cols-1 gap-4 p-8 w-full">
        {data.map((item, index) => (
          <Card key={index} data={item} size={CardSize.Medium} />
        ))}
      </div>
    </div>
  );
}
