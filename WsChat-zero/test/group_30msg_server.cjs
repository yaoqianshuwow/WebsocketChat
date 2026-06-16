const { chromium } = require('playwright');

(async () => {
  const b = await chromium.launch({ headless: true });
  const URL = 'http://120.77.251.18';
  const N = 30;

  const all = {};
  for (const [name, [user, pass]] of Object.entries({
    admin: ['admin','123456'], u1: ['u1','111111'], u2: ['u2','111111']
  })) {
    const ctx = await b.newContext();
    const p = await ctx.newPage();
    p.on('pageerror', e => console.log(`[${name}] ${e}`));
    await p.goto(`${URL}/login`, { waitUntil: 'networkidle', timeout: 15000 });
    await p.fill('input[placeholder*="用户名"]', user);
    await p.fill('input[placeholder*="密码"]', pass);
    await p.click('button:has-text("登录")');
    await p.waitForTimeout(2000);
    console.log(`${name}: ${p.url().includes('/chat') ? 'OK' : 'FAIL'}`);
    all[name] = { ctx, page: p };
  }

  // Enter group chat via Groups page
  console.log('\nOpening group...');
  for (const name of ['u1','admin','u2']) {
    const p = all[name].page;
    await p.goto(`${URL}/groups`, { waitUntil: 'networkidle', timeout: 10000 });
    await p.waitForTimeout(1000);
    const fam = await p.$('button:has-text("一家人")');
    if (fam) await fam.click();
    await p.waitForTimeout(500);
    const enter = await p.$('button:has-text("进入群聊")');
    if (enter) { await enter.click(); await p.waitForTimeout(1500); }
    console.log(`${name}: textarea=${await p.$('textarea') ? 'YES' : 'NO'}`);
  }

  // Send
  console.log(`\nSending ${N} each...`);
  const start = Date.now();
  await Promise.all(Object.entries(all).map(async ([name, { page }]) => {
    for (let i = 1; i <= N; i++) {
      const ta = await page.$('textarea');
      if (!ta) break;
      await ta.fill(`SERVER-${name}-${i}`);
      await page.click('button:has-text("发送")');
      if (i % 10 === 0) console.log(`  ${name}: ${i}`);
      await page.waitForTimeout(80);
    }
    console.log(`  ${name}: DONE`);
  }));
  console.log(`Done in ${((Date.now()-start)/1000).toFixed(1)}s`);

  await all['u1'].page.waitForTimeout(2000);
  for (const name of ['u1','admin','u2']) {
    await all[name].page.screenshot({ path: `test/server30_${name}.png` });
  }
  for (const v of Object.values(all)) await v.ctx.close();
  await b.close();
  console.log('DONE');
})();
