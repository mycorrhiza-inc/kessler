# What we need k8s to do:

So we currently have an app that consists of 3 elements that communicate with each other:

- a nextjs app that serves the frontend and does misc server side rendering.

- a backend server api written in go that handles communication with an external postgres and text search database.

- an ingest api server written in go to validate and transform incoming data and perform long running 


these container images will be accessible with names like 

mycorrhiza/frontend:a357e01e28dd38369fae277e7182b96f0bb77b85
mycorrhiza/backend-server:a357e01e28dd38369fae277e7182b96f0bb77b85
mycorrhiza/backend-ingest:a357e01e28dd38369fae277e7182b96f0bb77b85


where the key a357e01e28dd38369fae277e7182b96f0bb77b85 represents the commit hash of the monorepo that built all the images. 

currently we have a simple script that manually updates a docker compose file to download the new images, but we would also need a solution to get this to work with kubernetes, so we could update the nightly deploy regularly with new images. And once tested update production.


All of these services should be accessible through different endpoints like:


<server-ip>:3000 for the frontend 

<server-ip>:4041 for the backend api

<server-ip>:4042 for the ingest


However there should be multiple services deployed for each of these, for now there are only two enviornments:
- "prod"
- "nightly"
(maybe more added later)

And the front of the kubernetes cluster there should be a router that picks what deployment each request gets routed too.

It should make this decision based on a cookie included along with the request. So something like 
```
Cookie: target_deployment_override=prod
```

if the cookie is not set default to routing to the prod target.





