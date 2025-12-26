import type { PageServerLoad } from './$types';
import { api } from '$lib/api';

export const load: PageServerLoad = async ({ cookies, url }) => {
    const token = cookies.get('token');
    const page = Number(url.searchParams.get('page')) || 1;
    const per_page = Number(url.searchParams.get('per_page')) || 10;
    const keyword = url.searchParams.get('keyword') || '';
    
    // Construct query string manually since api helper handles path
    const query = new URLSearchParams({
        page: page.toString(),
        per_page: per_page.toString(),
        keyword
    }).toString();

    try {
        const res = await api('GET', `/users?${query}`, null, fetch, token, cookies);
        if (res.success) {
            return {
                users: res.data,
                pagination: res.result.pagination,
                search: res.result.search
            };
        }
        return { users: [], pagination: null, error: 'Failed' };
     } catch (err) {
        console.error('Fetch users error:', err);
        return { users: [], pagination: null, error: 'Failed to load users' };
    }
};
