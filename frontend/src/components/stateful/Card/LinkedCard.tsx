
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
        return <MediumCard data={data} enableClickAnimation={!disableHref} />;
      default:
        return <MediumCard data={data} />;
    }
  })()
  if (disableHref) {
    return rawCard
  }
  const href = calculateHref(data.object_uuid, data.type)
  // Handle click only if not clicking on an <a> tag inside the card
  const handleClick = (e: React.MouseEvent) => {
    // If the click target or any of its parents up to the div is an <a> tag, do nothing
    let target = e.target as HTMLElement | null
    while (target && target !== e.currentTarget) {
      if (target.tagName === 'A') {
        // Cancel router navigation
        return
      }
      target = target.parentElement
    }
    router.push(href)
  }

  return <div onClick={handleClick}>{rawCard}</div>

};
export default Card;
