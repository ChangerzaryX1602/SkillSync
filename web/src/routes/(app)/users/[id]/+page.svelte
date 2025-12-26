<script lang="ts">
  import { enhance } from "$app/forms";
  import { ArrowLeft, Trash2 } from "lucide-svelte";
  let { data, form } = $props();
  let loading = $state(false);
  let deleteLoading = $state(false);
</script>

<div class="md:flex md:items-center md:justify-between mb-6">
  <div class="flex-1 min-w-0">
    <nav class="flex" aria-label="Breadcrumb">
      <ol class="flex items-center space-x-4">
        <li>
          <div>
            <a href="/users" class="text-gray-400 hover:text-gray-500">
              <ArrowLeft class="flex-shrink-0 h-5 w-5" aria-hidden="true" />
              <span class="sr-only">Back</span>
            </a>
          </div>
        </li>
      </ol>
    </nav>
    <h2
      class="mt-2 text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate"
    >
      Edit User: {data.user.username}
    </h2>
  </div>
  <div class="mt-4 flex md:mt-0 md:ml-4">
    <form
      action="?/delete"
      method="POST"
      use:enhance={({ cancel }) => {
        const confirmed = confirm("Are you sure you want to delete this user?");
        if (!confirmed) return cancel();
        deleteLoading = true;
        return async ({ update }) => {
          deleteLoading = false;
          update();
        };
      }}
    >
      <button
        type="submit"
        disabled={deleteLoading}
        class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50"
      >
        <Trash2 class="-ml-1 mr-2 h-5 w-5" aria-hidden="true" />
        Delete
      </button>
    </form>
  </div>
</div>

<div class="mt-5 md:mt-0 grid grid-cols-1 gap-6 lg:grid-cols-2">
  <!-- User Details Form -->
  <div class="bg-white shadow px-4 py-5 sm:rounded-lg sm:p-6">
    <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">
      User Details
    </h3>
    {#if form?.error && form?.action === "update"}
      <div class="rounded-md bg-red-50 p-4 mb-4">
        <div class="flex">
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">
              {form.error}
            </h3>
          </div>
        </div>
      </div>
    {/if}

    <form
      method="POST"
      action="?/update"
      use:enhance={() => {
        loading = true;
        return async ({ update }) => {
          loading = false;
          update();
        };
      }}
      class="space-y-6"
    >
      <div>
        <label for="username" class="block text-sm font-medium text-gray-700"
          >Username</label
        >
        <input
          type="text"
          name="username"
          id="username"
          value={data.user.username}
          required
          class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
        />
      </div>

      <div>
        <label for="email" class="block text-sm font-medium text-gray-700"
          >Email address</label
        >
        <input
          type="email"
          name="email"
          id="email"
          value={data.user.email}
          required
          class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
        />
      </div>

      <div>
        <label for="password" class="block text-sm font-medium text-gray-700"
          >New Password (Optional)</label
        >
        <input
          type="password"
          name="password"
          id="password"
          class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
          placeholder="Leave blank to keep current"
        />
      </div>

      <div class="flex justify-end">
        <button
          type="submit"
          disabled={loading}
          class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
        >
          {#if loading}
            Saving...
          {:else}
            Save Changes
          {/if}
        </button>
      </div>
    </form>
  </div>

  <!-- Roles Form -->
  <div class="bg-white shadow px-4 py-5 sm:rounded-lg sm:p-6">
    <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">Roles</h3>
    <p class="text-sm text-gray-500 mb-4">
      Select roles to assign to this user.
    </p>

    <form
      method="POST"
      action="?/updateRoles"
      use:enhance={() => {
        loading = true; // reusing loading or create separate state
        return async ({ update }) => {
          loading = false;
          update();
        };
      }}
    >
      <div class="space-y-4 max-h-[600px] overflow-y-auto pr-2">
        {#if data.roles && data.roles.length > 0}
          {#each data.roles as role}
            <div class="relative flex items-start">
              <div class="flex items-center h-5">
                <input
                  id="role-{role.id}"
                  name="role_ids"
                  type="checkbox"
                  value={role.id}
                  checked={data.assignedRoleIds.includes(role.id)}
                  class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded"
                />
              </div>
              <div class="ml-3 text-sm">
                <label for="role-{role.id}" class="font-medium text-gray-700"
                  >{role.name}</label
                >
              </div>
            </div>
          {/each}
        {:else}
          <p class="text-sm text-gray-500">No roles available.</p>
        {/if}
      </div>

      <div class="flex justify-end mt-6 pt-4 border-t border-gray-200">
        <button
          type="submit"
          class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          Update Roles
        </button>
      </div>
    </form>
  </div>
</div>
