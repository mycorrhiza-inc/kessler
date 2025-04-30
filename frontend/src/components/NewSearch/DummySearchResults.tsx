import Card, { CardSize } from "./GenericResultCard";

const exampleData = [
  {
    type: "author",
    name: "Dr. Sarah Johnson",
    description: "Lead researcher on climate policy",
    timestamp: "2024-03-15",
  },
  {
    type: "docket",
    name: "EPA-HQ-OAR-2023-1234",
    description: "Clean Air Act Section 111(d) Emissions Guidelines",
    timestamp: "2024-04-01",
    extraInfo: "Comment period closes June 30",
  },
  {
    type: "document",
    name: "2023 Annual Emissions Report",
    description: "National greenhouse gas inventory update",
    timestamp: "2024-01-15",
    extraInfo: "PDF, 1.2MB",
    authors: ["EPA Office of Air and Radiation", "Climate Analysis Team"],
  },
];

export default function DummyResults() {
  return (
    <div className="flex w-full">
      {/* Main search results content */}
      <div className="grid grid-cols-1 gap-4 p-8 w-full">
        {exampleData.map((data, index) => (
          <Card key={index} data={data} size={CardSize.Medium} />
        ))}
      </div>
    </div>
  );
}
