import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { api } from '$lib/api';

export const actions = {
    default: async ({ request, cookies }) => {
        const token = cookies.get('token');
        const data = await request.formData();
        const { name } = Object.fromEntries(data);

        try {
            await api('POST', '/roles', { name }, fetch, token, cookies);
        } catch (err: any) {
            console.error('Create role error:', err);
            return fail(err.status || 500, { error: err.message || 'Failed to create role' });
        }

        throw redirect(303, '/roles');
    }
} satisfies Actions;
