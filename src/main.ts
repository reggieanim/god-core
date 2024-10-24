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

    for (const url of startingUrls) {
      await chrome.storage.session.set({ [`startingUrl_${url}`]: url });
      await chrome.storage.session.set({ [`instructions_${url}`]: rawInstructions });
    }

    await addListenersForStartingUrls();
    await CreateNewWindowOrTab(rawInstructions);
  } catch (error) {
    console.error("Error starting instruction processor:", error);
  }
});
