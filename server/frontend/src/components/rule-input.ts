import { LitElement, html, unsafeCSS } from "lit";
import { customElement, state } from "lit/decorators.js";
import type { Rule } from "../interface";
import globalStyles from "../index.css?inline";

@customElement("rule-input")
export class RuleInput extends LitElement {
  static styles = [unsafeCSS(globalStyles)];

  _rules: Array<Rule> = [];
  @state()
  set rules(value: Array<Rule>) {
    this.dispatchEvent(
      new CustomEvent("change", {
        detail: value,
      })
    );
    this._rules = value;
  }
  get rules() {
    return this._rules;
  }
  render() {
    return html`<!-- Rule -->
      <div class="form-control mb-3">
        <label class="label mb-1 pl-1">
          <span class="label-text">规则</span>
          <button
            class="btn btn-primary btn-xs"
            type="button"
            @click="${() => {
              let updatedRules = this.rules ? [...this.rules] : [];
              updatedRules?.push({
                rule: "",
                prepend: false,
              });
              this.rules = updatedRules;
            }}">
            +
          </button>
        </label>
      </div>

      <div class="mb-3">
        ${this.rules?.map((_, i) => this.RuleTemplate(i))}
      </div>`;
  }

  RuleTemplate(index: number) {
    return html`<div class="join mb-1">
      <input
        class="input join-item"
        placeholder="规则"
        .value="${this.rules![index].rule}"
        @change="${(e: Event) => {
          const target = e.target as HTMLInputElement;
          let updatedRules = this.rules;
          updatedRules![index].rule = target.value;
          this.rules = updatedRules;
        }}" />
      <div class="tooltip" data-tip="是否置于规则列表最前">
        <select
          class="select join-item w-fit"
          .value="${String(this.rules![index].prepend)}"
          @change="${(e: Event) => {
            const target = e.target as HTMLInputElement;
            let updatedRules = this.rules;
            updatedRules![index].prepend = Boolean(target.value);
            this.rules = updatedRules;
          }}">
          <option value="true">是</option>
          <option value="false" selected>否</option>
        </select>
      </div>
      <button
        class="btn join-item bg-error"
        type="button"
        @click="${() => {
          let updatedRules = this.rules?.filter((_, i) => i !== index);
          this.rules = updatedRules;
        }}">
        删除
      </button>
    </div>`;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    "rule-input": RuleInput;
  }
}
