So I have this problem. Since this search functionality is shared all throughout the app its split up to be at the main page holding all the page logic 

frontend/src/app/(application)/search/page.tsx

a hook containing a bunch of the state logic
frontend/src/lib/hooks/useSearchState.ts

and another component handling the infinite scroll
frontend/src/components/Search/SearchResults.tsx

However we have a bit of a problem, and a bit of a conflicting requirement.

1. It needs to be fully generic across an absolute ton of search interfaces. Including stuff like:
- Handling the search bar on the top of the home page `/`, including changing the URL to `/search?query=<input>` when that happens.
- Handling the search for routes like `/search?query=<input>`, and on the initial load of that page, executing a search for these initial components. And then having the same functionality as the regular search bar, including infinite search and updating the query in the url.
- A Command-K search that exists on every page, and has the same functionality for queries, infinite scroll. It should also change the url while searching, to make sharing a search result easy. But after the modal is closed it should remove all history from the browser and return to the page from before.
- (IN PROGRESS) Document specific search on individual pages. This should essentially behave exactly like the search on `/search` except it should automatically include filters for specific things. (Like only seeing documents made by that org, and excluding things like cases, or other authors from appearing in results.) (If you change the filters or search params on this page, it shouldnt change the url to `/search`)
- (IN PROGRESS) It should also handle other stuff like displaying a list of recently updated documents. Even though this doesnt share a lot of functionality, like needing to update on new queries or filters. It does need infinite scroll and hits the same API'
s as everything else. So doing this is an important simplifying assumption for the whole project.

2. The page needs access to a bunch of hard dynamic features. These include stuff like:
- Under some situations changing the url of a page, for example: On the home page when you use the search bar, everything else on the bottom of the page animates away. And the new results will come up, and when closed the page will animate back to its original state. Similar stuff going on with Command-K search as well.
- Initially only getting a small subset of search features, and when the user scrolls to the bottom incrementally fetching more search results, creating an infinite scroll.
- Having a dynamic search and filter setting system, so that the searches execute and change depending on how you adjust your filters and queries.


3. (In Progress) It needs to have MAXIMUM PERFORMANCE with server side rendering. Under the current client side architecture, in order to actually show search results the javascript needs to get delivered to the page, then and only then can it make an API request to get the rest of the data. However for a bunch of different locations like `/search?query=<input>` this introduces a bunch of inefficiency, since our nextjs is directly next to our API server, so it could just ask for the data directly from the api server. Render the componet there as a React Server Component, and stream it down to the client. This is a hugely important performance optimisation expecially for the following 2 use cases 
- Individual organization that shows all the documents that an org has authored. Which means that every time you ask for an org page, you need to wait for a waterfall and an long query just to see the documents. If this could be rendered on the server and cached it could immediately show up and provide a super good experience.
- The final use case of displaying a list of recently updated docs on the homepage absolutely requires this, otherwise it costs to much and takes way to long.

However it does run into some problems, mainly all of these use cases require infinite scroll. And in react getting these components that need a bunch of really heavy client state to work with infinite scroll. Luckily the API you need to hit for both is accessible from both the client and server.

Here was some of my ideas for how to implement it.

1. The user sends a request for the page.

2. The server recieves the page request, sends a search request to the database.

3. The user recieves the rest of the page, with the search results behind a <Suspense> component.

4. The server recieves the search results, renders the html for the search results on the server, and streams it to the user.

5. The user recieves the html for the search results, they load and become interactive. And at the same time the infinite scroll and search functionality start to work. So that if the user scrolls to the bottom of the search results it automatically loads more. And if they enter a new query or adjust the filters it clears away all results, including the server pregenerated html, and begins the search using the current logic.

