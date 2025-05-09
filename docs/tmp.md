I have a kinda hackish method for getting the server runtime enviornment variables on the client, that are in both of these files.

kessler/frontend/src/lib/env_variables_hydration_script.tsx
kessler/frontend/src/lib/env_variables_root_script.tsx

it does work, but I am needing to refactor this a bit to better support functions that need access to these variables on both the client and the server and I need good interfaces to support that, and thought it would be a good time to improve the architecture.

Throw any ideas for improvement in kessler/docs/frontend_env_refactor.md
