import React from 'react';
import { SearchResult } from '@/lib/types/new_search_types';
import Card, { CardSize } from '../NewSearch/GenericResultCard';

export default function RawSearchResults({ data }: { data: SearchResult[] }) {
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
