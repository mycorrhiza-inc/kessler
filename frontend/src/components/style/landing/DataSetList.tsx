import { FileCheckIcon, PublicCommentIcon } from "@/components/style/misc/Icons";

interface StateData {
  name: string;
  abrev: string;
  comments: boolean;
  filing: boolean;
}

const StatesList = () => { };

export const SupportedStates = ({ className }: { className: string }) => {
  const DataSets = [
    // { name: "Alabama", abrev: "AL", comments: false, filing: false },
    // { name: "Alaska", abrev: "AK", comments: false, filing: false },
    // { name: "Arizona", abrev: "AZ", comments: false, filing: false },
    // { name: "Arkansas", abrev: "AR", comments: false, filing: false },
    // { name: "California", abrev: "CA", comments: false, filing: false },
    { name: "Colorado", abrev: "CO", comments: true, filing: true },
    // { name: "Connecticut", abrev: "CT", comments: false, filing: false },
    // { name: "Delaware", abrev: "DE", comments: false, filing: false },
    // { name: "Florida", abrev: "FL", comments: false, filing: false },
    // { name: "Georgia", abrev: "GA", comments: false, filing: false },
    // { name: "Hawaii", abrev: "HI", comments: false, filing: false },
    // { name: "Idaho", abrev: "ID", comments: false, filing: false },
    // { name: "Illinois", abrev: "IL", comments: false, filing: false },
    // { name: "Indiana", abrev: "IN", comments: false, filing: false },
    // { name: "Iowa", abrev: "IA", comments: false, filing: false },
    // { name: "Kansas", abrev: "KS", comments: false, filing: false },
    // { name: "Kentucky", abrev: "KY", comments: false, filing: false },
    // { name: "Louisiana", abrev: "LA", comments: false, filing: false },
    // { name: "Maine", abrev: "ME", comments: false, filing: false },
    // { name: "Maryland", abrev: "MD", comments: false, filing: false },
    // { name: "Massachusetts", abrev: "MA", comments: false, filing: false },
    // { name: "Michigan", abrev: "MI", comments: false, filing: false },
    // { name: "Minnesota", abrev: "MN", comments: false, filing: false },
    // { name: "Mississippi", abrev: "MS", comments: false, filing: false },
    // { name: "Missouri", abrev: "MO", comments: false, filing: false },
    // { name: "Montana", abrev: "MT", comments: false, filing: false },
    // { name: "Nebraska", abrev: "NE", comments: false, filing: false },
    // { name: "Nevada", abrev: "NV", comments: false, filing: false },
    // { name: "New Hampshire", abrev: "NH", comments: false, filing: false },
    // { name: "New Jersey", abrev: "NJ", comments: false, filing: false },
    // { name: "New Mexico", abrev: "NM", comments: false, filing: false },
    { name: "New York", abrev: "NY", comments: true, filing: true },
    // { name: "North Carolina", abrev: "NC", comments: false, filing: false },
    // { name: "North Dakota", abrev: "ND", comments: false, filing: false },
    // { name: "Ohio", abrev: "OH", comments: false, filing: false },
    // { name: "Oklahoma", abrev: "OK", comments: false, filing: false },
    // { name: "Oregon", abrev: "OR", comments: false, filing: false },
    // { name: "Pennsylvania", abrev: "PA", comments: false, filing: false },
    // { name: "Rhode Island", abrev: "RI", comments: false, filing: false },
    // { name: "South Carolina", abrev: "SC", comments: false, filing: false },
    // { name: "South Dakota", abrev: "SD", comments: false, filing: false },
    // { name: "Tennessee", abrev: "TN", comments: false, filing: false },
    // { name: "Texas", abrev: "TX", comments: false, filing: false },
    // { name: "Utah", abrev: "UT", comments: false, filing: false },
    // { name: "Vermont", abrev: "VT", comments: false, filing: false },
    // { name: "Virginia", abrev: "VA", comments: false, filing: false },
    // { name: "Washington", abrev: "WA", comments: false, filing: false },
    // { name: "West Virginia", abrev: "WV", comments: false, filing: false },
    // { name: "Wisconsin", abrev: "WI", comments: false, filing: false },
    // { name: "Wyoming", abrev: "WY", comments: false, filing: false },
  ];

  return (
    <div className="flex justify-center">
      <section
        id="supported-states"
        className={"pb-20 pt-15 lg:pb-25 xl:pb-30 text-center " + className}
      >
        <div className="mx-auto px-4 md:px-8 xl:px-0">
          {/* <div className="mx-auto max-w-c-1315 px-4 md:px-8 xl:px-0"> */}
          <div id="supported-states-title" className="text-center">
            <h2 className="text-4xl">
              <br />
              <br />
              Supported States
            </h2>
            <div className="flex row justify-between">
              <div className="flex row">
                <PublicCommentIcon /> = Public Comments
              </div>{" "}
              <div className="flex row">
                <FileCheckIcon /> = Filing
              </div>
            </div>
          </div>
          <div className="divider" />
          <div className="flex column columns-" id="states">
            {DataSets.map((state: StateData, _) => (
              <div
                className="flex row rounded-box border-primary border-4 outline-secondary shadow-xl p-5 w-60"
                key={state.abrev}
              >
                <b>{state.name}</b>({state.abrev})
                {state.comments ? <PublicCommentIcon /> : ""}{" "}
                {state.comments ? <FileCheckIcon /> : ""}{" "}
              </div>
            ))}
          </div>
        </div>
      </section>
    </div>
  );
};
