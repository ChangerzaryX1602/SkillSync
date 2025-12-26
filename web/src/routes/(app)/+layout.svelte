<script lang="ts">
  import { page } from "$app/stores";
  import {
    LayoutDashboard,
    Users,
    Shield,
    Key,
    LogOut,
    Menu,
    X,
  } from "lucide-svelte";

  let { children, data } = $props();
  let sidebarOpen = $state(false);

  const navigation = [
    { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
    { name: "Users", href: "/users", icon: Users },
    { name: "Roles", href: "/roles", icon: Shield },
    { name: "Permissions", href: "/permissions", icon: Key },
  ];

  function isCurrent(path: string) {
    return $page.url.pathname.startsWith(path);
  }
</script>

<div class="min-h-screen bg-gray-100">
  <!-- Mobile sidebar -->
  {#if sidebarOpen}
    <div
      class="fixed inset-0 flex z-40 lg:hidden"
      role="dialog"
      aria-modal="true"
    >
      <div
        class="fixed inset-0 bg-gray-600 bg-opacity-75"
        aria-hidden="true"
        onclick={() => (sidebarOpen = false)}
      ></div>
      <div class="relative flex-1 flex flex-col max-w-xs w-full bg-white">
        <div class="absolute top-0 right-0 -mr-12 pt-2">
          <button
            class="ml-1 flex items-center justify-center h-10 w-10 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
            onclick={() => (sidebarOpen = false)}
          >
            <span class="sr-only">Close sidebar</span>
            <X class="h-6 w-6 text-white" />
          </button>
        </div>
        <div class="flex-1 h-0 pt-5 pb-4 overflow-y-auto">
          <div class="flex-shrink-0 flex items-center px-4">
            <span class="text-xl font-bold text-indigo-600">SkillSync</span>
          </div>
          <nav class="mt-5 px-2 space-y-1">
            {#each navigation as item}
              <a
                href={item.href}
                class="{isCurrent(item.href)
                  ? 'bg-indigo-50 text-indigo-600'
                  : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'} group flex items-center px-2 py-2 text-base font-medium rounded-md"
              >
                <item.icon
                  class="{isCurrent(item.href)
                    ? 'text-indigo-600'
                    : 'text-gray-400 group-hover:text-gray-500'} mr-4 h-6 w-6"
                />
                {item.name}
              </a>
            {/each}
          </nav>
        </div>
        <div class="flex-shrink-0 flex border-t border-gray-200 p-4">
          <div class="flex items-center">
            <div>
              <p class="text-base font-medium text-gray-700">
                {data.user?.username}
              </p>
              <form action="/logout" method="POST">
                <button
                  type="submit"
                  class="text-sm font-medium text-gray-500 hover:text-gray-700"
                  >Logout</button
                >
              </form>
            </div>
          </div>
        </div>
      </div>
      <div class="flex-shrink-0 w-14"></div>
    </div>
  {/if}

  <!-- Desktop sidebar -->
  <div
    class="hidden lg:flex lg:flex-col lg:w-64 lg:fixed lg:inset-y-0 lg:border-r lg:border-gray-200 lg:pt-5 lg:pb-4 lg:bg-white"
  >
    <div class="flex items-center flex-shrink-0 px-6">
      <span class="text-2xl font-bold text-indigo-600">SkillSync</span>
    </div>
    <div class="mt-6 h-0 flex-1 flex flex-col overflow-y-auto">
      <!-- User Account Dropdown -->

      <nav class="px-3 mt-6">
        <div class="space-y-1">
          {#each navigation as item}
            <a
              href={item.href}
              class="{isCurrent(item.href)
                ? 'bg-indigo-50 text-indigo-600'
                : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'} group flex items-center px-2 py-2 text-sm font-medium rounded-md transition-colors"
            >
              <item.icon
                class="{isCurrent(item.href)
                  ? 'text-indigo-600'
                  : 'text-gray-400 group-hover:text-gray-500'} mr-3 h-5 w-5"
              />
              {item.name}
            </a>
          {/each}
        </div>
      </nav>
    </div>
    <div class="flex-shrink-0 flex border-t border-gray-200 p-4">
      <div class="flex-shrink-0 w-full group block">
        <div class="flex items-center">
          <div class="ml-3">
            <p class="text-sm font-medium text-gray-700">
              {data.user?.username}
            </p>
            <form action="/logout" method="POST" class="mt-1">
              <button
                type="submit"
                class="flex items-center text-xs font-medium text-gray-500 hover:text-gray-700 group-hover:text-gray-900"
              >
                <LogOut class="mr-1.5 h-3 w-3" />
                Sign out
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="lg:pl-64 flex flex-col flex-1">
    <div
      class="sticky top-0 z-10 flex-shrink-0 flex h-16 bg-white shadow lg:hidden"
    >
      <button
        type="button"
        class="px-4 border-r border-gray-200 text-gray-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 lg:hidden"
        onclick={() => (sidebarOpen = true)}
      >
        <span class="sr-only">Open sidebar</span>
        <Menu class="h-6 w-6" />
      </button>
      <div class="flex-1 px-4 flex justify-between">
        <div class="flex-1 flex"></div>
      </div>
    </div>

    <main class="flex-1">
      <div class="py-6">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
          {@render children()}
        </div>
      </div>
    </main>
  </div>
</div>
