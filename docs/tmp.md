There is currently a bunch of documentation for the older search architecture in 

kessler/docs/frontend_search_legacy_design_doc.md

for this part I want you to architect out how you would write addapters to convert the data from the following api endpoints

kessler/frontend/src/lib/requests/conversations.ts
kessler/frontend/src/lib/requests/organizations.ts
kessler/frontend/src/lib/requests/search.ts

And figure out how you would write adapters to convert them into data that could be used by the generic cards at 

kessler/frontend/src/components/NewSearch/GenericResultCard.tsx

For some broader context, we are wanting to convert the entire architecture into a new format 

kessler/docs/frontend_search_current_design.md

but for now just ignore that and architect out how you would write the adapters, and throw your final architecture plan in docs/generic_card_architecture.md
