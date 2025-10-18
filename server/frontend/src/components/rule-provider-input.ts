import { LitElement, html, unsafeCSS } from "lit";
import { customElement, state } from "lit/decorators.js";
import type { RuleProvider } from "../interface";
import globalStyles from "../index.css?inline";

@customElement("rule-provider-input")
export class RuleProviderInput extends LitElement {
  static styles = [unsafeCSS(globalStyles)];

  _ruleProviders: Array<RuleProvider> = [];

  @state()
  set ruleProviders(value) {
    this.dispatchEvent(
      new CustomEvent("change", {
        detail: value,
      })
    );
    this._ruleProviders = value;
  }

  get ruleProviders() {
    return this._ruleProviders;
  }

  RuleProviderTemplate(index: number) {
    return html`
      <div class="join mb-1">
        <div class="tooltip" data-tip="不能重复">
          <input
            class="input join-item"
            placeholder="名称"
            .value="${this.ruleProviders![index].name}"
            @change="${(e: Event) => {
              const target = e.target as HTMLInputElement;
              let updatedRuleProviders = this.ruleProviders;
              updatedRuleProviders![index].name = target.value;
              this.ruleProviders = updatedRuleProviders;
            }}" />
        </div>
        <div class="tooltip" data-tip="类型">
          <select
            class="select join-item w-fit"
            .value="${this.ruleProviders![index].behavior}"
            @change="${(e: Event) => {
              const target = e.target as HTMLInputElement;
              let updatedRuleProviders = this.ruleProviders;
              updatedRuleProviders![index].behavior = target.value;
              this.ruleProviders = updatedRuleProviders;
            }}">
            <option value="classical" selected>classical</option>
            <option value="domain">domain</option>
            <option value="ipcidr">ipcidr</option>
          </select>
        </div>
        <div>
          <input
            class="input join-item"
            placeholder="Url"
            .value="${this.ruleProviders![index].url}"
            @change="${(e: Event) => {
              const target = e.target as HTMLInputElement;
              let updatedRuleProviders = this.ruleProviders;
              updatedRuleProviders![index].url = target.value;
              this.ruleProviders = updatedRuleProviders;
            }}" />
        </div>
        <input
          class="input join-item"
          placeholder="出站策略组"
          .value="${this.ruleProviders![index].group}"
          @change="${(e: Event) => {
            const target = e.target as HTMLInputElement;
            let updatedRuleProviders = this.ruleProviders;
            updatedRuleProviders![index].group = target.value;
            this.ruleProviders = updatedRuleProviders;
          }}" />
        <div class="tooltip" data-tip="是否置于规则列表最前">
          <select
            class="select join-item w-fit"
            .value="${String(this.ruleProviders![index].prepend)}"
            @change="${(e: Event) => {
              const target = e.target as HTMLInputElement;
              let updatedRuleProviders = this.ruleProviders;
              updatedRuleProviders![index].prepend = Boolean(target.value);
              this.ruleProviders = updatedRuleProviders;
            }}">
            <option value="true">是</option>
            <option value="false" selected>否</option>
          </select>
        </div>
        <button
          class="btn join-item bg-error"
          type="button"
          @click="${() => {
            let updatedRuleProviders = this.ruleProviders?.filter(
              (_, i) => i !== index
            );
            this.ruleProviders = updatedRuleProviders;
          }}">
          删除
        </button>
      </div>
    `;
  }

  render() {
    return html` <!-- Rule Provider -->
      <div class="form-control mb-3">
        <label class="label mb-1 pl-1">
          <span class="label-text">Rule Provider</span>
          <button
            class="btn btn-primary btn-xs"
            type="button"
            @click="${() => {
              let updatedRuleProviders = this.ruleProviders
                ? [...this.ruleProviders]
                : [];
              updatedRuleProviders.push({
                behavior: "classical",
                url: "",
                name: "",
                prepend: false,
                group: "",
              });
              this.ruleProviders = updatedRuleProviders;
            }}">
            +
          </button>
        </label>
      </div>

      <div class="mb-3">
        ${this.ruleProviders?.map((_, i) => this.RuleProviderTemplate(i))}
      </div>`;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    "rule-provider-input": RuleProviderInput;
  }
}
