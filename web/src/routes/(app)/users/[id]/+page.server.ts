import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import { api } from '$lib/api';

interface Role {
    id: number;
    name: string;
}

interface UserRole {
    role_id: number;
}

export const load: PageServerLoad = async ({ params, cookies }) => {
    const token = cookies.get('token');
    const id = params.id;

    try {
        const [userRes, userRolesRes, allRolesRes] = await Promise.all([
             api('GET', `/users/${id}`, null, fetch, token, cookies),
             api('GET', `/user_roles/user/${id}`, null, fetch, token, cookies),
             api('GET', `/roles?per_page=100`, null, fetch, token, cookies)
        ]);

        if (userRes.success) {
            const assignedRoleIds = new Set<number>();
            if (userRolesRes.success && userRolesRes.data) {
                userRolesRes.data.forEach((ur: UserRole) => assignedRoleIds.add(ur.role_id));
            }

            return { 
                user: userRes.data,
                assignedRoleIds: Array.from(assignedRoleIds),
                roles: allRolesRes.success ? allRolesRes.data : []
            };
        }
        throw new Error('User not found');
    } catch (err) {
        console.error('Fetch user details error:', err);
        throw redirect(303, '/users');
    }
};

export const actions = {
    update: async ({ request, cookies, params }) => {
        const token = cookies.get('token');
        const id = params.id;
        const data = await request.formData();
        const username = data.get('username') as string;
        const email = data.get('email') as string;
        const password = data.get('password') as string;

         const payload: any = { username, email };
         if (password) {
             payload.password = password;
         }

        try {
            await api('PATCH', `/users/${id}`, payload, fetch, token, cookies);
        } catch (err: any) {
             console.error('Update user error:', err);
             return fail(err.status || 500, { action: 'update', error: err.message || 'Failed to update user' });
        }

        return { success: true, action: 'update' };
    },
    updateRoles: async ({ request, cookies, params }) => {
        const token = cookies.get('token');
        const userId = params.id;
        const data = await request.formData();
        const roleIds = data.getAll('role_ids').map(id => Number(id));

        try {
            // 1. Delete all existing roles for this user
            await api('DELETE', `/user_roles/user/${userId}`, null, fetch, token, cookies);

            // 2. Add new roles
            const promises = roleIds.map(roleId => 
                api('POST', '/user_roles', { user_id: Number(userId), role_id: roleId }, fetch, token, cookies)
            );
            
            await Promise.all(promises);

        } catch (err: any) {
             console.error('Update roles error:', err);
             return fail(err.status || 500, { action: 'updateRoles', error: err.message || 'Failed to update roles' });
        }

        return { success: true, action: 'updateRoles' };
    },
    delete: async ({ cookies, params }) => {
        const token = cookies.get('token');
        const id = params.id;

        try {
             await api('DELETE', `/users/${id}`, null, fetch, token, cookies);
        } catch (err: any) {
             console.error('Delete user error:', err);
             return fail(err.status || 500, { action: 'delete', error: err.message || 'Failed to delete user' });
        }
        throw redirect(303, '/users');
    }
} satisfies Actions;
