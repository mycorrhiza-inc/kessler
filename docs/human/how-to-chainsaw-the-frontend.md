
# Philosophy:

Havier Milei chainsaw meme, but for frontend code. 

Part of the MVP means stripping out all but the most essential pages namely `/`,  `/search`, `/orgs/<org-id>`, `/dockets/<docket-id>`, and maybe `/file/<file-id>`

Big part of the simplification is going to be cutting out essentially all state for the frontend app, for now the rule might be something like:

**URL's are the only engine of application state**

Build the initial components with as little state as possible



## Use server actions as the universal solution to the environment variable problem

So TLDR, the big problem with enviornment variables is essentially that the command process.env that gets enviornment variables uses a nextjs specific api.





# Client components


Root Page 

- Search Results (Server Data Rendered, No Client state)

- File Query (Server Static, Has Client State)

- 




All network requests and interactions with the backend occur on the server 

# Pages we want:

`/`

Page type: Static

Has a single client component, for search. When you click this
