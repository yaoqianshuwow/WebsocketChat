const { chromium } = require('playwright');

async function loginUser(page, name, user, pass) {
  console.log(`[${name}] Login...`);
  await page.goto('http://120.77.251.18/login', { waitUntil: 'networkidle', timeout: 15000 });
  await page.fill('input[placeholder*="用户名"]', user);
  await page.fill('input[placeholder*="密码"]', pass);
  await page.click('button:has-text("登录")');
  await page.waitForTimeout(2000);
  const url = page.url();
  console.log(`[${name}] Logged in: ${url.includes('chat') ? 'OK' : 'FAIL'} - ${url}`);
  return url.includes('chat');
}

(async () => {
  const browser = await chromium.launch({ headless: false }); // visible browser

  // 1. Admin login
  const adminCtx = await browser.newContext();
  const admin = await adminCtx.newPage();
  await loginUser(admin, 'admin', 'admin', '123456');

  // 2. u1 login
  const u1Ctx = await browser.newContext();
  const u1 = await u1Ctx.newPage();
  await loginUser(u1, 'u1', 'u1', '111111');

  // 3. u2 login
  const u2Ctx = await browser.newContext();
  const u2 = await u2Ctx.newPage();
  await loginUser(u2, 'u2', 'u2', '111111');

  // 4. Check page elements for admin
  console.log('\n--- Checking admin page ---');
  const checks = ['搜消息', 'PPT小助手', '会话列表', '群聊'];
  for (const c of checks) {
    const el = await admin.$(`text=${c}`);
    console.log(`  ${c}: ${el ? 'OK' : 'MISSING'}`);
  }

  // 5. Click on group tab
  console.log('\n--- Group chat test ---');
  const groupTab = await u1.$('button:has-text("群聊")');
  if (groupTab) {
    await groupTab.click();
    await u1.waitForTimeout(500);
  }

  // Look for 一家人
  const family = await u1.$('text=一家人');
  if (family) {
    console.log('  一家人 found, opening...');
    await family.click();
    await u1.waitForTimeout(1000);

    // Send a group message
    const input = await u1.$('textarea');
    if (input) {
      await input.fill('大家好，测试群聊广播！');
      await u1.click('button:has-text("发送")');
      console.log('  u1 sent: 大家好，测试群聊广播！');
      await u1.waitForTimeout(2000);
    }
  }

  // Screenshots
  await admin.screenshot({ path: 'test/group_admin.png' });
  await u1.screenshot({ path: 'test/group_u1.png' });
  await u2.screenshot({ path: 'test/group_u2.png' });
  console.log('\nScreenshots saved: test/group_*.png');

  console.log('\nDone. Browser stays open for manual inspection.');
  // await browser.close();
})();
