import {
  addListenersForStartingUrls,
  CreateNewWindowOrTab,
  getStartingURLs,
} from "./helpers/functions/serviceWorker";
import { Instruction } from "./types/types";

document.getElementById("startButton")!.addEventListener("click", async () => {
  try {
    const response = await fetch("/bb.json");
    const rawInstructions = await response.text();
    const parsedInstructions: Instruction[] = JSON.parse(rawInstructions);
    console.log(parsedInstructions);

    const startingUrls = getStartingURLs(parsedInstructions);
    console.log(startingUrls);

    for (const url of startingUrls) {
      await chrome.storage.session.set({ [`startingUrl_${url}`]: url });
      await chrome.storage.session.set({ [`instructions_${url}`]: parsedInstructions });
    }

    await addListenersForStartingUrls();
    await CreateNewWindowOrTab(parsedInstructions);
  } catch (error) {
    console.error("Error starting instruction processor:", error);
  }
});
