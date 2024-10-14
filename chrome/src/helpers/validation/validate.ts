import { FormInstructions } from "../../types/types";

const kindMap: { [key: string]: boolean } = {
  text: true,
  notify: true,
  select: true,
  leftClick: true,
  rightClick: true,
  wait: true,
  eval: true,
  condEval: true,
};

class Validator {
  errors: string[] = [];

  check(condition: boolean, field: string, message: string) {
    if (!condition) {
      this.errors.push(`${field}: ${message}`);
    }
  }

  valid(): boolean {
    return this.errors.length === 0;
  }
}

export const validate = (ins: FormInstructions): boolean => {
  const isValidDescription = typeof ins.description === "string";
  const isValidField = typeof ins.field === "string";
  const isValidValue = typeof ins.value === "string";
  const insKind = ins.kind;
  const isValidKind = typeof insKind === "string";
  const doesNotMatchExpectedKind = insKind ? kindMap[insKind] : false;

  const v = new Validator();
  v.check(isValidDescription, "description", "Needs a 'description' property of string");
  v.check(isValidField, "field", "Needs a 'field' property of 'text'");
  v.check(isValidValue, "value", "Needs a 'value' property of string");
  v.check(isValidDescription, "description", "Needs a description property");
  v.check(isValidKind, "kind", "Needs a 'kind' property");
  v.check(
    doesNotMatchExpectedKind,
    "noKind",
    "Needs a 'kind' property of 'text' || condEval || select || leftClick || rightClick' || 'wait' || 'notify' || 'eval'"
  );

  if (!v.valid()) {
    for (const error of v.errors) {
      console.error(error);
    }
    console.log("Could not validate");
  }

  return v.valid();
};
