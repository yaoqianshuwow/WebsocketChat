const { chromium } = require('playwright');
(async () => {
  const b = await chromium.launch({ headless: true });
  const URL = 'http://120.77.251.18';
  const all = {};

  console.log('Logging in 10 users...');
  for (let i = 1; i <= 10; i++) {
    const ctx = await b.newContext();
    const p = await ctx.newPage();
    await p.goto(`${URL}/login`, { waitUntil: 'networkidle', timeout: 15000 });
    await p.fill('input[placeholder*="用户名"]', `u${i}`);
    await p.fill('input[placeholder*="密码"]', '111111');
    await p.click('button:has-text("登录")');
    await p.waitForTimeout(1000);
    const ok = p.url().includes('chat');
    all[`u${i}`] = { ctx, page: p, ok };
    if (i % 3 === 0 || i === 10) console.log(`  logged in ${i}/10`);
  }

  console.log('Entering group...');
  let taCount = 0;
  for (const [name, { page }] of Object.entries(all)) {
    await page.goto(`${URL}/groups`, { waitUntil: 'networkidle', timeout: 10000 });
    await page.waitForTimeout(500);
    const fam = await page.$('button:has-text("一家人")');
    if (fam) await fam.click();
    await page.waitForTimeout(300);
    const enter = await page.$('button:has-text("进入群聊")');
    if (enter) { await enter.click(); await page.waitForTimeout(800); }
    if (await page.$('textarea')) taCount++;
  }
  console.log(`  textarea: ${taCount}/10`);

  const start = Date.now();
  const users = Object.keys(all);
  for (const name of users) {
    const p = all[name].page;
    for (let i = 1; i <= 3; i++) {
      await p.fill('textarea', `10p-${name}-${i}`);
      await p.click('button:has-text("发送")');
      await p.waitForTimeout(50);
    }
  }
  console.log(`Sent 30 msgs in ${((Date.now()-start)/1000).toFixed(1)}s`);

  for (const v of Object.values(all)) await v.ctx.close();
  await b.close();
})();
