
'use client'
import { CardData, CardType } from "@/lib/types/generic_card_types";
import { useRouter } from 'next/navigation'
import React, { ReactNode } from "react";
import { LargeCard, MediumCard, SmallCard, CardSize } from "@/components/style/cards/SizedCards";

const calculateHref = (objectId: string, objectType: CardType) => {
  // Cuts of any segmentation info added on at the end of the object_id
  objectId = objectId.slice(0, 36)
  switch (objectType) {
    case CardType.Author:
      return `/orgs/${objectId}`
    case CardType.Docket:
      return `/docket/${objectId}`
    case CardType.Document:
      return `/filling/${objectId}`
  }
}

const Card = ({ data, disableHref, size = CardSize.Medium }: { data: CardData, disableHref?: boolean, size?: CardSize }) => {
  const router = useRouter();
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
  const href = calculateHref(data.object_uuid, data.type)
  return <div onClick={() => router.push(href)}>{rawCard}</div>

};
export default Card;
