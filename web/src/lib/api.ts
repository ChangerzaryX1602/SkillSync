import { env } from '$env/dynamic/public';
import type { Cookies } from '@sveltejs/kit';

const BASE_URL = env.PUBLIC_API_URL || 'http://localhost:8080';

type Fetch = typeof fetch;

export async function api(
    method: string,
    path: string,
    data?: any,
    customFetch: Fetch = fetch,
    token?: string,
    // Optional cookies object for server-side refresh
    cookies?: Cookies 
) {
    const opts: RequestInit = {
        method,
        headers: {
            'Content-Type': 'application/json',
        },
    };

    if (data) {
        opts.body = JSON.stringify(data);
    }

    if (token) {
        opts.headers = {
            ...opts.headers,
            Authorization: `Bearer ${token}`,
        };
    }

    // Ensure path starts with /
    const normalizedPath = path.startsWith('/') ? path : `/${path}`;
    const url = `${BASE_URL}${normalizedPath}`;

    let res = await customFetch(url, opts);

    // Handle 401 Unauthorized - Refresh Token Flow
    if (res.status === 401 && cookies) {
        const refreshToken = cookies.get('refresh_token');
        if (refreshToken) {
            try {
                // Attempt to refresh
                const refreshRes = await customFetch(`${BASE_URL}/auth/refresh`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ refresh_token: refreshToken })
                });

                if (refreshRes.ok) {
                    const refreshData = await refreshRes.json();
                    if (refreshData.success && refreshData.access_token) {
                        // 1. Update Cookies
                        const commonOpts = { 
                            path: '/', 
                            httpOnly: true, 
                            secure: false, // process.env is not available in browser, and we are generic here. 
                            // In a real prod environment we might want to check $app/environment building
                            sameSite: 'lax' as const,
                            maxAge: 60 * 60 * 24 * 7 // 7 days
                        };
                        cookies.set('token', refreshData.access_token, commonOpts);
                        if (refreshData.refresh_token) {
                            cookies.set('refresh_token', refreshData.refresh_token, commonOpts);
                        }

                        // 2. Retry original request with new token
                        opts.headers = {
                            ...opts.headers,
                            Authorization: `Bearer ${refreshData.access_token}`
                        };
                        res = await customFetch(url, opts);
                    }
                } else {
                     // Refresh failed, maybe refresh token expired too
                     // Let the 401 pass through or clear cookies?
                     // Safer to let the app redirect to login naturally
                }
            } catch (err) {
                console.error('Auto-refresh failed:', err);
            }
        }
    }
    
    // Parse JSON safely
    let json = null;
    try {
        json = await res.json();
    } catch (e) {
        // Ignore empty or non-json response
    }

    if (!res.ok) {
        throw {
            status: res.status,
            message: json?.message || res.statusText,
            errors: json?.errors,
        };
    }

    return json;
}
