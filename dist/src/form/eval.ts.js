import { waitForElement } from "/src/helpers/functions/functions.ts.js";
export class Eval {
  constructor() {
  }
  conditionalEvaluate = async (instruction, page) => {
    const result = this.detectFieldPresence(instruction, page);
    if (result !== void 0) {
    }
  };
  detectFieldPresence = async (instruction, page) => {
    try {
      const element = await waitForElement(instruction.field, instruction.timeout, page);
      if (element instanceof HTMLElement) {
        const propertyOrMethod = instruction.evalExpression;
        if (propertyOrMethod in element) {
          const result = typeof element[propertyOrMethod] === "function" ? element[propertyOrMethod]() : element[propertyOrMethod];
          return result === instruction.value;
        } else {
          throw new Error(`Property or method ${propertyOrMethod} does not exist on the element`);
        }
      }
    } catch (err) {
      console.error(err);
      console.log(
        `Error checking presence of element: ${instruction.field}, with property: ${instruction.evalExpression}, that has the value: ${instruction.value}`
      );
      return false;
    }
  };
  evaluate = async (instruction, page) => {
    function evaluate(element, codeToExecute) {
      const func = new Function("return " + codeToExecute).bind(element);
      return func();
    }
    try {
      const element = await waitForElement(instruction.field, instruction.timeout, page);
      if (element instanceof HTMLElement) {
        const evaluation = evaluate(element, instruction.evalExpression);
        return evaluation();
      }
    } catch (err) {
      console.error(err);
      console.log(
        `Error executing code on element: ${instruction.field} with function: ${instruction.evalExpression}`
      );
    }
  };
}
