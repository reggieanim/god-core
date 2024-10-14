import { InstructionProcessor } from "./core/core.ts";

document.getElementById("startButton")!.addEventListener("click", async () => {
  try {
    const response = await fetch("/a.json");
    const rawInstructions = await response.text();
    const processor = new InstructionProcessor(rawInstructions);
    processor.start();
  } catch (error) {
    console.error("Error starting instruction processor:", error);
  }
});
