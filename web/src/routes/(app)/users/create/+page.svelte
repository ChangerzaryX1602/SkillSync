<script lang="ts">
  import { enhance } from "$app/forms";
  import { ArrowLeft } from "lucide-svelte";
  let { form } = $props();
  let loading = $state(false);
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
      Create User
    </h2>
  </div>
</div>

<div class="mt-5 md:mt-0 md:col-span-2">
  {#if form?.error}
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
    use:enhance={() => {
      loading = true;
      return async ({ update }) => {
        loading = false;
        update();
      };
    }}
    class="space-y-6 bg-white shadow px-4 py-5 sm:rounded-lg sm:p-6"
  >
    <div class="grid grid-cols-6 gap-6">
      <div class="col-span-6 sm:col-span-3">
        <label for="username" class="block text-sm font-medium text-gray-700"
          >Username</label
        >
        <input
          type="text"
          name="username"
          id="username"
          required
          class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
        />
      </div>

      <div class="col-span-6 sm:col-span-3">
        <label for="email" class="block text-sm font-medium text-gray-700"
          >Email address</label
        >
        <input
          type="email"
          name="email"
          id="email"
          required
          class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
        />
      </div>

      <div class="col-span-6 sm:col-span-4">
        <label for="password" class="block text-sm font-medium text-gray-700"
          >Password</label
        >
        <input
          type="password"
          name="password"
          id="password"
          required
          class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
        />
      </div>
    </div>

    <div class="flex justify-end">
      <a
        href="/users"
        class="bg-white py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
      >
        Cancel
      </a>
      <button
        type="submit"
        disabled={loading}
        class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
      >
        {#if loading}
          Saving...
        {:else}
          Save
        {/if}
      </button>
    </div>
  </form>
</div>
