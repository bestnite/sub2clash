function setInputReadOnly(input, readonly) {
  if (readonly) {
    input.readOnly = true;
    input.style.cursor = 'not-allowed';
  } else {
    input.readOnly = false;
    input.style.cursor = 'auto';
  }
}

function clearExistingValues() {
  // 清除简单输入框和复选框的值
  document.getElementById("endpoint").value = "clash";
  document.getElementById("sub").value = "";
  document.getElementById("proxy").value = "";
  document.getElementById("refresh").checked = false;
  document.getElementById("autoTest").checked = false;
  document.getElementById("lazy").checked = false;
  document.getElementById("igcg").checked = false;
  document.getElementById("useUDP").checked = false;
  document.getElementById("template").value = "";
  document.getElementById("sort").value = "nameasc";
  document.getElementById("remove").value = "";
  document.getElementById("apiLink").value = "";
  document.getElementById("apiShortLink").value = "";

  // 恢复短链ID和密码输入框状态
  const customIdInput = document.getElementById("customId");
  const passwordInput = document.getElementById("password");
  const generateButton = document.querySelector('button[onclick="generateShortLink()"]');

  customIdInput.value = "";
  setInputReadOnly(customIdInput, false);

  passwordInput.value = "";
  setInputReadOnly(passwordInput, false);

  // 恢复生成短链按钮状态
  generateButton.disabled = false;
  generateButton.classList.remove('btn-secondary');
  generateButton.classList.add('btn-primary');

  document.getElementById("nodeList").checked = false;

  // 清除由 createRuleProvider, createReplace, 和 createRule 创建的所有额外输入组
  clearInputGroup("ruleProviderGroup");
  clearInputGroup("replaceGroup");
  clearInputGroup("ruleGroup");
}

function generateURI() {
  const config = {};

  config.clashType = parseInt(document.getElementById("endpoint").value);

  let subLines = document
    .getElementById("sub")
    .value.split("\n")
    .filter((line) => line.trim() !== "");
  if (subLines.length > 0) {
    config.subscriptions = subLines;
  }

  let proxyLines = document
    .getElementById("proxy")
    .value.split("\n")
    .filter((line) => line.trim() !== "");
  if (proxyLines.length > 0) {
    config.proxies = proxyLines;
  }

  if (
    (config.subscriptions === undefined || config.subscriptions.length === 0) &&
    (config.proxies === undefined || config.proxies.length === 0)
  ) {
    return "";
  }

  config.userAgent = document.getElementById("user-agent").value;

  config.refresh = document.getElementById("refresh").checked;
  config.autoTest = document.getElementById("autoTest").checked;
  config.lazy = document.getElementById("lazy").checked;
  config.nodeList = document.getElementById("nodeList").checked;
  config.ignoreCountryGroup = document.getElementById("igcg").checked;
  config.useUDP = document.getElementById("useUDP").checked;

  const template = document.getElementById("template").value;
  if (template.trim() !== "") {
    config.template = template;
  }

  const ruleProvidersElements = document.getElementsByName("ruleProvider");
  if (ruleProvidersElements.length > 0) {
    const ruleProviders = [];
    for (let i = 0; i < ruleProvidersElements.length / 5; i++) {
      let baseIndex = i * 5;
      let behavior = ruleProvidersElements[baseIndex].value;
      let url = ruleProvidersElements[baseIndex + 1].value;
      let group = ruleProvidersElements[baseIndex + 2].value;
      let prepend = ruleProvidersElements[baseIndex + 3].value;
      let name = ruleProvidersElements[baseIndex + 4].value;
      if (
        behavior.trim() === "" ||
        url.trim() === "" ||
        group.trim() === "" ||
        prepend.trim() === "" ||
        name.trim() === ""
      ) {
        return "";
      }
      ruleProviders.push({
        behavior: behavior,
        url: url,
        group: group,
        prepend: prepend.toLowerCase() === "true",
        name: name,
      });
    }
    if (ruleProviders.length > 0) {
      config.ruleProviders = ruleProviders;
    }
  }

  const rulesElements = document.getElementsByName("rule");
  if (rulesElements.length > 0) {
    const rules = [];
    for (let i = 0; i < rulesElements.length / 2; i++) {
      if (rulesElements[i * 2].value.trim() !== "") {
        let rule = rulesElements[i * 2].value;
        let prepend = rulesElements[i * 2 + 1].value;
        if (rule.trim() === "" || prepend.trim() === "") {
          return "";
        }
        rules.push({
          rule: rule,
          prepend: prepend.toLowerCase() === "true",
        });
      }
    }
    if (rules.length > 0) {
      config.rules = rules;
    }
  }

  config.sort = document.getElementById("sort").value;

  const remove = document.getElementById("remove").value;
  if (remove.trim() !== "") {
    config.remove = remove;
  }

  const replacesElements = document.getElementsByName("replace");
  if (replacesElements.length > 0) {
    const replace = {};
    for (let i = 0; i < replacesElements.length / 2; i++) {
      let replaceStr = replacesElements[i * 2].value;
      let replaceTo = replacesElements[i * 2 + 1].value;
      if (replaceStr.trim() === "") {
        return "";
      }
      replace[replaceStr] = replaceTo;
    }
    if (Object.keys(replace).length > 0) {
      config.replace = replace;
    }
  }

  const jsonString = JSON.stringify(config);
  // 解决 btoa 中文报错，使用 TextEncoder 进行 UTF-8 编码再 base64
  function base64EncodeUnicode(str) {
    const bytes = new TextEncoder().encode(str);
    let binary = '';
    bytes.forEach((b) => binary += String.fromCharCode(b));
    return btoa(binary);
  }
  const encoded = base64EncodeUnicode(jsonString);
  const urlSafeBase64 = encoded
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");

  return `convert/${urlSafeBase64}`;
}

// 将输入框中的 URL 解析为参数
async function parseInputURL() {
  // 获取输入框中的 URL
  const inputURL = document.getElementById("urlInput").value;
  // 清除现有的输入框值
  clearExistingValues();
  if (!inputURL) {
    alert("请输入有效的链接！");
    return;
  }

  let url;
  try {
    url = new URL(inputURL);
  } catch (_) {
    alert("无效的链接！");
    return;
  }
  if (url.pathname.includes("/s/")) {
    let hash = url.pathname.substring(url.pathname.lastIndexOf("/s/") + 3);
    let q = new URLSearchParams();
    let password = url.searchParams.get("password");
    if (password === null) {
      alert("仅可解析加密短链");
      return;
    }
    q.append("hash", hash);
    q.append("password", password);
    try {
      const response = await axios.get("./short?" + q.toString());
      url = new URL(response.data, window.location.href);

      // 回显配置链接
      const apiLinkInput = document.querySelector("#apiLink");
      apiLinkInput.value = url.href;
      setInputReadOnly(apiLinkInput, true);

      // 回显短链相关信息
      const apiShortLinkInput = document.querySelector("#apiShortLink");
      apiShortLinkInput.value = inputURL;
      setInputReadOnly(apiShortLinkInput, true);

      // 设置短链ID和密码，并设置为只读
      const customIdInput = document.querySelector("#customId");
      const passwordInput = document.querySelector("#password");
      const generateButton = document.querySelector('button[onclick="generateShortLink()"]');

      customIdInput.value = hash;
      setInputReadOnly(customIdInput, true);

      passwordInput.value = password;
      setInputReadOnly(passwordInput, true);

      // 禁用生成短链按钮
      generateButton.disabled = true;
      generateButton.classList.add('btn-secondary');
      generateButton.classList.remove('btn-primary');
    } catch (error) {
      console.log(error);
      alert("获取短链失败，请检查密码！");
    }
  }
  const pathSections = url.pathname.split("/");
  const convertIndex = pathSections.findIndex((s) => s === "convert");

  if (convertIndex === -1 || convertIndex + 1 >= pathSections.length) {
    alert("无效的配置链接，请确认链接为新版格式。");
    return;
  }
  const base64Config = pathSections[convertIndex + 1];
  let config;
  try {
    const regularBase64 = base64Config.replace(/-/g, "+").replace(/_/g, "/");
    const decodedStr = atob(regularBase64);
    config = JSON.parse(decodeURIComponent(escape(decodedStr)));
  } catch (e) {
    alert("解析配置失败！");
    console.error(e);
    return;
  }

  document.getElementById("endpoint").value = config.clashType || "1";

  if (config.subscriptions) {
    document.getElementById("sub").value = config.subscriptions.join("\n");
  }

  if (config.proxies) {
    document.getElementById("proxy").value = config.proxies.join("\n");
  }

  if (config.refresh) {
    document.getElementById("refresh").checked = config.refresh;
  }

  if (config.autoTest) {
    document.getElementById("autoTest").checked = config.autoTest;
  }

  if (config.lazy) {
    document.getElementById("lazy").checked = config.lazy;
  }

  if (config.template) {
    document.getElementById("template").value = config.template;
  }

  if (config.sort) {
    document.getElementById("sort").value = config.sort;
  }

  if (config.remove) {
    document.getElementById("remove").value = config.remove;
  }

  if (config.userAgent) {
    document.getElementById("user-agent").value = config.userAgent;
  }

  if (config.ignoreCountryGroup) {
    document.getElementById("igcg").checked = config.ignoreCountryGroup;
  }

  if (config.replace) {
    const replaceGroup = document.getElementById("replaceGroup");
    for (const original in config.replace) {
      const div = createReplace();
      div.children[0].value = original;
      div.children[1].value = config.replace[original];
      replaceGroup.appendChild(div);
    }
  }

  if (config.ruleProviders) {
    const ruleProviderGroup = document.getElementById("ruleProviderGroup");
    for (const p of config.ruleProviders) {
      const div = createRuleProvider();
      div.children[0].value = p.behavior;
      div.children[1].value = p.url;
      div.children[2].value = p.group;
      div.children[3].value = p.prepend;
      div.children[4].value = p.name;
      ruleProviderGroup.appendChild(div);
    }
  }

  if (config.rules) {
    const ruleGroup = document.getElementById("ruleGroup");
    for (const r of config.rules) {
      const div = createRule();
      div.children[0].value = r.rule;
      div.children[1].value = r.prepend;
      ruleGroup.appendChild(div);
    }
  }

  if (config.nodeList) {
    document.getElementById("nodeList").checked = config.nodeList;
  }

  if (config.useUDP) {
    document.getElementById("useUDP").checked = config.useUDP;
  }
}

function clearInputGroup(groupId) {
  // 清空第二个之后的child
  const group = document.getElementById(groupId);
  while (group.children.length > 2) {
    group.removeChild(group.lastChild);
  }
}

async function copyToClipboard(elem, e) {
  const apiLinkInput = document.querySelector(`#${elem}`).value;
  try {
    await navigator.clipboard.writeText(apiLinkInput);
    let text = e.textContent;
    e.addEventListener("mouseout", function () {
      e.textContent = text;
    });
    e.textContent = "复制成功";
  } catch (err) {
    console.error("复制到剪贴板失败:", err);
  }
}

function createRuleProvider() {
  const div = document.createElement("div");
  div.classList.add("input-group", "mb-2");
  div.innerHTML = `
            <input type="text" class="form-control" name="ruleProvider" placeholder="Behavior">
            <input type="text" class="form-control" name="ruleProvider" placeholder="Url">
            <input type="text" class="form-control" name="ruleProvider" placeholder="Group">
            <input type="text" class="form-control" name="ruleProvider" placeholder="Prepend">
            <input type="text" class="form-control" name="ruleProvider" placeholder="Name">
            <button type="button" class="btn btn-danger" onclick="removeElement(this)">删除</button>
        `;
  return div;
}

function createReplace() {
  const div = document.createElement("div");
  div.classList.add("input-group", "mb-2");
  div.innerHTML = `
            <input type="text" class="form-control" name="replace" placeholder="原字符串（正则表达式）">
            <input type="text" class="form-control" name="replace" placeholder="替换为（可为空）">
            <button type="button" class="btn btn-danger" onclick="removeElement(this)">删除</button>
        `;
  return div;
}

function createRule() {
  const div = document.createElement("div");
  div.classList.add("input-group", "mb-2");
  div.innerHTML = `
            <input type="text" class="form-control" name="rule" placeholder="Rule">
            <input type="text" class="form-control" name="rule" placeholder="Prepend">
            <button type="button" class="btn btn-danger" onclick="removeElement(this)">删除</button>
        `;
  return div;
}

function listenInput() {
  let selectElements = document.querySelectorAll("select");
  let inputElements = document.querySelectorAll("input");
  let textAreaElements = document.querySelectorAll("textarea");
  inputElements.forEach(function (element) {
    element.addEventListener("input", function () {
      generateURL();
    });
  });
  textAreaElements.forEach(function (element) {
    element.addEventListener("input", function () {
      generateURL();
    });
  });
  selectElements.forEach(function (element) {
    element.addEventListener("change", function () {
      generateURL();
    });
  });
}

function addRuleProvider() {
  const div = createRuleProvider();
  document.getElementById("ruleProviderGroup").appendChild(div);
  listenInput();
}

function addRule() {
  const div = createRule();
  document.getElementById("ruleGroup").appendChild(div);
  listenInput();
}

function addReplace() {
  const div = createReplace();
  document.getElementById("replaceGroup").appendChild(div);
  listenInput();
}

function removeElement(button) {
  button.parentElement.remove();
}

function generateURL() {
  const apiLink = document.getElementById("apiLink");
  let uri = generateURI();
  if (uri === "") {
    return;
  }
  apiLink.value = `${window.location.origin}${window.location.pathname}${uri}`;
  setInputReadOnly(apiLink, true);
}

function generateShortLink() {
  const apiShortLink = document.getElementById("apiShortLink");
  const password = document.getElementById("password");
  const customId = document.getElementById("customId");
  let uri = generateURI();
  if (uri === "") {
    return;
  }

  axios
    .post(
      "./short",
      {
        url: uri,
        password: password.value.trim(),
        customId: customId.value.trim()
      },
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    )
    .then((response) => {
      // 设置返回的短链ID和密码
      customId.value = response.data.hash;
      password.value = response.data.password;
      // 生成完整的短链接
      const shortLink = `${window.location.origin}${window.location.pathname}s/${response.data.hash}?password=${response.data.password}`;
      apiShortLink.value = shortLink;
    })
    .catch((error) => {
      console.log(error);
      if (error.response && error.response.data) {
        alert(error.response.data);
      } else {
        alert("生成短链失败，请重试！");
      }
    });
}

function updateShortLink() {
  const password = document.getElementById("password");
  const apiShortLink = document.getElementById("apiShortLink");
  let hash = apiShortLink.value;
  if (hash.startsWith("http")) {
    let u = new URL(hash);
    hash = u.pathname.substring(u.pathname.lastIndexOf("/s/") + 3);
  }
  if (password.value.trim() === "") {
    alert("请输入原密码进行验证！");
    return;
  }
  let uri = generateURI();
  if (uri === "") {
    return;
  }
  axios
    .put(
      "./short",
      {
        hash: hash,
        url: uri,
        password: password.value.trim(),
      },
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    )
    .then((response) => {
      alert(`短链 ${hash} 更新成功！`);
    })
    .catch((error) => {
      console.log(error);
      if (error.response && error.response.status === 401) {
        alert("密码错误，请输入正确的原密码！");
      } else if (error.response && error.response.data) {
        alert(error.response.data);
      } else {
        alert("更新短链失败，请重试！");
      }
    });
}


// 主题切换功能
function initTheme() {
  const html = document.querySelector('html');
  const themeIcon = document.getElementById('theme-icon');
  let theme;

  // 从localStorage获取用户偏好的主题
  const savedTheme = localStorage.getItem('theme');

  if (savedTheme) {
    // 如果用户之前设置过主题，使用保存的主题
    theme = savedTheme;
  } else {
    // 如果没有设置过，检测系统主题偏好
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    theme = prefersDark ? 'dark' : 'light';
  }

  // 设置主题
  html.setAttribute('data-bs-theme', theme);

  // 更新图标
  if (theme === 'dark') {
    themeIcon.textContent = '☀️';
  } else {
    themeIcon.textContent = '🌙';
  }
}

function toggleTheme() {
  const html = document.querySelector('html');
  const currentTheme = html.getAttribute('data-bs-theme');
  // 切换主题
  const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
  html.setAttribute('data-bs-theme', newTheme);

  // 更新图标
  if (newTheme === 'dark') {
    themeIcon.textContent = '☀️';
  } else {
    themeIcon.textContent = '🌙';
  }

  // 保存用户偏好到localStorage
  localStorage.setItem('theme', newTheme);
}

listenInput();
initTheme();