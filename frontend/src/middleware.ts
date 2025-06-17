import { NextResponse } from 'next/server';

export function middleware(request: Request) {
  const url = new URL(request.url);
  const origin = url.origin;
  const pathname = url.pathname;
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set('x-url', request.url);      // Full URL
  requestHeaders.set('x-origin', origin);        // Origin (scheme + host)
  requestHeaders.set('x-pathname', pathname);    // Path only

  return NextResponse.next({
    request: {
      headers: requestHeaders,
    }
  });
}
