Last time you generated an architecture for taking the old search 

kessler/docs/frontend_search_legacy_design_doc.md

and migrating it to use the new format at 

kessler/docs/frontend_search_current_design_doc.md


implement your previous plan that you wrote down at 

/home/nicole/Documents/mycorrhizae/kessler/docs/table_migration_architecture.md

To limit scope creep, only implement this plan for the FileSearchView at 

kessler/frontend/src/components/Search/FileSearch/FileSearchView.tsx

What you should do is make a new file called FileSearchViewNew.tsx that has the same interface in order to avoid messing up the depenandices of the existing project.




--- 
There is currently a bunch of documentation for the older search architecture in 

kessler/docs/frontend_search_legacy_design_doc.md

with the addition of adapters to use the new card components like so:

kessler/frontend/src/lib/adapters/genericCardAdapters.ts

And I am wanting to transfer all the components, like


File Search View (`frontend/src/components/Search/FileSearch/FileSearchView.tsx`)
Conversation Lookup Table (`frontend/src/components/LookupPages/ConvoLookup/ConversationTable.tsx`)
and 
Organization Lookup Table (`frontend/src/components/LookupPages/OrgLookup/OrganizationTable.tsx`)

to use the new search architecture described in:

kessler/docs/frontend_search_current_design.md

For now I dont want you to actually write any code, instead think about the architecture for how all of this would work, and write your thoughts to docs/table_rearchitect.md 

To do this you should use the existing fetchers, and combine them with a genericCardAdapters to create the lookup search functionality needed to create the callbacks for each method.

Also all of these pages will absolutely require the SSR functionality, this is already implemented for 
kessler/frontend/src/app/(application)/search/page.tsx
if you need a reference for how its implemented. But it just uses traditional React Server Components and doesnt need anything nextjs specific. But the functionality for this component is still a bit specific, so if you can think of any simplifications or generalizations please include them.


so if you can think of any simplifications to that please include in your thoughts on the architecture.

but for now just ignore that and architect out how you would write the adapters, and throw your final architecture plan in docs/generic_card_architecture.md
