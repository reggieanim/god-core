export const waitForElement = async (selector, timeout = 30, page) => {
  const startTime = Date.now();
  while (Date.now() - startTime < timeout * 1e3) {
    const element = page.querySelector(selector);
    if (element) {
      return element;
    }
    await new Promise((resolve) => setTimeout(resolve, timeout));
  }
  throw new Error(`Timeout waiting for element: ${selector}`);
};
export const observeAndWaitForElement = async (selector, timeout = 5e3, page) => {
  return new Promise((resolve, reject) => {
    const element = page.querySelector(selector);
    if (element) {
      resolve(element);
    }
    const observer = new MutationObserver((_mutationsList, observer2) => {
      const element2 = page.querySelector(selector);
      if (element2) {
        observer2.disconnect();
        resolve(element2);
      }
    });
    observer.observe(page, { childList: true, subtree: true });
    setTimeout(() => {
      observer.disconnect();
      reject(new Error(`Timeout waiting for element: ${selector}`));
    }, timeout * 1e3);
  });
};
export const insertCustomBanner = (value) => {
  if (!document.getElementById("customBanner")) {
    document.body.insertAdjacentHTML(
      "beforeend",
      '<style>@import url("https://fonts.googleapis.com/css2?family=Roboto:wght@400;500&display=swap"); .mui-button { display: inline-block; padding: 10px 20px; font-size: 13px; color: black; text-transform: uppercase; background-color: #fff; border: 1px solid black; border-radius: 3px; cursor: pointer; font-family: "Roboto", sans-serif; transition: background-color 0.3s, color 0.3s, border-color 0.3s; } .mui-button:hover { background-color: black; color: #fff; border-color: black; } #customBanner { position: fixed; top: 50%; right: -200px; transform: translateY(-50%); width: 250px; background-color: #fff; color: #333; text-align: center; padding: 10px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.08); z-index: 9999; transition: right 0.5s ease-out; font-family: "Roboto", sans-serif; } #logo { width: 150px; height: auto; margin-bottom: 10px; } .log { font-size: 12px; text-transform: uppercase; color: rgb(132, 81, 225); opacity: 0.5; }</style><div id="customBanner"><p class="log">My Approval Engine</p><button id="startAutofill" class="mui-button">' + value + "</button></div>"
    );
    setTimeout(function() {
      document.getElementById("customBanner").style.right = "0";
    }, 1e3);
  }
};
export const removeCustomBanner = () => document.getElementById("startAutofill").addEventListener("click", function() {
  window.startAutofill = true;
  const element = document.getElementById("customBanner");
  if (element) {
    element.remove();
  }
});
export const setWindowToFalse = () => {
  window.startAutofill = false;
};
export const until = (conditionFunction) => {
  const poll = (resolve) => {
    if (conditionFunction()) resolve();
    else setTimeout(() => poll(resolve), 400);
  };
  return new Promise(poll);
};
export const clearStorage = async () => {
  try {
    await chrome.storage.local.clear();
    await chrome.storage.session.clear();
    await chrome.storage.sync.clear();
  } catch (error) {
    console.error("Error clearing storage:", error);
  }
};
