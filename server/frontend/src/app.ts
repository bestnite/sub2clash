import { LitElement, html, unsafeCSS } from "lit";
import { customElement, query, state } from "lit/decorators.js";
import globalStyles from "./index.css?inline";
import { type Config, type Rule, type RuleProvider } from "./interface.js";
import axios, { AxiosError } from "axios";
import { base64EncodeUnicode, base64decodeUnicode } from "./utils.js";
import "./components/rule-provider-input.js";
import "./components/rule-input.js";
import "./components/rename-input.js";

@customElement("sub2clash-app")
export class Sub2clashApp extends LitElement {
  static styles = [unsafeCSS(globalStyles)];

  private _config: Config = {
    clashType: 2,
    subscriptions: [],
    proxies: [],
    refresh: false,
    autoTest: false,
    lazy: false,
    ignoreCountryGroup: false,
    useUDP: false,
    template: "",
    sort: "nameasc",
    remove: "",
    nodeList: false,
    ruleProviders: [],
    replace: undefined,
    rules: [],
  };

  @state()
  set config(value: Config) {
    console.log(JSON.stringify(value));
    if (
      (value.subscriptions == null || value.subscriptions.length == 0) &&
      (value.proxies == null || value.proxies.length == 0)
    ) {
      this.configUrl = "";
      return;
    }
    const oldValue = this._config;
    this.configUrl = `${
      window.location.origin
    }${window.location.pathname.replace(
      /\/$/,
      ""
    )}/convert/${base64EncodeUnicode(JSON.stringify(value))
      .replace(/\+/g, "-")
      .replace(/\//g, "_")}`;
    this._config = value;
    this.requestUpdate("config", oldValue);
  }

  get config(): Config {
    return this._config;
  }

  @state({
    hasChanged(value: boolean) {
      localStorage.setItem("theme", value ? "dark" : "light");
      document
        .querySelector("html")
        ?.setAttribute("data-theme", value ? "dark" : "light");
      return true;
    },
  })
  darkTheme: boolean = this.initTheme();

  initTheme(): boolean {
    const savedTheme = localStorage.getItem("theme");
    if (savedTheme != null) {
      return savedTheme === "dark" ? true : false;
    }
    const prefersDark = window.matchMedia(
      "(prefers-color-scheme: dark)"
    ).matches;
    return prefersDark;
  }

  @state()
  reverseUrl: string = "";

  @state()
  dialogMessage: string = "";

  @state()
  dialogTitle: string = "";

  @query("dialog#my_modal")
  dialog!: HTMLDialogElement;

  showDialog(title: string, message: string): void {
    if (title.trim() === "") {
      title = "警告";
    }
    this.dialogTitle = title;
    this.dialogMessage = message;
    this.dialog.showModal();
  }

  @state()
  configUrl: string = "";

  @state()
  shortLinkID: string = "";

  @state()
  shortLinkPasswd: string = "";

  async copyToClipboard(content: string, e: HTMLButtonElement) {
    try {
      await navigator.clipboard.writeText(content);
      let text = e.textContent;
      e.addEventListener("mouseout", function () {
        e.textContent = text;
      });
      e.textContent = "复制成功";
    } catch (err) {
      console.error("复制到剪贴板失败:", err);
    }
  }

  generateShortLink() {
    if (this.configUrl === "") {
      this.showDialog("", "还未填写配置");
      return;
    }
    axios
      .post(
        "./short",
        {
          config: this.config,
          password: this.shortLinkPasswd,
          id: this.shortLinkID,
        },
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      )
      .then((response) => {
        // 设置返回的短链ID和密码
        this.shortLinkID = response.data.id;
        this.shortLinkPasswd = response.data.password;
      })
      .catch((error) => {
        if (error.response && error.response.data) {
          this.showDialog("", "生成短链失败：" + error.response.data);
        } else {
          this.showDialog("", "生成短链失败");
        }
      });
  }

  updateShortLink() {
    if (this.shortLinkID.trim() === "") {
      this.showDialog("", "请输入ID");
      return;
    }
    if (this.shortLinkPasswd.trim() === "") {
      this.showDialog("", "请输入密码");
      return;
    }
    if (this.configUrl === "") {
      this.showDialog("", "还未填写配置");
      return;
    }
    axios
      .put(
        "./short",
        {
          id: this.shortLinkID,
          config: this.config,
          password: this.shortLinkPasswd,
        },
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      )
      .then(() => {
        this.showDialog("成功", "更新成功");
      })
      .catch((error) => {
        if (error.response && error.response.status === 401) {
          this.showDialog("", "密码错误");
        } else if (error.response && error.response.data) {
          this.showDialog("", "更新短链失败：" + error.response.data);
        } else {
          this.showDialog("", "更新短链失败");
        }
      });
  }

  deleteShortLink() {
    if (this.shortLinkID.trim() === "") {
      this.showDialog("", "请输入ID");
      return;
    }
    if (this.shortLinkPasswd.trim() === "") {
      this.showDialog("", "请输入密码");
      return;
    }
    const params = new URLSearchParams();
    params.append("password", this.shortLinkPasswd);
    axios
      .delete(`./short/${this.shortLinkID}?${params.toString()}`, {
        headers: {
          "Content-Type": "application/json",
        },
      })
      .then(() => {
        this.showDialog("成功", "删除成功");
      })
      .catch((error) => {
        if (error.response && error.response.status === 401) {
          this.showDialog("", "短链不存在或密码错误");
        } else if (error.response && error.response.data) {
          this.showDialog("", "删除短链失败：" + error.response.data);
        } else {
          this.showDialog("", "删除短链失败");
        }
      });
  }

  getRawConfigFromShortLink() {
    const s = this.reverseUrl.split("/s/");
    if (s.length != 2) {
      this.showDialog("", "解析失败");
      return;
    }

    axios
      .get(`./short/${s[1]}`)
      .then((resp) => {
        this.config = resp.data;
      })
      .catch((err: AxiosError) => {
        if (err.response && err.response.status == 401) {
          this.showDialog("", "短链不存在或密码错误");
        } else if (err.response && err.response.data) {
          this.showDialog("", "获取配置失败：" + err.response.data);
        } else {
          this.showDialog("", "获取配置失败");
        }
      });
  }

  parseConfig() {
    if (this.reverseUrl.trim() === "") {
      this.showDialog("", "无法解析，链接为空");
    }
    if (this.reverseUrl.indexOf("/s/") != -1) {
      this.getRawConfigFromShortLink();
      return;
    }
    let url = new URL(this.reverseUrl);
    const pathSections = url.pathname.split("/");
    if (pathSections.length < 2) {
      this.showDialog("", "无法解析，链接格式错误");
    }
    if (pathSections[pathSections.length - 2] == "convert") {
      let base64Data = pathSections[pathSections.length - 1];
      base64Data = base64Data.replace(/-/g, "+").replace(/_/g, "/");
      try {
        const configData = base64decodeUnicode(base64Data);
        this.config = JSON.parse(configData) as Config;
      } catch (e: any) {
        this.showDialog("", "无法解析 Base64，配置格式错误");
        return;
      }
    } else {
      this.showDialog("", "无法解析，链接格式错误");
    }
  }

  render() {
    return html`
      <dialog id="my_modal" class="modal">
        <div class="modal-box">
          <h3 class="text-lg font-bold">${this.dialogTitle}</h3>
          <p class="py-4">${this.dialogMessage}</p>
          <div class="modal-action">
            <form method="dialog">
              <button class="btn">关闭</button>
            </form>
          </div>
        </div>
      </dialog>
      <div class="max-w-4xl mx-auto p-4 flex flex-col items-center">
        <form class="w-full max-w-2xl bg-base-100">
          <fieldset class="fieldset mb-6">
            <div class="flex flex-row justify-between items-center my-6">
              <legend
                class="fieldset-legend text-2xl font-semibold inline-block m-0 p-0">
                sub2clash
              </legend>
              <label class="swap swap-rotate h-7 w-7">
                <!-- this hidden checkbox controls the state -->
                <input
                  type="checkbox"
                  class="theme-controller"
                  .checked="${!this.darkTheme}"
                  @change="${() => (this.darkTheme = !this.darkTheme)}" />

                <!-- sun icon -->
                <svg
                  class="swap-off h-7 w-7 fill-current"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 24 24">
                  <path
                    d="M5.64,17l-.71.71a1,1,0,0,0,0,1.41,1,1,0,0,0,1.41,0l.71-.71A1,1,0,0,0,5.64,17ZM5,12a1,1,0,0,0-1-1H3a1,1,0,0,0,0,2H4A1,1,0,0,0,5,12Zm7-7a1,1,0,0,0,1-1V3a1,1,0,0,0-2,0V4A1,1,0,0,0,12,5ZM5.64,7.05a1,1,0,0,0,.7.29,1,1,0,0,0,.71-.29,1,1,0,0,0,0-1.41l-.71-.71A1,1,0,0,0,4.93,6.34Zm12,.29a1,1,0,0,0,.7-.29l.71-.71a1,1,0,1,0-1.41-1.41L17,5.64a1,1,0,0,0,0,1.41A1,1,0,0,0,17.66,7.34ZM21,11H20a1,1,0,0,0,0,2h1a1,1,0,0,0,0-2Zm-9,8a1,1,0,0,0-1,1v1a1,1,0,0,0,2,0V20A1,1,0,0,0,12,19ZM18.36,17A1,1,0,0,0,17,18.36l.71.71a1,1,0,0,0,1.41,0,1,1,0,0,0,0-1.41ZM12,6.5A5.5,5.5,0,1,0,17.5,12,5.51,5.51,0,0,0,12,6.5Zm0,9A3.5,3.5,0,1,1,15.5,12,3.5,3.5,0,0,1,12,15.5Z" />
                </svg>

                <!-- moon icon -->
                <svg
                  class="swap-on h-7 w-7 fill-current"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 24 24">
                  <path
                    d="M21.64,13a1,1,0,0,0-1.05-.14,8.05,8.05,0,0,1-3.37.73A8.15,8.15,0,0,1,9.08,5.49a8.59,8.59,0,0,1,.25-2A1,1,0,0,0,8,2.36,10.14,10.14,0,1,0,22,14.05,1,1,0,0,0,21.64,13Zm-9.5,6.69A8.14,8.14,0,0,1,7.08,5.22v.27A10.15,10.15,0,0,0,17.22,15.63a9.79,9.79,0,0,0,2.1-.22A8.11,8.11,0,0,1,12.14,19.73Z" />
                </svg>
              </label>
            </div>

            <!-- Input URL -->
            <div class="form-control mb-5">
              <label class="label mb-1 pl-1">解析链接</label>
              <div class="join w-full">
                <input
                  class="input input-bordered w-full join-item"
                  type="text"
                  @change="${(e: Event) => {
                    this.reverseUrl = (e.target as HTMLInputElement).value;
                  }}"
                  placeholder="通过生成的链接重新填写下方设置" />
                <button
                  class="btn btn-primary join-item"
                  @click="${this.parseConfig}"
                  type="button">
                  解析
                </button>
              </div>
            </div>

            <!-- API Endpoint -->
            <div class="form-control mb-3">
              <label class="label mb-1 pl-1" for="endpoint">客户端类型</label>
              <select
                class="select select-bordered w-full"
                name="endpoint"
                .value="${this.config.clashType == 1 ? "1" : "2"}"
                @change="${(e: Event) => {
                  this.config = {
                    ...this.config,
                    clashType: Number((e.target as HTMLInputElement).value),
                  };
                }}">
                <option value="1">Clash</option>
                <option value="2" selected>Clash.Meta</option>
              </select>
            </div>

            <!-- Template -->
            <div class="form-control mb-3">
              <label class="label mb-1 pl-1" for="template">模板链接</label>
              <input
                class="input input-bordered w-full"
                name="template"
                placeholder="输入模板链接"
                type="text"
                .value="${this.config.template ?? ""}"
                @change="${(e: Event) => {
                  this.config = {
                    ...this.config,
                    template: (e.target as HTMLInputElement).value,
                  };
                }}" />
            </div>

            <!-- Subscription Link -->
            <div class="form-control mb-3">
              <label class="label mb-1 pl-1" for="sub">订阅链接</label>
              <div>
                <textarea
                  class="textarea textarea-bordered h-24 w-full"
                  name="sub"
                  placeholder="每行输入一个订阅链接"
                  .value="${this.config.subscriptions
                    ? this.config.subscriptions.join("\n")
                    : ""}"
                  @change="${(e: Event) => {
                    this.config = {
                      ...this.config,
                      subscriptions: (e.target as HTMLInputElement).value
                        .split("\n")
                        .filter((e) => e.trim() !== ""),
                    };
                  }}"></textarea>
              </div>
            </div>

            <!-- Proxy Link -->
            <div class="form-control mb-3">
              <label class="label mb-1 pl-1" for="proxy">节点分享链接</label>
              <div>
                <textarea
                  class="textarea textarea-bordered h-24 w-full"
                  name="proxy"
                  placeholder="每行输入一个节点分享链接"
                  .value="${this.config.proxies
                    ? this.config.proxies.join("\n")
                    : ""}"
                  @change="${(e: Event) => {
                    this.config = {
                      ...this.config,
                      proxies: (e.target as HTMLInputElement).value
                        .split("\n")
                        .filter((e) => e.trim() !== ""),
                    };
                  }}"></textarea>
              </div>
            </div>

            <!-- User Agent -->
            <div class="form-control mb-3">
              <label class="label mb-1 pl-1" for="user-agent">UA 标识</label>
              <div>
                <textarea
                  class="textarea textarea-bordered h-20 w-full"
                  name="user-agent"
                  placeholder="用于获取订阅的 http 请求中的 User-Agent 标识"
                  .value="${this.config.userAgent ?? ""}"
                  @change="${(e: Event) => {
                    this.config = {
                      ...this.config,
                      userAgent: (e.target as HTMLInputElement).value,
                    };
                  }}"></textarea>
              </div>
            </div>

            <!-- Sort -->
            <div class="form-control mb-3">
              <label class="label mb-1 pl-1" for="sort">
                国家策略组排序规则
              </label>
              <select
                class="select select-bordered w-full"
                name="sort"
                .value="${this.config.sort ?? "nameasc"}"
                @change="${(e: Event) => {
                  this.config = {
                    ...this.config,
                    sort: (e.target as HTMLInputElement).value,
                  };
                }}">
                <option value="nameasc">名称（升序）</option>
                <option value="namedesc">名称（降序）</option>
                <option value="sizeasc">节点数量（升序）</option>
                <option value="sizedesc">节点数量（降序）</option>
              </select>
            </div>

            <!-- Remove -->
            <div class="form-control mb-3">
              <label class="label mb-1 pl-1" for="remove">
                <span class="label-text">排除节点</span>
              </label>
              <input
                class="input input-bordered w-full"
                type="text"
                name="remove"
                placeholder="正则表达式"
                .value="${this.config.remove ?? ""}"
                @change="${(e: Event) => {
                  this.config = {
                    ...this.config,
                    remove: (e.target as HTMLInputElement).value,
                  };
                }}" />
            </div>

            <!-- Checkboxes -->
            <div class="form-control mb-3">
              <label class="label cursor-pointer">
                <input
                  type="checkbox"
                  name="refresh"
                  class="checkbox"
                  .checked="${this.config.refresh ?? false}"
                  @change="${(e: Event) => {
                    this.config = {
                      ...this.config,
                      refresh: (e.target as HTMLInputElement).checked,
                    };
                  }}" />
                强制重新获取订阅
              </label>
            </div>
            <div class="form-control mb-3">
              <label class="label cursor-pointer">
                <input
                  type="checkbox"
                  name="nodeList"
                  class="checkbox"
                  .checked="${this.config.nodeList ?? false}"
                  @change="${(e: Event) => {
                    this.config = {
                      ...this.config,
                      nodeList: (e.target as HTMLInputElement).checked,
                    };
                  }}" />
                输出为 Node List
              </label>
            </div>
            <div class="form-control mb-3">
              <label class="label cursor-pointer">
                <input
                  type="checkbox"
                  name="autoTest"
                  class="checkbox"
                  .checked="${this.config.autoTest ?? false}"
                  @change="${(e: Event) => {
                    this.config = {
                      ...this.config,
                      autoTest: (e.target as HTMLInputElement).checked,
                    };
                  }}" />
                国家策略组自动测速
              </label>
            </div>
            <div class="form-control mb-3">
              <label class="label cursor-pointer">
                <input
                  type="checkbox"
                  name="lazy"
                  class="checkbox"
                  .checked="${this.config.lazy ?? false}"
                  @change="${(e: Event) => {
                    this.config = {
                      ...this.config,
                      lazy: (e.target as HTMLInputElement).checked,
                    };
                  }}" />
                自动测速启用 lazy 模式
              </label>
            </div>
            <div class="form-control mb-3">
              <label class="label cursor-pointer">
                <input
                  type="checkbox"
                  name="igcg"
                  class="checkbox"
                  .checked="${this.config.ignoreCountryGroup ?? false}"
                  @change="${(e: Event) => {
                    this.config = {
                      ...this.config,
                      ignoreCountryGroup: (e.target as HTMLInputElement)
                        .checked,
                    };
                  }}" />
                不输出国家策略组
              </label>
            </div>
            <div class="form-control mb-5">
              <label class="label cursor-pointer">
                <input
                  type="checkbox"
                  name="useUDP"
                  class="checkbox"
                  .checked="${this.config.useUDP ?? false}"
                  @change="${(e: Event) => {
                    this.config = {
                      ...this.config,
                      useUDP: (e.target as HTMLInputElement).checked,
                    };
                  }}" />
                使用 UDP
              </label>
            </div>

            <rule-provider-input
              @change="${(e: CustomEvent<Array<RuleProvider>>) => {
                this.config = {
                  ...this.config,
                  ruleProviders: e.detail,
                };
              }}"></rule-provider-input>

            <rule-input
              @change="${(e: CustomEvent<Array<Rule>>) => {
                this.config = {
                  ...this.config,
                  rules: e.detail,
                };
              }}"></rule-input>
            <rename-input
              @change="${(e: CustomEvent<{ [key: string]: string }>) => {
                this.config = {
                  ...this.config,
                  replace: e.detail,
                };
              }}"></rename-input>
          </fieldset>

          <fieldset class="fieldset mb-8">
            <legend
              class="fieldset-legend text-2xl font-semibold mb-4 text-center">
              输出配置
            </legend>

            <!-- Display the API Link -->
            <div class="form-control mb-5">
              <div class="join w-full mb-2">
                <input
                  class="input input-bordered w-full join-item cursor-not-allowed"
                  type="text"
                  placeholder="链接"
                  .value="${this.configUrl}"
                  readonly />
                <button
                  class="btn btn-primary join-item"
                  @click="${(e: Event) => {
                    this.copyToClipboard(
                      this.configUrl,
                      e.target as HTMLButtonElement
                    );
                  }}"
                  type="button">
                  复制链接
                </button>
              </div>
            </div>

            <div class="form-control mb-2">
              <div class="join w-full">
                <input
                  class="input input-bordered w-1/2 join-item"
                  type="text"
                  placeholder="ID（可选）"
                  .value="${this.shortLinkID}"
                  @change="${(e: Event) => {
                    this.shortLinkID = (e.target as HTMLInputElement).value;
                  }}" />
                <input
                  class="input input-bordered w-1/2 join-item"
                  type="text"
                  placeholder="密码"
                  .value="${this.shortLinkPasswd}"
                  @change="${(e: Event) => {
                    this.shortLinkPasswd = (e.target as HTMLInputElement).value;
                  }}" />
                <button
                  class="btn btn-primary join-item"
                  type="button"
                  @click="${this.generateShortLink}">
                  生成短链
                </button>
                <button
                  class="btn btn-primary join-item"
                  @click="${this.updateShortLink}"
                  type="button">
                  更新短链
                </button>
                <button
                  class="btn btn-primary join-item"
                  @click="${this.deleteShortLink}"
                  type="button">
                  删除短链
                </button>
                <button
                  class="btn btn-primary join-item"
                  type="button"
                  @click="${(e: Event) => {
                    this.copyToClipboard(
                      `${window.location.origin}${window.location.pathname}s/${this.shortLinkID}?password=${this.shortLinkPasswd}`,
                      e.target as HTMLButtonElement
                    );
                  }}">
                  复制短链
                </button>
              </div>
            </div>
          </fieldset>
        </form>
      </div>

      <footer class="footer footer-horizontal footer-center mb-8">
        <aside>
          <p>
            Powered by
            <a class="link" href="https://github.com/bestnite/sub2clash"
              >sub2clash</a
            >
          </p>
          <p>Version: ${import.meta.env.APP_VERSION ?? "dev"}</p>
        </aside>
      </footer>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    "sub2clash-app": Sub2clashApp;
  }
}
