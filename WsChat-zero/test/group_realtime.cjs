const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: true });
  const server = 'http://120.77.251.18';
  const users = [
    { name: 'admin', user: 'admin', pass: '123456' },
    { name: 'u1', user: 'u1', pass: '111111' },
    { name: 'u2', user: 'u2', pass: '111111' },
  ];

  const pages = {};
  const contexts = {};

  // 1. All 3 users login
  console.log('=== 登录 ===');
  for (const u of users) {
    const ctx = await browser.newContext();
    const page = await ctx.newPage();
    page.on('pageerror', e => console.log(`  [${u.name} ERROR] ${e.message}`));

    await page.goto(`${server}/login`, { waitUntil: 'networkidle', timeout: 15000 });
    await page.fill('input[placeholder*="用户名"]', u.user);
    await page.fill('input[placeholder*="密码"]', u.pass);
    await page.click('button:has-text("登录")');
    await page.waitForTimeout(1500);
    const ok = page.url().includes('chat');
    console.log(`  ${u.name}: ${ok ? 'OK' : 'FAIL'}`);
    contexts[u.name] = ctx;
    pages[u.name] = page;
  }

  // 2. u1 opens 一家人 group
  console.log('\n=== 打开一家人群聊 ===');
  const u1 = pages['u1'];
  // Click group tab
  const groupTab = await u1.$('button:has-text("群聊")');
  if (groupTab) await groupTab.click();
  await u1.waitForTimeout(800);
  // Click 一家人
  const family = await u1.$('text=一家人');
  if (family) {
    await family.click();
    await u1.waitForTimeout(1000);
    console.log('  u1 opened 一家人');
  }

  // 3. u1 sends a message
  console.log('\n=== u1 发消息 ===');
  const input1 = await u1.$('textarea');
  if (input1) {
    await input1.fill('大家好！测试群聊实时广播 1/3');
    await u1.click('button:has-text("发送")');
    await u1.waitForTimeout(1500);
    console.log('  u1 sent message');
  }

  // 4. admin also opens 一家人 and sends
  const admin = pages['admin'];
  const gtab = await admin.$('button:has-text("群聊")');
  if (gtab) await gtab.click();
  await admin.waitForTimeout(800);
  const fam = await admin.$('text=一家人');
  if (fam) {
    await fam.click();
    await admin.waitForTimeout(1500);
    console.log('  admin opened 一家人');
  }
  const inputA = await admin.$('textarea');
  if (inputA) {
    await inputA.fill('admin 来了！测试群聊 2/3');
    await admin.click('button:has-text("发送")');
    await admin.waitForTimeout(1500);
    console.log('  admin sent message');
  }

  // 5. u2 opens and sends
  const u2 = pages['u2'];
  const gtab2 = await u2.$('button:has-text("群聊")');
  if (gtab2) await gtab2.click();
  await u2.waitForTimeout(800);
  const fam2 = await u2.$('text=一家人');
  if (fam2) {
    await fam2.click();
    await u2.waitForTimeout(1500);
    console.log('  u2 opened 一家人');
  }
  const input2 = await u2.$('textarea');
  if (input2) {
    await input2.fill('u2 也来了！群聊测试 3/3');
    await u2.click('button:has-text("发送")');
    await u2.waitForTimeout(1500);
    console.log('  u2 sent message');
  }

  // 6. Go back to u1 and check messages
  await u1.waitForTimeout(1000);

  // Screenshots
  for (const u of users) {
    await pages[u.name].screenshot({ path: `test/realtime_${u.name}.png` });
  }
  console.log('\nScreenshots: test/realtime_*.png');

  // 7. Check message text visible for u1
  const msgs = await u1.$$eval('div[style*="wordBreak"]', els => els.map(e => e.textContent?.slice(0,30)));
  console.log('\n=== u1 看到的最后几条消息 ===');
  msgs.slice(-5).forEach(m => console.log(`   ${m}`));

  for (const u of users) await contexts[u.name].close();
  await browser.close();
  console.log('\nDone');
})();
