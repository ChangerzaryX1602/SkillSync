import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import { api } from '$lib/api';

interface Permission {
    id: number;
    name: string;
    group: string;
}

interface RolePermission {
    permission_id: number;
}

export const load: PageServerLoad = async ({ params, cookies }) => {
    const token = cookies.get('token');
    const roleId = params.id;

    try {
        const [roleRes, rolePermissionsRes, permissionsRes] = await Promise.all([
            api('GET', `/roles/${roleId}`, null, fetch, token, cookies),
            api('GET', `/role_permissions/role/${roleId}`, null, fetch, token, cookies),
            api('GET', `/permissions?per_page=1000`, null, fetch, token, cookies) // Fetch all permissions
        ]);

        if (roleRes.success) {
            let availablePermissions: Permission[] = [];
            if (permissionsRes.success) {
                // Ensure permissionsRes.data.items exists if it's paginated, or just data if list
                // Based on other endpoints it seems to return data.items for lists
                // But /permissions?per_page=1000 might return { items: [...] } 
                // Let's check api.ts result for permissions list
                availablePermissions = permissionsRes.data.items || permissionsRes.data || [];
            }
 
            const assignedPermissionIds = new Set<number>();
            if (rolePermissionsRes.success && rolePermissionsRes.data) {
                rolePermissionsRes.data.forEach((rp: RolePermission) => assignedPermissionIds.add(rp.permission_id));
            }

            return { 
                role: roleRes.data,
                assignedPermissionIds: Array.from(assignedPermissionIds),
                permissions: availablePermissions
            };
        }
        throw new Error('Role not found');
    } catch (err) {
        console.error('Fetch role details error:', err);
        throw redirect(303, '/roles');
    }
};

export const actions = {
    update: async ({ request, cookies, params }) => {
        const token = cookies.get('token');
        const id = params.id;
        const data = await request.formData();
        const { name } = Object.fromEntries(data);

        try {
            await api('PUT', `/roles/${id}`, { name }, fetch, token, cookies);
        } catch (err: any) {
             return fail(err.status || 500, { action: 'update', error: err.message || 'Failed to update role' });
        }

        return { success: true, action: 'update' };
    },
    updatePermissions: async ({ request, cookies, params }) => {
        const token = cookies.get('token');
        const roleId = params.id;
        const data = await request.formData();
        const permissionIds = data.getAll('permission_ids').map(id => Number(id));

        try {
            // 1. Delete all existing permissions for this role
            // We use DELETE /role_permissions/role/{roleId}
            await api('DELETE', `/role_permissions/role/${roleId}`, null, fetch, token, cookies);

            // 2. Add new permissions
            // Parallel requests might be too much if many, but for reasonable amount (10-50) it's ok.
            // A better way would be a bulk insert API but we don't have it.
            const promises = permissionIds.map(permId => 
                api('POST', '/role_permissions', { role_id: Number(roleId), permission_id: permId }, fetch, token, cookies)
            );
            
            await Promise.all(promises);

        } catch (err: any) {
             console.error('Update permissions error:', err);
             return fail(err.status || 500, { action: 'updatePermissions', error: err.message || 'Failed to update permissions' });
        }

        return { success: true, action: 'updatePermissions' };
    },
    delete: async ({ cookies, params }) => {
        const token = cookies.get('token');
        const id = params.id;

        try {
             await api('DELETE', `/roles/${id}`, null, fetch, token, cookies);
        } catch (err: any) {
             return fail(err.status || 500, { action: 'delete', error: err.message || 'Failed to delete role' });
        }
        throw redirect(303, '/roles');
    }
} satisfies Actions;
