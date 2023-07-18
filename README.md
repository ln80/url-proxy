## url-proxy (wip)

url-proxy a serverless proxy mainly used in combination with Cloudfront dists to adjust and cache third-party assets.

It takes advantage of the recent partial support of lambda `RESPONSE_STREAM` in Go runtime.

STATUS:

    Work in progress / Experimental

TODO:

- Protect the public lambda function URL using a custom header.
- Pass down `DeniedHost` list as stack parameter
- Add `ContentTypes` allow list as stack parameter
- Implement a custom **Round tripper** to cache assets in S3
