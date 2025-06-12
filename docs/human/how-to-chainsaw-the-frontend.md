
# Philosophy:

Havier Milei chainsaw meme, but for frontend code. 

Part of the MVP means stripping out all but the most essential pages namely `/`,  `/search`, `/orgs/<org-id>`, `/dockets/<docket-id>`, and maybe `/file/<file-id>`

Big part of the simplification is going to be cutting out essentially all state for the frontend app, for now the rule might be something like:

**URL's are the only engine of application state**

Build the initial components with as little state as possible



## Use server actions as the universal solution to the environment variable problem

So TLDR, the big problem with environment variables is essentially that the command process.env that gets environment variables uses a nodejs specific api. So in order to get environment variables on the frontend you need to invent some kind of scheme to send them to the frontend, typically this is done in 2 ways:

1. During the compilation step replacing a call for `process.env.<VARNAME>` with the value of `<VARNAME>` at compile time. (Default NextJS behavior)

2. Serialize the environment variables into a json object, and send it to the frontend as an html element with a specific tag. Then when you need to look up a env variable you lookup the DOM element with that id, take the data, deserialize it and read the variables.

Process 1 is what next does by default, but it means that you cant set env variables at runtime, which makes deployment on any complex deployment system much much harder then it needs to be.

Option 2 is cursed and SHOULD NEVER BE USED UNDER ANY CIRCUMSTANCES

So how does this get solved? Well right now we have 2 problems 

1. How does the frontend know what backend endpoint to hit depending on its enviornment. (Typically either api.kessler.xyz, nightly-api.kessler.xyz, or localhost)

2. When running code, how does the code know what enviornment its running in, and hit PUBLIC_KESSLER_API, declared above, or http://backend-server:4041

I think the easiest way to solve this is to route all requests for data on the frontend through the backend:

Path Right Now:
Client -> Golang Backend && NextJS SSR -> Golang Backend


Proposed:
Client -> NextJS -> Golang Backend 

Using something like Server Functions
https://nextjs.org/docs/app/getting-started/updating-data

This solves both problems, with the downside of introducing some performance regressions on client side data requests. But this is only an extra JSON de/reserialization and doesnt include any more round trip requests.


# Component Structure


Root Page 

- Search Results (Server Data Rendered, No Client state)

- File Query (Server Static, Has Client State)

- 




All network requests and interactions with the backend occur on the server 

# Pages we want:

`/`

Page type: Static

Has a single client component, for search. When you click this
