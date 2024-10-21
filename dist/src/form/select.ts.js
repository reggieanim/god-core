import { waitForElement } from "/src/helpers/functions/functions.ts.js";
export class Select {
  constructor() {
  }
  inputSelect = async (instruction, page) => {
    const element = await waitForElement(instruction.field, instruction.timeout, page);
    if (element instanceof HTMLElement) {
      element.focus();
      return this.select(element, instruction.value);
    }
    throw new Error("Could not find element or Element is not an HTML Element.");
  };
  select(element, selectors) {
    const options = element.querySelectorAll("option");
    const matcher = (item) => item.matches(selectors);
    let found = false;
    const matchingOption = Array.from(options).find(matcher);
    if (matchingOption) {
      matchingOption.selected = true;
      found = true;
    }
    element.dispatchEvent(new Event("input", { bubbles: true }));
    element.dispatchEvent(new Event("change", { bubbles: true }));
    return found;
  }
}
