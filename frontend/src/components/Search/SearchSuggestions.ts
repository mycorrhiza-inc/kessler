import { getRuntimeEnv } from "@/lib/env_variables_hydration_script";
import { PageContextMode, Suggestion } from "@/lib/types/SearchTypes";

export const getRawSuggestions = (
  PageContext: PageContextMode,
): Suggestion[] => {
  if (PageContext === PageContextMode.Conversations) {
    return [
      {
        id: "fc001a23-4b5c-6d7e-8f9g-0h1i2j3k4l5",
        type: "nypuc_docket_industry",
        label: "Miscellaneous",
        value: "Miscellaneous",
      },
      {
        id: "fc002b34-5c6d-7e8f-9g0h-1i2j3k4l5m6",
        type: "nypuc_docket_industry",
        label: "Gas",
        value: "Gas",
      },
      {
        id: "fc003c45-6d7e-8f9g-0h1i-2j3k4l5m6n7",
        type: "nypuc_docket_industry",
        label: "Electric",
        value: "Electric",
      },
      {
        id: "fc004d56-7e8f-9g0h-1i2j-3k4l5m6n7o8",
        type: "nypuc_docket_industry",
        label: "Facility Gen.",
        value: "Facility Gen.",
      },
      {
        id: "fc005e67-8f9g-0h1i-2j3k-4l5m6n7o8p9",
        type: "nypuc_docket_industry",
        label: "Transmission",
        value: "Transmission",
      },
      {
        id: "fc006f78-9g0h-1i2j-3k4l-5m6n7o8p9q0",
        type: "nypuc_docket_industry",
        label: "Water",
        value: "Water",
      },
      {
        id: "fc007g89-0h1i-2j3k-4l5m-6n7o8p9q0r1",
        type: "nypuc_docket_industry",
        label: "Communication",
        value: "Communication",
      },
    ];
  }
  if (PageContext === PageContextMode.Organizations) {
    return [];
  }
  if (PageContext === PageContextMode.Files) {
    const file_extensions = [
      {
        id: "eb096148-7944-4f02-8c7b-16d0d8549e91",
        type: "extension",
        label: "pdf",
        value: "pdf",
      },
      {
        id: "2e809511-33ac-4a5b-a3bd-14fdaa8694e1",
        type: "extension",
        label: "xlsx",
        value: "xlsx",
      },
      {
        id: "ec05fe86-4ab6-415f-a3a5-a7c724753a8c",
        type: "extension",
        label: "docx",
        value: "docx",
      },
    ];
    const file_types = [
      {
        id: "fc001a23-5f7e-4b3c-9d2a-8f6e4c7d9e0b",
        type: "file_class",
        label: "Plans and Proposals",
        value: "Plans and Proposals",
      },
      {
        id: "fc002b34-6a8d-4c5f-ae3b-9g7f5d8e0f1c",
        type: "file_class",
        label: "Corrospondence",
        value: "Corrospondence",
      },
      {
        id: "fc003c45-7b9e-5d6g-bf4c-ah8g6e9f1g2d",
        type: "file_class",
        label: "Exhibits",
        value: "Exhibits",
      },
      {
        id: "fc004d56-8c0f-6e7h-cg5d-bi9h7f0g2h3e",
        type: "file_class",
        label: "Testimony",
        value: "Testimony",
      },
      {
        id: "fc005e67-9d1g-7f8i-dh6e-cj0i8g1h3i4f",
        type: "file_class",
        label: "Reports",
        value: "Reports",
      },
      {
        id: "fc006f78-0e2h-8g9j-ei7f-dk1j9h2i4j5g",
        type: "file_class",
        label: "Comments",
        value: "Comments",
      },
      {
        id: "fc007g89-1f3i-9h0k-fj8g-el2k0i3j5k6h",
        type: "file_class",
        label: "Attachment",
        value: "Attachment",
      },
      {
        id: "e1ed3a38-7164-47a4-ad7f-5bd3715ec894",
        type: "file_class",
        label: "Rulings",
        value: "Rulings",
      },
      {
        id: "b4c9e0e0-5510-4cb0-a5f8-ead7ebbc61a8",
        type: "file_class",
        label: "Orders",
        value: "Orders",
      },
      {
        id: "3693bac9-f2c8-4f5a-9f8b-b023ad029f93",
        type: "file_class",
        label: "Transcripts",
        value: "Transcripts",
      },
      {
        id: "040f8217-2afd-4939-8117-35bbf333cfc1",
        type: "file_class",
        label: "Letters",
        value: "Letters",
      },
    ];
    return file_extensions.concat(file_types);
  }
  console.error(
    "Unknown page context for generating raw suggestions",
    PageContext,
  );
  return [];
};

interface RawSuggestion {
  id: string;
  name: string;
  type: string;
}

const rawToRealSuggestions = (sug: RawSuggestion): Suggestion => {
  return { id: sug.id, type: sug.type, label: sug.name, value: sug.id };
};

export const fetchSuggestionsQuickwitAsync = async (
  query: string,
): Promise<Suggestion[]> => {
  const runtimeClientConfig = getRuntimeEnv();
  // IF issues replace this line
  // const apiUrl =
  // runtimeClientConfig?.public_api_url || "https://api.kessler.xyz";
  const apiUrl = runtimeClientConfig?.public_api_url;
  const url = `${apiUrl}/v2/autocomplete/files-basic?query=${query}`;
  const res = await fetch(url);
  const suggestions = await res.json();
  const return_sugs = (suggestions as RawSuggestion[]).map(
    rawToRealSuggestions,
  );
  return return_sugs;
};

export const mockFetchSuggestions = async (
  query: string,
  PageContext: PageContextMode,
): Promise<Suggestion[]> => {
  // Simulate API delay
  // await new Promise((resolve) => setTimeout(resolve, 300));

  let suggestions: Suggestion[] = getRawSuggestions(PageContext).filter(
    (s) =>
      s.label.toLowerCase().includes(query.toLowerCase()) ||
      s.type.toLowerCase().includes(query.toLowerCase()),
  );
  if (query.length < 3) return suggestions;
  const new_sugs = await fetchSuggestionsQuickwitAsync(query);
  suggestions = [...suggestions, ...new_sugs];

  return suggestions;
};
