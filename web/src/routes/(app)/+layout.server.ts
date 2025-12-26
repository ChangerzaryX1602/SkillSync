import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { api } from '$lib/api';

export const load: LayoutServerLoad = async ({ cookies, fetch, url }) => {
    const token = cookies.get('token');
    console.log(`[DEBUG] Layout Load: URL ${url.pathname}, Token present: ${!!token}, Token len: ${token?.length}`);

    if (token) {
        console.log('[DEBUG] Layout Load: Token found, verifying...');
        try {
            const res = await api('GET', '/users/me', null, fetch, token, cookies);
            console.log('[DEBUG] Layout Load: API Response success:', res.success);
            if (res.success) {
                return {
                    user: res.data
                };
            } else {
                console.log('[DEBUG] Layout Load: API success false');
                // Token invalid or API returned an error
                cookies.delete('token', { path: '/' });
                cookies.delete('refresh_token', { path: '/' });
                throw redirect(303, '/login');
            }
        } catch (err) {
            console.error('[DEBUG] Layout Load: Auth check error:', err);
            cookies.delete('token', { path: '/' });
            cookies.delete('refresh_token', { path: '/' });
            throw redirect(303, '/login');
        }
    } else {
        // Only redirect to login if we are in a protected route? 
        // Actually +layout.server.ts runs for ALL routes in (app).
        // If (app) is the protected group, then yes.
        // /login is in (auth), so it does NOT use this layout.
        // So redirection is correct for (app).
        console.log('[DEBUG] Layout Load: No token, redirecting from', url.pathname);
        throw redirect(303, '/login');
    }
};
