import { waitForElement } from "../helpers/functions/functions";
import { FormInstructions } from "../types/types";

export class Click {
  constructor() {}

  public executeLeftClick = async (instruction: FormInstructions, page: Document) => {
    try {
      const element = await waitForElement(instruction.field, instruction.timeout, page);

      if (element instanceof HTMLElement) {
        element.focus();
        element.click();
      }
    } catch (err) {
      console.error("Error finding element", err);
      console.log(`Error left clicking: ${instruction.field} when: ${instruction.description}`);
    }
  };

  public executeRightClick = async (instruction: FormInstructions, page: Document) => {
    try {
      const element = await waitForElement(instruction.field, instruction.timeout, page);

      if (element instanceof HTMLElement) {
        const event = new MouseEvent("contextmenu", {
          bubbles: true,
          cancelable: true,
          view: window,
          button: 2,
        });
        element.focus();
        element.dispatchEvent(event);
      }
    } catch (err) {
      console.error("Error finding element", err);
      console.log(`Error right clicking: ${instruction.field} when: ${instruction.description}`);
    }
  };
}
