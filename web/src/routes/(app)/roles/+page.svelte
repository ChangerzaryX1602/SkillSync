<script lang="ts">
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import {
    Plus,
    Search,
    Edit2,
    ChevronLeft,
    ChevronRight,
  } from "lucide-svelte";

  let { data } = $props();
  let searchKeyword = $state(data.search?.keyword || "");

  $effect(() => {
    searchKeyword = data.search?.keyword || "";
  });

  function handleSearch(e: Event) {
    e.preventDefault();
    const url = new URL($page.url);
    url.searchParams.set("keyword", searchKeyword);
    url.searchParams.set("page", "1");
    goto(url.toString());
  }

  function changePage(newPage: number) {
    const url = new URL($page.url);
    url.searchParams.set("page", String(newPage));
    goto(url.toString());
  }
</script>

<div class="sm:flex sm:items-center">
  <div class="sm:flex-auto">
    <h1 class="text-xl font-semibold text-gray-900">Roles</h1>
    <p class="mt-2 text-sm text-gray-700">
      Manage user roles and their associated permissions.
    </p>
  </div>
  <div class="mt-4 sm:mt-0 sm:ml-16 sm:flex-none">
    <a
      href="/roles/create"
      class="inline-flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:w-auto"
    >
      <Plus class="mr-2 h-4 w-4" />
      Create Role
    </a>
  </div>
</div>

<div class="mt-8 flex flex-col">
  <div class="mb-4">
    <form onsubmit={handleSearch} class="flex gap-2">
      <div class="relative flex-grow max-w-md">
        <div
          class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3"
        >
          <Search class="h-5 w-5 text-gray-400" />
        </div>
        <input
          type="text"
          name="search"
          bind:value={searchKeyword}
          class="block w-full rounded-md border-gray-300 pl-10 focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm py-2 border shadow-sm"
          placeholder="Search roles..."
        />
      </div>
      <button
        type="submit"
        class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
      >
        Search
      </button>
    </form>
  </div>

  <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
    <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
      <div
        class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg"
      >
        <table class="min-w-full divide-y divide-gray-300">
          <thead class="bg-gray-50">
            <tr>
              <th
                scope="col"
                class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6"
                >ID</th
              >
              <th
                scope="col"
                class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                >Name</th
              >
              <th
                scope="col"
                class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                >Created At</th
              >
              <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                <span class="sr-only">Edit</span>
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 bg-white">
            {#if data.roles && data.roles.length > 0}
              {#each data.roles as role}
                <tr>
                  <td
                    class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6"
                    >{role.id}</td
                  >
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500"
                    >{role.name}</td
                  >
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500"
                    >{new Date(role.created_at).toLocaleDateString()}</td
                  >
                  <td
                    class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6"
                  >
                    <a
                      href="/roles/{role.id}"
                      class="text-indigo-600 hover:text-indigo-900 mr-4 inline-flex items-center"
                    >
                      <Edit2 class="h-4 w-4 mr-1" /> Edit
                    </a>
                  </td>
                </tr>
              {/each}
            {:else}
              <tr>
                <td
                  colspan="4"
                  class="px-6 py-4 text-center text-sm text-gray-500"
                  >No roles found.</td
                >
              </tr>
            {/if}
          </tbody>
        </table>

        <!-- Pagination -->
        {#if data.pagination}
          <div
            class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6"
          >
            <div
              class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between"
            >
              <div>
                <p class="text-sm text-gray-700">
                  Showing
                  <span class="font-medium"
                    >{(data.pagination.page - 1) * data.pagination.size +
                      1}</span
                  >
                  to
                  <span class="font-medium"
                    >{Math.min(
                      data.pagination.page * data.pagination.size,
                      data.pagination.total_item
                    )}</span
                  >
                  of
                  <span class="font-medium">{data.pagination.total_item}</span>
                  results
                </p>
              </div>
              <div>
                <nav
                  class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px"
                  aria-label="Pagination"
                >
                  <button
                    disabled={data.pagination.page <= 1}
                    onclick={() => changePage(data.pagination.page - 1)}
                    class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <span class="sr-only">Previous</span>
                    <ChevronLeft class="h-5 w-5" />
                  </button>
                  <span
                    class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700"
                  >
                    Page {data.pagination.page} of {data.pagination.total_page}
                  </span>
                  <button
                    disabled={data.pagination.page >=
                      data.pagination.total_page}
                    onclick={() => changePage(data.pagination.page + 1)}
                    class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <span class="sr-only">Next</span>
                    <ChevronRight class="h-5 w-5" />
                  </button>
                </nav>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
