import { Card } from "@mui/joy";
import { CardSize } from "./GenericResultCard";

import { faker } from "@faker-js/faker";

// Helper function for random date in last 6 months
const randomRecentDate = () => {
  const date = new Date();
  date.setMonth(date.getMonth() - Math.floor(Math.random() * 6));
  return date.toISOString().split("T")[0];
};

export const generateFakeResults = (count: number) => {
  const resultTypes = ["author", "docket", "document"];
  const agencies = ["EPA", "DOT", "FDA", "USDA", "DOE"];
  const documentFormats = ["PDF", "DOCX", "HTML", "Markdown"];

  return Array.from({ length: count }, (_, i) => {
    const type = resultTypes[i % resultTypes.length];
    const base = {
      id: `fake-result-${i + 1}`,
      timestamp: randomRecentDate(),
    };

    switch (type) {
      case "author":
        return {
          ...base,
          type: "author",
          name: faker.person.fullName(),
          description: faker.person.jobTitle(),
          affiliation: faker.company.name(),
        };

      case "docket":
        const agency = agencies[Math.floor(Math.random() * agencies.length)];
        return {
          ...base,
          type: "docket",
          name: `${agency}-${faker.string.alpha({ length: 3, casing: "upper" })}-${faker.date.recent().getFullYear()}-${faker.string.numeric(4)}`,
          description: `${agency} ${fmerce.productmerce.productAdjective()} ${faker.commerce.product()} Regulations`,
          extraInfo: `Comment period closes ${faker.date.soon({ days: 30 }).toLocaleDateString()}`,
        };

      case "document":
        return {
          ...base,
          type: "document",
          name: `${faker.date.recent().getFullYear()} ${faker.commerce.department()} Report`,
          description: faker.lorem.sentence(),
          extraInfo: `${documentFormats[Math.floor(Math.random() * documentFormats.length)]}, ${faker.number.float({ min: 0.1, max: 5.0 }).toFixed(1)}MB`,
          authors: Array.from(
            { length: faker.number.int({ min: 1, max: 3 }) },
            () => faker.person.fullName(),
          ),
        };

      default:
        return base;
    }
  });
};
