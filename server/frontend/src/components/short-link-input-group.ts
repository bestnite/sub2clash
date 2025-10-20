import { html, LitElement, unsafeCSS } from "lit";
import { customElement, property } from "lit/decorators.js";
import globalStyles from "../index.css?inline";

@customElement("short-link-input-group")
export class ShortLinkInputGroup extends LitElement {
  static styles = unsafeCSS(globalStyles);

  @property()
  id: string = "";

  @property({ type: Number })
  _screenSizeLevel: number = 0;

  @property()
  passwd: string = "";

  connectedCallback() {
    super.connectedCallback();
    window.addEventListener("resize", this._checkScreenSize);
    this._checkScreenSize(); // Initial check
  }

  disconnectedCallback() {
    window.removeEventListener("resize", this._checkScreenSize);
    super.disconnectedCallback();
  }

  _checkScreenSize = () => {
    const width = window.innerWidth;
    if (width < 365) {
      this._screenSizeLevel = 0; // sm
    } else if (width < 640) {
      this._screenSizeLevel = 1; // md
    } else {
      this._screenSizeLevel = 2; // other
    }
  };

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

  idInputTemplate() {
    return html`<input
      class="input input-bordered w-1/2 join-item"
      type="text"
      placeholder="ID（可选）"
      .value="${this.id}"
      @change="${(e: Event) => {
        this.id = (e.target as HTMLInputElement).value;
        this.dispatchEvent(
          new CustomEvent("id-change", {
            detail: this.id,
          })
        );
      }}" />`;
  }

  passwdInputTemplate() {
    return html`<input
      class="input input-bordered w-1/2 join-item"
      type="text"
      placeholder="密码"
      .value="${this.passwd}"
      @change="${(e: Event) => {
        this.passwd = (e.target as HTMLInputElement).value;
        this.dispatchEvent(
          new CustomEvent("passwd-change", {
            detail: this.passwd,
          })
        );
      }}" />`;
  }

  generateBtnTemplate(extraClass: string = "") {
    return html`<button
      class="btn btn-primary join-item ${extraClass}"
      type="button"
      @click="${(e: Event) => {
        this.dispatchEvent(
          new CustomEvent("generate-btn-click", { detail: e })
        );
      }}">
      生成短链
    </button>`;
  }

  updateBtnTemplate(extraClass: string = "") {
    return html`<button
      class="btn btn-primary join-item ${extraClass}"
      @click="${(e: Event) => {
        this.dispatchEvent(new CustomEvent("update-btn-click", { detail: e }));
      }}"
      type="button">
      更新短链
    </button>`;
  }

  deleteBtnTemplate(extraClass: string = "") {
    return html`<button
      class="btn btn-primary join-item ${extraClass}"
      @click="${(e: Event) => {
        this.dispatchEvent(new CustomEvent("delete-btn-click", { detail: e }));
      }}"
      type="button">
      删除短链
    </button>`;
  }

  copyBtnTemplate(extraClass: string = "") {
    return html`<button
      class="btn btn-primary join-item ${extraClass}"
      type="button"
      @click="${(e: Event) => {
        this.copyToClipboard(
          `${window.location.origin}${window.location.pathname}s/${this.id}?password=${this.passwd}`,
          e.target as HTMLButtonElement
        );
      }}">
      复制短链
    </button>`;
  }

  render() {
    const sm = html`<div class="form-control mb-2">
      <div class="join w-full mb-1">
        ${this.idInputTemplate()} ${this.passwdInputTemplate()}
      </div>
      <div class="join w-full mb-1">
        ${this.generateBtnTemplate("w-1/2")} ${this.updateBtnTemplate("w-1/2")}
      </div>
      <div class="join w-full">
        ${this.deleteBtnTemplate("w-1/2")} ${this.copyBtnTemplate("w-1/2")}
      </div>
    </div>`;

    const md = html`<div class="form-control mb-2">
      <div class="join w-full mb-1">
        ${this.idInputTemplate()} ${this.passwdInputTemplate()}
      </div>
      <div class="join w-full">
        ${this.generateBtnTemplate("w-1/4")} ${this.updateBtnTemplate("w-1/4")}
        ${this.deleteBtnTemplate("w-1/4")} ${this.copyBtnTemplate("w-1/4")}
      </div>
    </div>`;

    const other = html`<div class="form-control mb-2">
      <div class="join w-full">
        ${this.idInputTemplate()} ${this.passwdInputTemplate()}
        ${this.generateBtnTemplate()} ${this.updateBtnTemplate()}
        ${this.deleteBtnTemplate()} ${this.copyBtnTemplate()}
      </div>
    </div>`;

    switch (this._screenSizeLevel) {
      case 0:
        return sm;
      case 1:
        return md;
      default:
        return other;
    }
  }
}

declare global {
  interface HTMLElementTagNameMap {
    "short-link-input-group": ShortLinkInputGroup;
  }
}
