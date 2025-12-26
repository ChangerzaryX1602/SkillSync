import type { PageServerLoad } from './$types';
import { api } from '$lib/api';

export const load: PageServerLoad = async ({ cookies, url }) => {
    const token = cookies.get('token');
    const page = Number(url.searchParams.get('page')) || 1;
    const per_page = Number(url.searchParams.get('per_page')) || 100;
    
    // Construct query string manually since api helper handles path
    const query = new URLSearchParams({
        page: String(page),
        per_page: String(per_page),
    }).toString();

    try {
        const res = await api('GET', `/permissions?${query}`, null, fetch, token, cookies);
        if (res.success) {
            return {
                permissions: res.data.items,
            };
        }
        return { permissions: [], pagination: null, error: 'Failed' };
     } catch (err) {
        console.error('Fetch permissions error:', err);
        return { permissions: [], pagination: null, error: 'Failed to load permissions' };
    }
};
