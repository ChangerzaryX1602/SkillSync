import { test, expect } from '@playwright/test';

const randomEmail = `adminuser_${Date.now()}@example.com`;
const password = 'Password123!';

test.describe('User Management Flows', () => {

// Helper to login before tests
test.beforeEach(async ({ page }) => {
    // Login with provided admin credentials
    await page.goto('/login');
    await page.fill('input[name="email"]', 'changnoi2547@gmail.com');
    await page.fill('input[name="password"]', 'cn16022547');
    await page.click('button[type="submit"]');
    
    // Wait for login to complete and redirect to dashboard
    await page.waitForURL('**/dashboard');
});

  test('User can navigate to users list and create a new user', async ({ page }) => {
    // 1. Navigate to Users List
    await page.goto('/users');
    // await expect(page).toHaveTitle(/Users/); // Title might vary, check h2
    await expect(page.locator('h1')).toContainText('Users'); // Based on +page.svelte (assuming it has h2 Users)
    
    // 2. Click Create User
    await page.click('a[href="/users/create"]');
    await expect(page).toHaveURL(/\/users\/create/);
    await expect(page.locator('h2')).toContainText('Create User'); // From users/create/+page.svelte

    // 3. Create User
    const newUserName = `newuser_${Date.now()}`;
    const newUserEmail = `new_${Date.now()}@example.com`;
    
    await page.fill('input[name="username"]', newUserName);
    await page.fill('input[name="email"]', newUserEmail);
    await page.fill('input[name="password"]', 'NewPass123!');
    
    await page.click('button[type="submit"]');

    // 4. Assert Success - should probably redirect back to users list
    await page.waitForURL('**/users');
    await expect(page.locator('tbody')).toContainText(newUserName);
  });
});
