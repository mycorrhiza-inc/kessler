


Could you make pages for the following routes 
/ 
/search 
/org/<org-id>
/convo/<convo-id>


where every page will parse these url paramaters:


    ?q= must populate the search box
        a new search should update this for the next page
    ?f:{filterkey}= must populate its respective filter
        these should be validated, when submitted, by the backend
    ?dataset= must have an indicator of which dataset it is in
        a new search should update this for the next page
        the filter struct already has an enabled field and can be greyed out if the previously picked filters do no apply to the new dataset


try to generalize the code so that the url parsing code will work across multiple endpoint
