import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { api } from '$lib/api';

export const actions = {
    register: async ({ request }) => {
        const data = await request.formData();
        const username = data.get('username') as string;
        const email = data.get('email') as string;
        const password = data.get('password') as string;

        if (!username || !email || !password) {
            return fail(400, { error: 'All fields are required' });
        }

        try {
            await api('POST', '/auth/register', { username, email, password });
        } catch (err: any) {
            console.error('Register error:', err);
            return fail(err.status || 500, { error: err.message || 'Registration failed' });
        }
        throw redirect(303, '/login');
    }
} satisfies Actions;
