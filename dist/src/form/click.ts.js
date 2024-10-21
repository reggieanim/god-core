import { waitForElement } from "/src/helpers/functions/functions.ts.js";
export class Click {
  constructor() {
  }
  executeLeftClick = async (instruction, page) => {
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
  executeRightClick = async (instruction, page) => {
    try {
      const element = await waitForElement(instruction.field, instruction.timeout, page);
      if (element instanceof HTMLElement) {
        const event = new MouseEvent("contextmenu", {
          bubbles: true,
          cancelable: true,
          view: window,
          button: 2,
          buttons: 2,
          clientX: element.getBoundingClientRect().x,
          clientY: element.getBoundingClientRect().y
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
