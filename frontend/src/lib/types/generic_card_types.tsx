export enum CardType {
  Author = "author",
  Docket = "docket",
  Document = "document",
}

export interface BaseCardData {
  name: string;
  description: string;
  timestamp: string;
  authors?: Array<string>;
  extraInfo?: string;
}

export interface AuthorCardData extends BaseCardData {
  type: CardType.Author;
}

export interface DocketCardData extends BaseCardData {
  type: CardType.Docket;
}

export interface DocumentCardData extends BaseCardData {
  type: CardType.Document;
}

export type CardData = AuthorCardData | DocketCardData | DocumentCardData;
