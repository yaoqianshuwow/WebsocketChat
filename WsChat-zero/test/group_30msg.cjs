const { chromium } = require('playwright');

(async () => {
  const b = await chromium.launch({ headless: false, slowMo: 20 }); // visible browser
  const URL = 'http://localhost:3000';
  const N = 30;

  const all = {};
  for (const [name, [user, pass]] of Object.entries({
    admin: ['admin','123456'], u1: ['u1','111111'], u2: ['u2','111111']
  })) {
    const ctx = await b.newContext();
    const p = await ctx.newPage();
    await p.goto(`${URL}/login`, { waitUntil: 'networkidle', timeout: 15000 });
    await p.fill('input[placeholder*="用户名"]', user);
    await p.fill('input[placeholder*="密码"]', pass);
    await p.click('button:has-text("登录")');
    await p.waitForTimeout(1500);
    console.log(`${name}: ${p.url().includes('chat') ? 'LOGIN OK' : 'LOGIN FAIL'}`);
    all[name] = { ctx, page: p };
  }

  // Open groups page and enter 一家人 chat
  console.log('\nEntering group chat...');
  for (const name of ['u1','admin','u2']) {
    const p = all[name].page;
    await p.goto(`${URL}/groups`, { waitUntil: 'networkidle', timeout: 10000 });
    await p.waitForTimeout(1000);

    // Click 一家人 group
    const fam = await p.$('button:has-text("一家人")');
    if (fam) await fam.click();
    await p.waitForTimeout(800);

    // Click "进入群聊" button
    const enter = await p.$('button:has-text("进入群聊")');
    if (enter) {
      await enter.click();
      await p.waitForTimeout(1500);
    }

    const hasInput = await p.$('textarea');
    console.log(`${name}: textarea=${hasInput ? 'YES' : 'NO'} url=${p.url()}`);
    all[name].page = p;
  }

  // Send 30 messages each in parallel
  console.log(`\nSending ${N} messages each...`);
  const start = Date.now();

  await Promise.all(Object.entries(all).map(async ([name, { page }]) => {
    for (let i = 1; i <= N; i++) {
      const ta = await page.$('textarea');
      if (!ta) { console.log(`${name}: lost textarea at ${i}`); break; }
      await ta.fill(`${name}-${i}`);
      await page.click('button:has-text("发送")');
      if (i % 10 === 0) console.log(`  ${name}: ${i}`);
      await page.waitForTimeout(50);
    }
    console.log(`  ${name}: DONE`);
  }));

  console.log(`\nAll done in ${((Date.now()-start)/1000).toFixed(1)}s`);

  // Wait for render
  await all['u1'].page.waitForTimeout(2000);

  // Check what u1 sees
  const texts = await all['u1'].page.$$eval('[style*="wordBreak"]', els =>
    els.map(e => e.textContent?.trim()?.slice(0,50))
  );
  const groupMsgs = texts.filter(t => t && t.includes('-'));
  console.log(`\nu1 message list: ${groupMsgs.length} visible (expecting 90)`);
  groupMsgs.slice(-5).forEach(t => console.log(`  ${t}`));

  for (const name of ['u1','admin','u2']) {
    await all[name].page.screenshot({ path: `test/30msg_${name}.png` });
  }

  for (const v of Object.values(all)) await v.ctx.close();
  await b.close();
  console.log('\nDONE');
})();
