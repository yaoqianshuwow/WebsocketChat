const { chromium } = require('playwright');
(async () => {
  const b = await chromium.launch({ headless: true });
  const URL = 'http://120.77.251.18';

  // Login u1 and u2
  const ctx1 = await b.newContext(), ctx2 = await b.newContext();
  const p1 = await ctx1.newPage(), p2 = await ctx2.newPage();

  console.log('=== Login ===');
  for (const [p, u, pw] of [[p1,'u1','111111'],[p2,'u2','111111']]) {
    await p.goto(`${URL}/login`, { waitUntil: 'networkidle', timeout: 15000 });
    await p.fill('input[placeholder*="用户名"]', u);
    await p.fill('input[placeholder*="密码"]', pw);
    await p.click('button:has-text("登录")');
    await p.waitForTimeout(1000);
    console.log(`  ${u}: ${p.url().includes('chat') ? 'OK' : 'FAIL'}`);
  }

  // u1 opens Contacts and adds u2 as friend
  console.log('=== Add friend ===');
  await p1.goto(`${URL}/contacts`, { waitUntil: 'networkidle' });
  await p1.waitForTimeout(500);
  // Click + button to open menu
  const plusBtn = await p1.$('button[title="更多操作"]');
  if (plusBtn) await plusBtn.click();
  await p1.waitForTimeout(300);
  // Click search user
  const searchBtn = await p1.$('text=查找用户');
  if (searchBtn) await searchBtn.click();
  await p1.waitForTimeout(500);
  // Check if already friends
  const alreadyFriend = await p1.$('text=暂无联系人');
  if (alreadyFriend) {
    console.log('  No friends yet, searching for u2...');
  }

  // u1 opens single chat with u2
  console.log('=== Open single chat ===');
  await p1.goto(`${URL}/contacts/search`, { waitUntil: 'networkidle' });
  await p1.waitForTimeout(500);
  const srchInput = await p1.$('input[placeholder*="搜索"]');
  if (srchInput) {
    await srchInput.fill('u2');
    await p1.click('button:has-text("搜索")');
    await p1.waitForTimeout(1000);
    const addBtn = await p1.$('button:has-text("添加好友")');
    if (addBtn) {
      await addBtn.click();
      await p1.waitForTimeout(500);
      console.log('  u1 applied to add u2');
    }
  }

  // u2 accepts
  await p2.goto(`${URL}/contacts`, { waitUntil: 'networkidle' });
  await p2.waitForTimeout(500);
  const plus2 = await p2.$('button[title="更多操作"]');
  if (plus2) await plus2.click();
  await p2.waitForTimeout(300);
  const newFriend = await p2.$('text=新朋友');
  if (newFriend) await newFriend.click();
  await p2.waitForTimeout(500);
  const agreeBtn = await p2.$('button:has-text("同意")');
  if (agreeBtn) { await agreeBtn.click(); await p2.waitForTimeout(500); console.log('  u2 accepted'); }

  // u1 starts chat with u2
  console.log('=== Chat ===');
  await p1.goto(`${URL}/contacts`, { waitUntil: 'networkidle' });
  await p1.waitForTimeout(500);
  // Click chat button for u2 if available
  const chatBtn = await p1.$('button:has-text("聊天")');
  if (chatBtn) {
    await chatBtn.click();
    await p1.waitForTimeout(1000);
  }

  const hasTA = await p1.$('textarea');
  console.log(`  u1 textarea: ${hasTA ? 'YES' : 'NO'}`);
  if (hasTA) {
    for (let i = 1; i <= 3; i++) {
      await hasTA.fill(`u1->u2: 你好${i}`);
      await p1.click('button:has-text("发送")');
      await p1.waitForTimeout(100);
    }
    console.log('  u1 sent 3 messages');
  }

  // u2 responds
  await p2.goto(`${URL}/contacts`, { waitUntil: 'networkidle' });
  await p2.waitForTimeout(500);
  const chatBtn2 = await p2.$('button:has-text("聊天")');
  if (chatBtn2) { await chatBtn2.click(); await p2.waitForTimeout(1000); }
  const ta2 = await p2.$('textarea');
  if (ta2) {
    for (let i = 1; i <= 3; i++) {
      await ta2.fill(`u2->u1: 回复${i}`);
      await p2.click('button:has-text("发送")');
      await p2.waitForTimeout(100);
    }
    console.log('  u2 sent 3 messages');
  }

  await p1.screenshot({ path: 'test/single_u1.png' });
  await p2.screenshot({ path: 'test/single_u2.png' });

  await ctx1.close(); await ctx2.close(); await b.close();
  console.log('DONE');
})();
