There is currently a bunch of documentation for the newer search architecture in 

/home/nicole/Documents/mycorrhizae/kessler/docs/frontend_search_current_design.md

The legacy search system is kinda sprawling throughout the codebase at the moment. But starts at the root of the project with these files that describe how the search requests hit various endpoints.

In the beginning could you just search the project and try to document the bottom part of the search architecture in an attempt to try and move everything into the new format.

/frontend/src/lib/requests/conversations.ts
/frontend/src/lib/requests/organizations.ts
/frontend/src/lib/requests/search.ts

One big part of this is the necessary conversion of the backend schemas in all these different files into a unified format that could get used by 

/home/nicole/Documents/mycorrhizae/kessler/frontend/src/components/NewSearch/GenericResultCard.tsx

Could you go ahead and throw all this documentation and throw it into docs/legacy_api_doc.md?
