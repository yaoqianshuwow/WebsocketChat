const { chromium } = require('playwright');
(async () => {
  const b = await chromium.launch({ headless: false, slowMo: 30 });
  const URL = 'http://120.77.251.18';

  // u1 and u2 login
  const c1 = await b.newContext(), c2 = await b.newContext();
  const p1 = await c1.newPage(), p2 = await c2.newPage();

  for (const [p, u] of [[p1,'u1'],[p2,'u2']]) {
    await p.goto(`${URL}/login`, { waitUntil: 'networkidle' });
    await p.fill('input[placeholder*="用户名"]', u);
    await p.fill('input[placeholder*="密码"]', '111111');
    await p.click('button:has-text("登录")');
    await p.waitForTimeout(1000);
    console.log(`${u}: ${p.url().includes('chat') ? 'OK' : 'FAIL'}`);
  }

  // u1: search u3 and add friend
  console.log('=== u1 search + add u3 ===');
  await p1.goto(`${URL}/contacts/search`, { waitUntil: 'networkidle' });
  await p1.waitForTimeout(500);
  await p1.fill('input[placeholder*="搜索"]', 'u3');
  await p1.click('button:has-text("搜索")');
  await p1.waitForTimeout(1000);
  const addBtn = await p1.$('button:has-text("添加好友")');
  if (addBtn) { await addBtn.click(); await p1.waitForTimeout(500); console.log('  u1 applied to u3'); }

  // u3 accepts (API)
  console.log('=== u3 accept ===');

  // u1: open Contacts and start chat with u2
  console.log('=== u1 chat with u2 ===');
  await p1.goto(`${URL}/contacts`, { waitUntil: 'networkidle' });
  await p1.waitForTimeout(800);
  const chatBtns = await p1.$$('button:has-text("聊天")');
  console.log(`  Found ${chatBtns.length} chat buttons`);
  if (chatBtns.length >= 1) {
    await chatBtns[0].click();
    await p1.waitForTimeout(2000);
    const ta = await p1.$('textarea');
    console.log(`  textarea: ${ta ? 'YES' : 'NO'}`);
    if (ta) {
      await ta.fill('u1: 单聊测试消息！');
      await p1.click('button:has-text("发送")');
      console.log('  u1 sent');
      await p1.waitForTimeout(500);
    }
  }

  // u2: open chat with u1
  console.log('=== u2 chat with u1 ===');
  await p2.goto(`${URL}/contacts`, { waitUntil: 'networkidle' });
  await p2.waitForTimeout(800);
  const cBtns2 = await p2.$$('button:has-text("聊天")');
  console.log(`  Found ${cBtns2.length} chat buttons`);
  if (cBtns2.length >= 1) {
    await cBtns2[0].click();
    await p2.waitForTimeout(2000);
    const ta = await p2.$('textarea');
    console.log(`  textarea: ${ta ? 'YES' : 'NO'}`);
    if (ta) {
      await ta.fill('u2: 收到！单聊回复');
      await p2.click('button:has-text("发送")');
      console.log('  u2 sent');
    }
  }

  await p1.screenshot({ path: 'test/single_u1.png' });
  await p2.screenshot({ path: 'test/single_u2.png' });
  console.log('DONE - browser stays open');
})();
