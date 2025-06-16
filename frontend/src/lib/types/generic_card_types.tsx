import { AuthorInformation } from "./backend_schemas";

export enum CardType {
  Author = "author",
  Docket = "docket",
  Document = "document",
}

export interface BaseCardData {
  name: string;
  description: string;
  timestamp: string;
  extraInfo?: string;
  index: number;
}


export interface AuthorCardData extends BaseCardData {
  type: CardType.Author;
}

export interface DocketCardData extends BaseCardData {
  type: CardType.Docket;
}

export interface DocumentCardData extends BaseCardData {
  type: CardType.Document;
  authors: {
    author_name: string;
    is_person: boolean;
    is_primary_author: boolean;
    author_id: string;
  }[];
}

export type CardData = AuthorCardData | DocketCardData | DocumentCardData;

