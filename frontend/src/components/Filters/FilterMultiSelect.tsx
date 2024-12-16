import React, { useState } from "react";
import Select from "react-select";

const options = [
  { value: "chocolate", label: "Chocolate" },
  { value: "strawberry", label: "Strawberry" },
  { value: "vanilla", label: "Vanilla" },
];

export default function TestMultiSelect() {
  const [selectedOption, setSelectedOption] = useState(null);

  return (
    <Select
      theme={(theme) => ({
        ...theme,
        // borderRadius: 0,
        colors: {
          ...theme.colors,
          neutral0: "oklch(var(--b1))",
          neutral5: "oklch(var(--b1))",
          neutral10: "oklch(var(--b2))",
          neutral20: "oklch(var(--b2))",
          neutral30: "oklch(var(--b3))",
          neutral40: "oklch(var(--b3))",
          neutral50: "oklch(var(--b3))",
          neutral60: "oklch(var(--bc))",
          neutral70: "oklch(var(--bc))",
          neutral80: "oklch(var(--bc))",
          neutral90: "oklch(var(--bc))",
          primary25: "oklch(var(--p))",
          primary: "oklch(var(--s))",
          danger: "oklch(var(--erc))",
          dangerlight: "oklch(var(--er))",
        },
      })}
      //https://react-select.com/styles
      isMulti
      defaultValue={selectedOption}
      onChange={setSelectedOption as any}
      options={options}
    />
  );
}
