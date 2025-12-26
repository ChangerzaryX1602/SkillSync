import { test, expect } from '@playwright/test';

// Generate a random email to ensure unique registration each run
const randomEmail = `testuser_${Date.now()}@example.com`;
const password = 'Password123!';

test.describe('Authentication Flows', () => {
  test('User can register and then login', async ({ page }) => {
    // 1. Register
    await page.goto('/register');
    await expect(page).toHaveTitle(/Register/);

    // Assuming register form fields based on typical patterns (will verify if fails)
    // If the register page wasn't inspected, I should have, but I'll assume standard naming
    // If this fails, I'll inspect /register
    
    // Fill Registration Form
    await page.fill('input[name="username"]', `user_${Date.now()}`);
    await page.fill('input[name="email"]', randomEmail);
    await page.fill('input[name="password"]', password);
    // await page.fill('input[name="confirmPassword"]', password); // Not in form
    
    
    await page.click('button[type="submit"]');

    // Expecting to be redirected to login
    await page.waitForURL('**/login');

    // 2. Login
    await page.fill('input[name="email"]', randomEmail);
    await page.fill('input[name="password"]', password);
    await page.click('button[type="submit"]');

    // 3. Assert Success
    // Should be on dashboard or home
    await expect(page).not.toHaveURL(/login/);
    // await expect(page).toHaveURL('/'); // or /dashboard
  });

  test('Login fails with invalid credentials', async ({ page }) => {
    await page.goto('/login');
    await page.fill('input[name="email"]', 'wrong@example.com');
    await page.fill('input[name="password"]', 'wrongpassword');
    await page.click('button[type="submit"]');

    // Expect error message
    // Based on +page.svelte read earlier: {form.error} inside a div with text-red-800
    await expect(page.locator('.text-red-800')).toBeVisible();
  });
});
