Ok so currently I am kinda stuck on actually finishing the implementation for the frontend stuff. I have implemented everything with dummy data, but in order to actually get it working I would need to how how the data is shaped when it gets passed from fugu to the frontend. Mainly because most of the hard work exists in functions that 

- Fetch, validate and handle errors from the backend

- Convert that data into the intermediate format for the cards.

So would it be possible to figure out what the schema the backend is going to use when 

1a. What schema is the backend returning for a search request on a table with filings.
1b. What schema is being returned when querying a specific filing ID?


2a. What schema is the backend returning for a search request on a table with organizations.
2b. What schema is being returned when querying a specific organization ID?


2a. What schema is the backend returning for a search request on a table with conversations.
2b. What schema is being returned when querying a specific conversations ID?


For simplicity I think its best to pick a schema <T> for each, and use a list of <T> for A, and just a bare <T> for B.
