There is currently a bunch of documentation for the newer search architecture in 

kessler/docs/frontend_search_current_design.md

The legacy search system is kinda sprawling throughout the codebase at the moment. But starts at the root of the project with these files that describe how the search requests hit various endpoints.

In the beginning could you just search the project and try to document the bottom part of the search architecture in an attempt to try and move everything into the new format.

kessler/frontend/src/lib/requests/conversations.ts
kessler/frontend/src/lib/requests/organizations.ts
kessler/frontend/src/lib/requests/search.ts

and document how these endpoints eventually end up getting consumed by the upper level components in 

kessler/frontend/src/components/Search/FileSearch/FileSearchView.tsx
kessler/frontend/src/components/LookupPages/ConvoLookup/ConversationTable.tsx
kessler/frontend/src/components/LookupPages/OrgLookup/OrganizationTable.tsx

One big part of this is the necessary conversion of the backend schemas in all these different files into a unified format that could get used by our filter card system. But all of that is a much lower priority, just try to understand and summarize how everythingt works and throw it in that design document.

Could you go ahead and throw all this documentation and throw it into docs/frontend_search_legacy_design_doc.md?
