import { LitElement, html, unsafeCSS } from "lit";
import { customElement, state } from "lit/decorators.js";
import globalStyles from "../index.css?inline";
import type { Rename } from "../interface";

@customElement("rename-input")
export class RenameInput extends LitElement {
  static styles = [unsafeCSS(globalStyles)];

  private _replaceArray: Array<Rename> = [];

  @state()
  set replaceArray(value: Array<Rename>) {
    this._replaceArray = value;
    let updatedReplaceMap: { [key: string]: string } = {};
    value.forEach((e) => {
      updatedReplaceMap[e.old] = e.new;
    });
    this.dispatchEvent(
      new CustomEvent("change", {
        detail: updatedReplaceMap,
      })
    );
  }

  get replaceArray(): Array<Rename> {
    return this._replaceArray;
  }

  render() {
    return html`<!-- Rename -->
      <div class="form-control mb-3">
        <label class="label mb-1 pl-1">
          <span class="label-text">节点名称替换</span>
          <button
            class="btn btn-primary btn-xs"
            type="button"
            @click="${() => {
              let updatedReplaceArray = [...this.replaceArray];
              updatedReplaceArray.push({ old: "", new: "" });
              this.replaceArray = updatedReplaceArray;
            }}">
            +
          </button>
        </label>
      </div>

      <div class="mb-3">
        ${this.replaceArray.map((_, i) => this.RenameTemplate(i))}
      </div>`;
  }

  RenameTemplate(index: number) {
    const replaceItem = this.replaceArray[index];
    return html`<div class="join mb-1">
      <input
        class="input join-item"
        placeholder="旧名称 (正则表达式)"
        .value="${replaceItem?.old ?? ""}"
        @change="${(e: Event) => {
          const target = e.target as HTMLInputElement;
          let updatedReplaceArray = [...this.replaceArray];
          updatedReplaceArray[index] = {
            ...updatedReplaceArray[index],
            old: target.value,
          };
          this.replaceArray = updatedReplaceArray;
        }}" />
      <input
        class="input join-item"
        placeholder="新名称"
        .value="${replaceItem?.new ?? ""}"
        @change="${(e: Event) => {
          const target = e.target as HTMLInputElement;
          let updatedReplaceArray = [...this.replaceArray];
          updatedReplaceArray[index] = {
            ...updatedReplaceArray[index],
            new: target.value,
          };
          this.replaceArray = updatedReplaceArray;
        }}" />
      <button
        class="btn join-item bg-error"
        type="button"
        @click="${() => {
          let updatedReplaceArray = this.replaceArray.filter(
            (_, i) => i !== index
          );
          this.replaceArray = updatedReplaceArray;
        }}">
        删除
      </button>
    </div>`;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    "rename-input": RenameInput;
  }
}
