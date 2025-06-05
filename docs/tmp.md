So I am getting this issue when trying to get clerk to work on the frontend in 
/home/nicole/Documents/mycorrhizae/kessler/frontend/src

when I load it I get this error on the frontend
Error: Clerk: Handshake token verification failed: Unable to find a signing key in JWKS that matches the kid='ins_2xqHt2jXH8V7QWbmLFqA8uNpn7j' of the provided session token. Please make sure that the __session cookie or the HTTP authorization header contain a Clerk-generated session JWT. The following kid is available: ins_2hYVBIXYyVR123QuRowB7icPZE3 Go to your Dashboard and validate your secret and public keys are correct. Contact support@clerk.com if the issue persists. (reason=jwk-kid-mismatch, token-carrier=undefined).

dispite the url that we are logging in from 
https://ace-wallaby-84.clerk.accounts.dev/.well-known/jwks.json 

returns the following key that its complaining about not having:

```json
{"keys":[{"use":"sig","kty":"RSA","kid":"ins_2xqHt2jXH8V7QWbmLFqA8uNpn7j","alg":"RS256","n":"w2F9jfoHTpN5_hriW--ljT6xx4xH3Z8QGB9YcjRr2pgTFcQi4dLqiraz-v04j-QNW9nNThHUyD1jDbT8-sxcx2J0MWzbDm1zik0mjnJLaSLR8HuWQuqNJeS6o41kHyC_nZDhoSlRKfgKpNYKnkrqPPXvmA8D4vlAT0u2Tcvvz6cyn8g3v9v0RUkdaan19ZruDm60OcwuvxudETgCvboB9xb1nhFI0fHnz8LnzBowcscdToH9dRiSzwu_jI07tNmSRqYiCDUUS4eVivu8G5WiAPHV8E1-Efj3Qeno6JYOUQwxhc4vw6X3W_zEk2C_c8stG-1aT0bCIHd9ZsgcSPuC3w","e":"AQAB"}]}
```


I was also able to get some more debug information from the server:
```
frontend-1        | [clerk debug start: clerkMiddleware]
frontend-1        |   options, {
frontend-1        |     "publishableKey": "pk_test_YWNlLXdhbGxhYnktODQuY2xlcmsuYWNjb3VudHMuZGV2JA",
frontend-1        |     "secretKey": "sk_test_*********MC8",
frontend-1        |     "signInUrl": "",
frontend-1        |     "signUpUrl": "",
frontend-1        |     "debug": true
frontend-1        |   }
frontend-1        |   url, {
frontend-1        |     "url": "http://localhost/",
frontend-1        |     "method": "GET",
frontend-1        |     "headers": "{\"accept\":\"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\",\"accept-encoding\":\"gzip, deflate, br, zstd\",\"accept-language\":\"en-US,en;q=0.5\",\"cookie\":\"session=86888454-ac80-418c-ad7e-4dd43b81f7b7.-YLsUC4wyi3SzPXpuMR-8hcZKgw; ph_phc_3ESMmY9SgqEAGBB6sMGK5ayYHkeUuknH2vP6FmWH9RA_posthog=%7B%22distinct_id%22%3A%220195d86b-d861-7ff9-9d9c-dc8ef967522b%22%2C%22%24sesid%22%3A%5Bnull%2Cnull%2Cnull%5D%2C%22%24initial_person_info%22%3A%7B%22r%22%3A%22%24direct%22%2C%22u%22%3A%22http%3A%2F%2Flocalhost%3A3000%2Fconversations%2Ffba516e19534489aa4f8f064c8d9ec3f%22%7D%7D; session-512cdb6b=eyJjc3JmVG9rZW4iOiJmZDBjNGRmMjBhMjk2NTJkZjFiMTE5MDgwM2UyNGNhMjMzNTRiZjhiOGZlNzFkOTQ1ZTU4ODA1MTUxZGYzYTJiIn0=; session-512cdb6b.sig=sKpR1BZ8pW0Gkep4BsLLZZZeKTs; grafana_session=e3b3fe90ac319d2df3afe2242a3ddd68; grafana_session_expiry=1748222312\",\"host\":\"localhost\",\"priority\":\"u=0, i\",\"sec-fetch-dest\":\"document\",\"sec-fetch-mode\":\"navigate\",\"sec-fetch-site\":\"none\",\"sec-fetch-user\":\"?1\",\"upgrade-insecure-requests\":\"1\",\"user-agent\":\"Mozilla/5.0 (X11; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0\",\"x-forwarded-for\":\"172.18.0.1\",\"x-forwarded-host\":\"localhost\",\"x-forwarded-port\":\"80\",\"x-forwarded-proto\":\"http\",\"x-forwarded-server\":\"69d4efca9b96\",\"x-real-ip\":\"172.18.0.1\"}",
frontend-1        |     "clerkUrl": "http://localhost/",
frontend-1        |     "cookies": "{\"session\":\"86888454-ac80-418c-ad7e-4dd43b81f7b7.-YLsUC4wyi3SzPXpuMR-8hcZKgw\",\"ph_phc_3ESMmY9SgqEAGBB6sMGK5ayYHkeUuknH2vP6FmWH9RA_posthog\":\"{\\\"distinct_id\\\":\\\"0195d86b-d861-7ff9-9d9c-dc8ef967522b\\\",\\\"$sesid\\\":[null,null,null],\\\"$initial_person_info\\\":{\\\"r\\\":\\\"$direct\\\",\\\"u\\\":\\\"http://localhost:3000/conversations/fba516e19534489aa4f8f064c8d9ec3f\\\"}}\",\"session-512cdb6b\":\"eyJjc3JmVG9rZW4iOiJmZDBjNGRmMjBhMjk2NTJkZjFiMTE5MDgwM2UyNGNhMjMzNTRiZjhiOGZlNzFkOTQ1ZTU4ODA1MTUxZGYzYTJiIn0=\",\"session-512cdb6b.sig\":\"sKpR1BZ8pW0Gkep4BsLLZZZeKTs\",\"grafana_session\":\"e3b3fe90ac319d2df3afe2242a3ddd68\",\"grafana_session_expiry\":\"1748222312\"}"
frontend-1        |   }
frontend-1        |   requestState, {
frontend-1        |     "status": "handshake",
frontend-1        |     "headers": "{\"cache-control\":\"no-store\",\"location\":\"https://ace-wallaby-84.clerk.accounts.dev/v1/client/handshake?redirect_url=http%3A%2F%2Flocalhost%2F&__clerk_api_version=2025-04-10&suffixed_cookies=false&__clerk_hs_reason=dev-browser-missing\",\"set-cookie\":\"__clerk_redirect_count=1; SameSite=Lax; HttpOnly; Max-Age=3\",\"x-clerk-auth-reason\":\"dev-browser-missing\",\"x-clerk-auth-status\":\"handshake\"}",
frontend-1        |     "reason": "dev-browser-missing"
frontend-1        |   }
frontend-1        | [clerk debug end: clerkMiddleware] (@clerk/nextjs=6.20.2,next=14.2.29,timestamp=1749130735)
frontend-1        | [clerk debug start: clerkMiddleware]
frontend-1        |   options, {
frontend-1        |     "publishableKey": "pk_test_YWNlLXdhbGxhYnktODQuY2xlcmsuYWNjb3VudHMuZGV2JA",
frontend-1        |     "secretKey": "sk_test_*********MC8",
frontend-1        |     "signInUrl": "",
frontend-1        |     "signUpUrl": "",
frontend-1        |     "debug": true
frontend-1        |   }
frontend-1        |   url, {
frontend-1        |     "url": "http://localhost/?__clerk_handshake=eyJhbGciOiJSUzI1NiIsImNhdCI6ImNsX0I3ZDRQRDExMUFBQSIsImtpZCI6Imluc18yeHFIdDJqWEg4VjdRV2JtTEZxQTh1TnBuN2oiLCJ0eXAiOiJKV1QifQ.eyJoYW5kc2hha2UiOlsiX19jbGllbnRfdWF0PTsgUGF0aD0vOyBFeHBpcmVzPVRodSwgMDEgSmFuIDE5NzAgMDA6MDA6MDAgR01UOyBTYW1lU2l0ZT1MYXgiLCJfX2NsaWVudF91YXQ9MDsgUGF0aD0vOyBEb21haW49bG9jYWxob3N0OyBNYXgtQWdlPTMxNTM2MDAwMDsgU2FtZVNpdGU9TGF4IiwiX19zZXNzaW9uPTsgUGF0aD0vOyBFeHBpcmVzPVRodSwgMDEgSmFuIDE5NzAgMDA6MDA6MDAgR01UOyBTYW1lU2l0ZT1MYXgiLCJfX2NsZXJrX2RiX2p3dD1kdmJfMnk1aUZzVmF6bTJNZVN1bnAwamZTRWY3QlhaOyBQYXRoPS87IEV4cGlyZXM9RnJpLCAwNSBKdW4gMjAyNiAxMzozODo1NCBHTVQ7IFNhbWVTaXRlPUxheCJdfQ.M7JOsHceyoWNnJIhh6l8wWv6qNtkmpbrAtqPkOCIefNudhyzEfezcTqQ6kN3LQ6VNtjjI95VpiIlpbKAN7iyCFLnkxDej71BQt6tw3H4Rv8we9gbiB_Z-lc8KhT8sTeLwxYIXfswH9tJhCR_Omn0-qXghFwy_yelXVwXb-7YyxpUQCgr-4oB6vjVEdH1pDuCK-mzDu0sX5R2ET7R5__xDk_Z4WOFesw1ohaW1QDgen2AmiGH2H__16UuzqZIuV0568sYFFTx0utlNDCc41AUPzf7ZtF-PAu-TK2-cSgdbbhrnA9drEj_BZ8fSc33-JTjizBvI-fWH_GZKmbyPekJow",
frontend-1        |     "method": "GET",
frontend-1        |     "headers": "{\"accept\":\"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\",\"accept-encoding\":\"gzip, deflate, br, zstd\",\"accept-language\":\"en-US,en;q=0.5\",\"cookie\":\"session=86888454-ac80-418c-ad7e-4dd43b81f7b7.-YLsUC4wyi3SzPXpuMR-8hcZKgw; ph_phc_3ESMmY9SgqEAGBB6sMGK5ayYHkeUuknH2vP6FmWH9RA_posthog=%7B%22distinct_id%22%3A%220195d86b-d861-7ff9-9d9c-dc8ef967522b%22%2C%22%24sesid%22%3A%5Bnull%2Cnull%2Cnull%5D%2C%22%24initial_person_info%22%3A%7B%22r%22%3A%22%24direct%22%2C%22u%22%3A%22http%3A%2F%2Flocalhost%3A3000%2Fconversations%2Ffba516e19534489aa4f8f064c8d9ec3f%22%7D%7D; grafana_session=e3b3fe90ac319d2df3afe2242a3ddd68; grafana_session_expiry=1748222312; __clerk_redirect_count=1\",\"host\":\"localhost\",\"priority\":\"u=0, i\",\"sec-fetch-dest\":\"document\",\"sec-fetch-mode\":\"navigate\",\"sec-fetch-site\":\"none\",\"sec-fetch-user\":\"?1\",\"upgrade-insecure-requests\":\"1\",\"user-agent\":\"Mozilla/5.0 (X11; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0\",\"x-forwarded-for\":\"172.18.0.1\",\"x-forwarded-host\":\"localhost\",\"x-forwarded-port\":\"80\",\"x-forwarded-proto\":\"http\",\"x-forwarded-server\":\"69d4efca9b96\",\"x-real-ip\":\"172.18.0.1\"}",
frontend-1        |     "clerkUrl": "http://localhost/?__clerk_handshake=eyJhbGciOiJSUzI1NiIsImNhdCI6ImNsX0I3ZDRQRDExMUFBQSIsImtpZCI6Imluc18yeHFIdDJqWEg4VjdRV2JtTEZxQTh1TnBuN2oiLCJ0eXAiOiJKV1QifQ.eyJoYW5kc2hha2UiOlsiX19jbGllbnRfdWF0PTsgUGF0aD0vOyBFeHBpcmVzPVRodSwgMDEgSmFuIDE5NzAgMDA6MDA6MDAgR01UOyBTYW1lU2l0ZT1MYXgiLCJfX2NsaWVudF91YXQ9MDsgUGF0aD0vOyBEb21haW49bG9jYWxob3N0OyBNYXgtQWdlPTMxNTM2MDAwMDsgU2FtZVNpdGU9TGF4IiwiX19zZXNzaW9uPTsgUGF0aD0vOyBFeHBpcmVzPVRodSwgMDEgSmFuIDE5NzAgMDA6MDA6MDAgR01UOyBTYW1lU2l0ZT1MYXgiLCJfX2NsZXJrX2RiX2p3dD1kdmJfMnk1aUZzVmF6bTJNZVN1bnAwamZTRWY3QlhaOyBQYXRoPS87IEV4cGlyZXM9RnJpLCAwNSBKdW4gMjAyNiAxMzozODo1NCBHTVQ7IFNhbWVTaXRlPUxheCJdfQ.M7JOsHceyoWNnJIhh6l8wWv6qNtkmpbrAtqPkOCIefNudhyzEfezcTqQ6kN3LQ6VNtjjI95VpiIlpbKAN7iyCFLnkxDej71BQt6tw3H4Rv8we9gbiB_Z-lc8KhT8sTeLwxYIXfswH9tJhCR_Omn0-qXghFwy_yelXVwXb-7YyxpUQCgr-4oB6vjVEdH1pDuCK-mzDu0sX5R2ET7R5__xDk_Z4WOFesw1ohaW1QDgen2AmiGH2H__16UuzqZIuV0568sYFFTx0utlNDCc41AUPzf7ZtF-PAu-TK2-cSgdbbhrnA9drEj_BZ8fSc33-JTjizBvI-fWH_GZKmbyPekJow",
frontend-1        |     "cookies": "{\"session\":\"86888454-ac80-418c-ad7e-4dd43b81f7b7.-YLsUC4wyi3SzPXpuMR-8hcZKgw\",\"ph_phc_3ESMmY9SgqEAGBB6sMGK5ayYHkeUuknH2vP6FmWH9RA_posthog\":\"{\\\"distinct_id\\\":\\\"0195d86b-d861-7ff9-9d9c-dc8ef967522b\\\",\\\"$sesid\\\":[null,null,null],\\\"$initial_person_info\\\":{\\\"r\\\":\\\"$direct\\\",\\\"u\\\":\\\"http://localhost:3000/conversations/fba516e19534489aa4f8f064c8d9ec3f\\\"}}\",\"grafana_session\":\"e3b3fe90ac319d2df3afe2242a3ddd68\",\"grafana_session_expiry\":\"1748222312\",\"__clerk_redirect_count\":\"1\"}"
frontend-1        |   }
frontend-1        | [clerk debug end: clerkMiddleware] (@clerk/nextjs=6.20.2,next=14.2.29,timestamp=1749130735)
frontend-1        |  ⨯ Error: Clerk: Handshake token verification failed: Unable to find a signing key in JWKS that matches the kid='ins_2xqHt2jXH8V7QWbmLFqA8uNpn7j' of the provided session token. Please make sure that the __session cookie or the HTTP authorization header contain a Clerk-generated session JWT. The following kid is available: ins_2hYVBIXYyVR123QuRowB7icPZE3 Go to your Dashboard and validate your secret and public keys are correct. Contact support@clerk.com if the issue persists. (reason=jwk-kid-mismatch, token-carrier=undefined).
frontend-1        |     at HandshakeService.handleTokenVerificationErrorInDevelopment (webpack-internal:///(middleware)/./node_modules/@clerk/backend/dist/chunk-MQOIIRZU.mjs:3414:11)
frontend-1        |     at authenticateRequestWithTokenInCookie (webpack-internal:///(middleware)/./node_modules/@clerk/backend/dist/chunk-MQOIIRZU.mjs:3742:28)
frontend-1        |     at process.processTicksAndRejections (node:internal/process/task_queues:105:5)
frontend-1        |     at async eval (webpack-internal:///(middleware)/./node_modules/@clerk/nextjs/dist/esm/server/clerkMiddleware.js:86:28)
frontend-1        |     at async adapter (webpack-internal:///(middleware)/./node_modules/next/dist/esm/server/web/adapter.js:178:16)
frontend-1        |     at async (file:///app/node_modules/next/dist/server/web/sandbox/sandbox.js:97:22)
frontend-1        |     at async runWithTaggedErrors (file:///app/node_modules/next/dist/server/web/sandbox/sandbox.js:94:9)
frontend-1        |     at async DevServer.runMiddleware (file:///app/node_modules/next/dist/server/next-server.js:1068:24)
frontend-1        |     at async DevServer.runMiddleware (file:///app/node_modules/next/dist/server/dev/next-dev-server.js:268:28)
frontend-1        |     at async NextNodeServer.handleCatchallMiddlewareRequest (file:///app/node_modules/next/dist/server/next-server.js:322:26)
frontend-1        |  GET /?__clerk_handshake=eyJhbGciOiJSUzI1NiIsImNhdCI6ImNsX0I3ZDRQRDExMUFBQSIsImtpZCI6Imluc18yeHFIdDJqWEg4VjdRV2JtTEZxQTh1TnBuN2oiLCJ0eXAiOiJKV1QifQ.eyJoYW5kc2hha2UiOlsiX19jbGllbnRfdWF0PTsgUGF0aD0vOyBFeHBpcmVzPVRodSwgMDEgSmFuIDE5NzAgMDA6MDA6MDAgR01UOyBTYW1lU2l0ZT1MYXgiLCJfX2NsaWVudF91YXQ9MDsgUGF0aD0vOyBEb21haW49bG9jYWxob3N0OyBNYXgtQWdlPTMxNTM2MDAwMDsgU2FtZVNpdGU9TGF4IiwiX19zZXNzaW9uPTsgUGF0aD0vOyBFeHBpcmVzPVRodSwgMDEgSmFuIDE5NzAgMDA6MDA6MDAgR01UOyBTYW1lU2l0ZT1MYXgiLCJfX2NsZXJrX2RiX2p3dD1kdmJfMnk1aUZzVmF6bTJNZVN1bnAwamZTRWY3QlhaOyBQYXRoPS87IEV4cGlyZXM9RnJpLCAwNSBKdW4gMjAyNiAxMzozODo1NCBHTVQ7IFNhbWVTaXRlPUxheCJdfQ.M7JOsHceyoWNnJIhh6l8wWv6qNtkmpbrAtqPkOCIefNudhyzEfezcTqQ6kN3LQ6VNtjjI95VpiIlpbKAN7iyCFLnkxDej71BQt6tw3H4Rv8we9gbiB_Z-lc8KhT8sTeLwxYIXfswH9tJhCR_Omn0-qXghFwy_yelXVwXb-7YyxpUQCgr-4oB6vjVEdH1pDuCK-mzDu0sX5R2ET7R5__xDk_Z4WOFesw1ohaW1QDgen2AmiGH2H__16UuzqZIuV0568sYFFTx0utlNDCc41AUPzf7ZtF-PAu-TK2-cSgdbbhrnA9drEj_BZ8fSc33-JTjizBvI-fWH_GZKmbyPekJow 404 in 1ms
frontend-1        |  ⨯ Error [ERR_HTTP_HEADERS_SENT]: Cannot append headers after they are sent to the client
frontend-1        |     at ServerResponse.appendHeader (node:_http_outgoing:755:11)
frontend-1        | digest: "747541165"
frontend-1        | <w> [webpack.cache.PackFileCacheStrategy] Serializing big strings (123kiB) impacts deserialization performance (consider using Buffer instead and decode when needed
```

This seems weird what should I do to fix this?


There is a function here that wraps an infinite scroll component in react.


/home/nicole/Documents/mycorrhizae/kessler/frontend/src/components/InfiniteScroll/InfiniteScrollPlus.tsx


This imported component has given us a ton of problems. Could you make a new component at

/home/nicole/Documents/mycorrhizae/kessler/frontend/src/components/InfiniteScroll/InfiniteScroll.tsx

With an optimised event listener that will load more results when scrolling to the bottom.
