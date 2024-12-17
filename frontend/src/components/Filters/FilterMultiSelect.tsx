import React from "react";
import chroma from "chroma-js";

import Select, { StylesConfig } from "react-select";
import {
  ConversationsAutocompleteList,
  GeneralizedOption,
  OrganizationsAutocompleteList,
  convoAutocompleteToGeneralOption,
  orgAutocompleteToGeneralOption,
} from "./HardcodedAutocompletes";

const colourStyles: StylesConfig<GeneralizedOption, true> = {
  control: (styles) => ({ ...styles }),
  option: (styles, { data, isDisabled, isFocused, isSelected }) => {
    const color = chroma(data.color);
    return {
      ...styles,
      backgroundColor: isDisabled
        ? undefined
        : isSelected
          ? data.color
          : isFocused
            ? color.alpha(0.1).css()
            : undefined,
      color: isDisabled
        ? "#ccc"
        : isSelected
          ? chroma.contrast(color, "white") > 2
            ? "white"
            : "black"
          : data.color,
      cursor: isDisabled ? "not-allowed" : "default",

      ":active": {
        ...styles[":active"],
        backgroundColor: !isDisabled
          ? isSelected
            ? data.color
            : color.alpha(0.3).css()
          : undefined,
      },
    };
  },
  multiValue: (styles, { data }) => {
    const color = chroma(data.color);
    return {
      ...styles,
      backgroundColor: color.alpha(0.1).css(),
    };
  },
  multiValueLabel: (styles, { data }) => ({
    ...styles,
    color: data.color,
  }),
  multiValueRemove: (styles, { data }) => ({
    ...styles,
    color: data.color,
    ":hover": {
      backgroundColor: data.color,
      color: "black",
    },
  }),
};

const testOptions: GeneralizedOption[] = [
  {
    value: "chocolate",
    label: "Chocolate",
    color: "oklch(37.54% 0.0783 58.24)",
  },
  { value: "strawberry", label: "Strawberry", color: "oklch(80% 0.1 2.71)" },
  { value: "vanilla", label: "Vanilla", color: "oklch(95.64% 0.0383 58.24)" },
];
export const SimpleOptionMultiSelect = ({
  options,
}: {
  options: GeneralizedOption[];
}) => (
  <Select
    className="text-base-content"
    closeMenuOnSelect={false}
    isMulti
    options={options}
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
    styles={colourStyles}
  />
);
export const OrgMultiSelect = () => {
  const options: GeneralizedOption[] = OrganizationsAutocompleteList.map(
    orgAutocompleteToGeneralOption,
  );
  return <SimpleOptionMultiSelect options={options} />;
};

export const ConvoMultiSelect = () => {
  const options: GeneralizedOption[] = ConversationsAutocompleteList.map(
    convoAutocompleteToGeneralOption,
  );
  return <SimpleOptionMultiSelect options={options} />;
};
