import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions = {
    default: async ({ cookies }) => {
        cookies.delete('token', { path: '/' });
        cookies.delete('refresh_token', { path: '/' });
        throw redirect(303, '/login');
    }
} satisfies Actions;
