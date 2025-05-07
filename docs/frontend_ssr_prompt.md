So I have this problem. Since this search functionality is shared all throughout the app its split up to be at the main page holding all the page logic 

frontend/src/app/(application)/search/page.tsx

a hook containing a bunch of the state logic
frontend/src/lib/hooks/useSearchState.ts

and another component handling the infinite scroll
frontend/src/components/Search/SearchResults.tsx

However I have a bit of a problem. For the main search page, and potentially in other places, we would like the search to behave in the following way.

1. The user sends a request for the page.

2. The server recieves the page request, sends a search request to the database.

3. The user recieves the rest of the page, with the search results behind a <Suspense> component.

4. The server recieves the search results, renders the html for the search results on the server, and streams it to the user.

5. The user recieves the html for the search results, they load and become interactive. And at the same time the infinite scroll and search functionality start to work. So that if the user scrolls to the bottom of the search results it automatically loads more. And if they enter a new query or adjust the filters it clears away all results, including the server pregenerated html, and begins the search using the current logic.

How should I implement this?
