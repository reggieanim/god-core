import {
  addListenersForStartingUrls,
  CreateNewWindowOrTab,
  getStartingURLs,
} from "./helpers/functions/serviceWorker";

document.getElementById("startButton")!.addEventListener("click", async () => {
  try {
    const response = await fetch("/bb.json");
    const rawInstructions = await response.text();
    const startingUrls = getStartingURLs(rawInstructions);
    await chrome.storage.session.set({ startingUrls: startingUrls });
    await chrome.storage.session.set({ instructions: rawInstructions });
    await addListenersForStartingUrls();
    await CreateNewWindowOrTab(rawInstructions);
  } catch (error) {
    console.error("Error starting instruction processor:", error);
  }
});
