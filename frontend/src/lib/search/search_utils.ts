import { faker } from "@faker-js/faker";
import {
  PaginationData,
  SearchResult,
  SearchResultsGetter,
} from "../types/new_search_types";
import { AuthorCardData, CardData, CardType, DocketCardData, DocumentCardData } from "../types/generic_card_types";
import { AuthorInformation } from "../types/backend_schemas";
import { sleep } from "../utils";
const randomRecentDate = () => {
  const date = new Date();
  date.setMonth(date.getMonth() - Math.floor(Math.random() * 6));
  return date.toISOString().split("T")[0];
};

export const generateFakeResults = async ({ page, limit }: PaginationData) => {
  await sleep(3000);
  return generateFakeResultsRaw(limit);
};

export const generateFakeResultsRaw = (count: number): CardData[] => {
  const resultTypes = ["author", "docket", "document"];
  const agencies = ["EPA", "DOT", "FDA", "USDA", "DOE"];
  const documentFormats = ["PDF", "DOCX", "HTML", "Markdown"];

  return Array.from({ length: count }, (_, i) => {
    const type =
      resultTypes[Math.floor(Math.random() * 100) % resultTypes.length];
    const base = {
      id: `fake-result-${i + 1}`,
      object_uuid: faker.string.uuid(),
      index: 0,
      timestamp: randomRecentDate(),
    };

    switch (type) {
      case "author":
        const author_data: AuthorCardData = {
          ...base,
          type: CardType.Author,
          name: faker.person.fullName(),
          description: faker.person.jobTitle(),
          // affiliation: faker.company.name(),
        };
        return author_data;

      case "docket":
        const agency = agencies[Math.floor(Math.random() * agencies.length)];
        const docket_data: DocketCardData = {
          ...base,
          type: CardType.Docket,
          name: `${agency}-${faker.string.alpha({ length: 3, casing: "upper" })}-${faker.date.recent().getFullYear()}-${faker.string.numeric(4)}`,
          description: `${agency} ${faker.commerce.product()} Regulations`,
          extraInfo: `Comment period closes ${faker.date.soon({ days: 30 }).toLocaleDateString()}`,
        };
        return docket_data;

      case "document":
        const document_data: DocumentCardData = {
          ...base,
          type: CardType.Document,
          name: `${faker.date.recent().getFullYear()} ${faker.commerce.department()} Report`,
          description: faker.lorem.sentence(),
          extraInfo: `${documentFormats[Math.floor(Math.random() * documentFormats.length)]}, ${faker.number.float({ min: 0.1, max: 5.0 }).toFixed(1)}MB`,
          authors: Array.from(
            { length: faker.number.int({ min: 1, max: 3 }) },
            () => {
              const data: AuthorInformation = { author_name: faker.person.fullName(), is_person: true, is_primary_author: false, author_id: faker.string.uuid() }
              return data
            },
          ),
        };
        return document_data;

      default:
        throw new Error("wrong card type")
    }
  });
};
