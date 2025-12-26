import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { api } from '$lib/api';

export const actions = {
    default: async ({ request, cookies }) => {
        const token = cookies.get('token');
        const data = await request.formData();
        const username = data.get('username') as string;
        const email = data.get('email') as string;
        const password = data.get('password') as string;

        if (!username || !email || !password) {
            return fail(400, { error: 'All fields are required' });
        }

        try {
            await api('POST', '/users', { username, email, password }, fetch, token, cookies);
        } catch (err: any) {
            console.error('Create user error:', err);
            return fail(err.status || 500, { error: err.message || 'Failed to create user' });
        }

        throw redirect(303, '/users');
    }
} satisfies Actions;
