import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { api } from '$lib/api';

export const actions = {
    login: async ({ request, cookies }) => {
        const data = await request.formData();
        const email = data.get('email') as string;
        const password = data.get('password') as string;

        if (!email || !password) {
            return fail(400, { error: 'Email and password are required' });
        }
        console.log('[DEBUG] Login Action: Attempting login for', email);

        try {
            const res = await api('POST', '/auth/login', { email, password });
            console.log('[DEBUG] Login Action: API Response', JSON.stringify(res));
            
            if (res.success) {
                cookies.set('token', res.access_token, {
                    path: '/',
                    httpOnly: true,
                    sameSite: 'lax',
                    secure: false,
                    maxAge: 60 * 60 * 24 // 1 day
                });
                
                // Also store refresh token if needed, usually safer to just keep access token short lived and refresh it, but for simplicity here:
                cookies.set('refresh_token', res.refresh_token, {
                    path: '/',
                    httpOnly: true,
                    sameSite: 'lax',
                    secure: false,
                    maxAge: 60 * 60 * 24 * 7 // 7 days
                });
            } else {
                return fail(401, { error: 'Invalid credentials' });
            }
        } catch (err: any) {
            console.error('Login error:', err);
            return fail(err.status || 500, { error: err.message || 'Login failed' });
        }

        throw redirect(303, '/dashboard');
    }
} satisfies Actions;
