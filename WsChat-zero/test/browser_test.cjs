const { chromium } = require('playwright');

async function testSite(name, url) {
  console.log(`\n========== ${name} : ${url} ==========`);
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: { width: 1280, height: 900 } });
  const errors = [];
  page.on('pageerror', err => errors.push(err.message));

  try {
    // 1. Login page
    console.log('1. Login page...');
    await page.goto(url + '/login', { waitUntil: 'networkidle', timeout: 15000 });
    await page.waitForTimeout(500);
    console.log('   Title:', await page.title());

    // 2. Login
    console.log('2. Login...');
    await page.fill('input[placeholder*="用户名"]', 'admin');
    await page.fill('input[placeholder*="密码"]', '123456');
    await page.click('button:has-text("登录")');
    await page.waitForTimeout(2000);
    console.log('   URL:', page.url());
    const loggedIn = page.url().includes('chat');
    console.log('   Login:', loggedIn ? 'OK' : 'FAIL');

    // 3. Check elements
    console.log('3. Page elements...');
    await page.waitForTimeout(500);
    const checks = [
      ['会话列表', 'text=会话列表'],
      ['搜消息按钮', 'text=搜消息'],
      ['PPT助手', 'text=PPT小助手'],
    ];
    for (const [label, sel] of checks) {
      const el = await page.$(sel);
      console.log(`   ${label}: ${el ? 'OK' : 'MISSING'}`);
    }

    // 4. Click group tab
    console.log('4. Groups...');
    const groupTab = await page.$('button:has-text("群聊")');
    if (groupTab) await groupTab.click();
    await page.waitForTimeout(500);
    const group = await page.$('text=一家人');
    console.log('   一家人 group:', group ? 'OK' : 'MISSING');

    // 5. Search page
    console.log('5. Message search...');
    await page.goto(url + '/messages/search', { waitUntil: 'networkidle' });
    await page.waitForTimeout(500);
    const searchEl = await page.$('text=消息搜索');
    console.log('   Search page:', searchEl ? 'OK' : 'MISSING');

    // Screenshot
    await page.screenshot({ path: `test/browser_${name.replace(/[:\/]/g,'_')}.png` });
    console.log('   Screenshot saved');

  } catch(e) {
    console.log('   ERROR:', e.message);
  }

  console.log('Errors:', errors.length > 0 ? errors.slice(0,3).join(' | ') : 'none');
  await browser.close();
}

(async () => {
  await testSite('LOCAL', 'http://localhost:3000');
  await testSite('SERVER', 'http://120.77.251.18');
  console.log('\n===== DONE =====');
})();
