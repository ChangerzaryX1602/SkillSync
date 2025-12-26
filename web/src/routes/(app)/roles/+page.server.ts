import type { PageServerLoad } from './$types';
import { api } from '$lib/api';

export const load: PageServerLoad = async ({ cookies, url }) => {
    const token = cookies.get('token');
    const page = Number(url.searchParams.get('page')) || 1;
    const per_page = Number(url.searchParams.get('per_page')) || 10;
    const keyword = url.searchParams.get('keyword') || '';

    const query = new URLSearchParams({
        page: String(page),
        per_page: String(per_page),
        keyword: keyword,
    }).toString();

    try {
        const res = await api('GET', `/roles?${query}`, null, fetch, token, cookies);
        if (res.success) {
            return {
                roles: res.data.items,
                pagination: res.result.pagination,
                search: res.result.search
            };
        }
        return { roles: [], pagination: null, error: 'Failed' };
    } catch (err) {
        console.error('Fetch roles error:', err);
        return { roles: [], pagination: null, error: 'Failed to load roles' };
    }
};
