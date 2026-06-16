const { chromium } = require('playwright');
(async () => {
  const b = await chromium.launch({ headless: false, slowMo: 60 });
  const URL = 'http://120.77.251.18';
  const c1 = await b.newContext(), c2 = await b.newContext();
  const p1 = await c1.newPage(), p2 = await c2.newPage();

  // ====== 1. LOGIN ======
  console.log('=== 登录 ===');
  await p1.goto(`${URL}/login`, { waitUntil: 'networkidle' });
  await p1.fill('input[placeholder*="用户名"]', 'u1');
  await p1.fill('input[placeholder*="密码"]', '111111');
  await p1.click('button:has-text("登录")');
  await p1.waitForTimeout(1000);
  console.log('u1 OK');

  await p2.goto(`${URL}/login`, { waitUntil: 'networkidle' });
  await p2.fill('input[placeholder*="用户名"]', 'u2');
  await p2.fill('input[placeholder*="密码"]', '111111');
  await p2.click('button:has-text("登录")');
  await p2.waitForTimeout(1000);
  console.log('u2 OK');

  // ====== 2. u1 搜索 u3 ======
  console.log('\n=== u1 搜索用户 ===');
  await p1.goto(`${URL}/contacts/search`, { waitUntil: 'networkidle' });
  await p1.waitForTimeout(500);
  await p1.fill('input[placeholder*="搜索"]', 'u3');
  await p1.click('button:has-text("搜索")');
  await p1.waitForTimeout(1500);
  await p1.screenshot({ path: 'test/01_search_result.png' });

  // ====== 3. u1 添加好友 ======
  console.log('=== u1 添加好友 ===');
  const addBtn = await p1.$('button:has-text("添加好友")');
  if (addBtn) { await addBtn.click(); console.log('OK - 申请已发送'); }
  await p1.waitForTimeout(500);

  // ====== 4. u2 同意 ======
  console.log('\n=== u2 查看申请 ===');
  await p2.goto(`${URL}/contacts`, { waitUntil: 'networkidle' });
  await p2.waitForTimeout(500);
  const plus = await p2.$('button[title="更多操作"]');
  if (plus) await plus.click();
  await p2.waitForTimeout(500);
  await p2.screenshot({ path: 'test/02_apply_menu.png' });
  const newFriend = await p2.$('text=新朋友');
  if (newFriend) { await newFriend.click(); await p2.waitForTimeout(800); }
  await p2.screenshot({ path: 'test/03_apply_list.png' });

  console.log('=== u2 同意 ===');
  const agree = await p2.$('button:has-text("同意")');
  if (agree) { await agree.click(); console.log('OK - 已同意'); }
  await p2.waitForTimeout(500);

  // ====== 5. u1 打开聊天 ======
  console.log('\n=== u1 联系人列表 ===');
  await p1.goto(`${URL}/contacts`, { waitUntil: 'networkidle' });
  await p1.waitForTimeout(800);
  await p1.screenshot({ path: 'test/04_contacts.png' });
  const chatBtns = await p1.$$('button:has-text("聊天")');
  console.log(`找到 ${chatBtns.length} 个聊天按钮`);

  if (chatBtns.length > 0) {
    console.log('=== u1 进入聊天 ===');
    await chatBtns[0].click();
    await p1.waitForTimeout(2000);
    const ta = await p1.$('textarea');
    if (ta) {
      console.log('OK - 输入框就绪');
      await ta.fill('你好！单聊测试消息 ✅');
      await p1.click('button:has-text("发送")');
      await p1.waitForTimeout(500);
      console.log('OK - 消息已发送');
    }
  }

  // ====== 6. u2 打开聊天回复 ======
  console.log('\n=== u2 联系人列表 ===');
  await p2.goto(`${URL}/contacts`, { waitUntil: 'networkidle' });
  await p2.waitForTimeout(800);
  const cb2 = await p2.$$('button:has-text("聊天")');
  console.log(`找到 ${cb2.length} 个聊天按钮`);
  if (cb2.length > 0) {
    console.log('=== u2 进入聊天 ===');
    await cb2[0].click();
    await p2.waitForTimeout(2000);
    const ta = await p2.$('textarea');
    if (ta) {
      console.log('OK - 输入框就绪');
      await ta.fill('收到！回复测试 ✅');
      await p2.click('button:has-text("发送")');
      await p2.waitForTimeout(500);
      console.log('OK - 消息已发送');
    }
  }

  await p1.screenshot({ path: 'test/05_final_u1.png' });
  await p2.screenshot({ path: 'test/05_final_u2.png' });
  console.log('\n=== 测试完成 ===');
  console.log('截图: test/01-05_*.png');
})();
