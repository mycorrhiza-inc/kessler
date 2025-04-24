"use client";
import Card from "@/components/NewSearch/GenericResultCard";

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
  {
    type: "author",
    name: "Michael Chen",
    description: "Environmental policy analyst",
    timestamp: "2024-04-10",
  },
  {
    type: "docket",
    name: "DOT-OST-2024-0045",
    description: "Vehicle Fuel Economy Standards",
    timestamp: "2024-04-18",
    extraInfo: "Proposed rulemaking",
  },
];

export default function Page() {
  return (
    <div className="grid grid-cols-1 gap-4 p-8">
      {exampleData.map((data, index) => (
        <Card key={index} data={data} />
      ))}
    </div>
  );
}
