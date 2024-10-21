import { waitForElement } from "../helpers/functions/functions";
import { FormInstructions } from "../types/types";

export class Select {
  constructor() {}

  public inputSelect = async (instruction: FormInstructions, page: Document) => {
    const element = await waitForElement(instruction.field, instruction.timeout, page);

    if (element instanceof HTMLElement) {
      element.focus();
      return this.select(element, instruction.value);
    }

    throw new Error("Could not find element or Element is not an HTML Element.");
  };

  private select(element: HTMLElement, selectors: string): boolean {
    const options = element.querySelectorAll<HTMLOptionElement>("option");
    const matcher = (item: HTMLOptionElement) => item.matches(selectors);

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
