const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: { width: 1280, height: 900 } });

  // 监听控制台错误
  const errors = [];
  page.on('console', msg => {
    if (msg.type() === 'error') errors.push(msg.text());
  });
  page.on('pageerror', err => errors.push(err.message));

  console.log('1. 打开登录页...');
  await page.goto('http://localhost:3000/login', { waitUntil: 'networkidle', timeout: 10000 });
  await page.waitForTimeout(1000);

  const title = await page.title();
  console.log('   标题:', title);

  // 查找输入框
  const inputs = await page.$$('input');
  console.log('   输入框数量:', inputs.length);

  // 登录
  console.log('2. 测试登录...');
  await page.fill('input[placeholder*="用户名"]', 'admin');
  await page.fill('input[placeholder*="密码"]', '123456');
  await page.click('button:has-text("登录")');
  await page.waitForTimeout(2000);

  const url = page.url();
  console.log('   登录后URL:', url);

  // 截图
  await page.screenshot({ path: 'test/login_result.png' });
  console.log('3. 截图保存: test/login_result.png');

  // 检查会话列表
  const sessionsVisible = await page.$('text=会话列表');
  console.log('   会话列表:', sessionsVisible ? '可见' : '不可见');

  // 检查消息搜索按钮
  const searchBtn = await page.$('text=搜索');
  console.log('   搜索按钮:', searchBtn ? '可见' : '不可见');

  // 检查PPT助手
  const pptBtn = await page.$('text=PPT小助手');
  console.log('   PPT助手:', pptBtn ? '可见' : '不可见');

  // 检查群组列表
  await page.click('text=群组');
  await page.waitForTimeout(1000);
  const groupVisible = await page.$('text=一家人');
  console.log('   一家人群组:', groupVisible ? '可见' : '不可见');

  // 点击聊天
  await page.click('text=聊天');
  await page.waitForTimeout(500);

  await page.screenshot({ path: 'test/chat_result.png' });
  console.log('4. 聊天截图: test/chat_result.png');

  // 检查消息搜索页面
  await page.goto('http://localhost:3000/messages/search', { waitUntil: 'networkidle' });
  await page.waitForTimeout(500);
  const searchPage = await page.$('text=消息搜索');
  console.log('   消息搜索页:', searchPage ? '可见' : '不可见');
  await page.screenshot({ path: 'test/search_result.png' });
  console.log('5. 搜索页截图: test/search_result.png');

  console.log('\n=== 前端错误 ===');
  if (errors.length === 0) {
    console.log('   无错误');
  } else {
    errors.slice(0, 10).forEach(e => console.log('   ', e));
  }

  await browser.close();
  console.log('\n✅ 前端检查完成');
})();
