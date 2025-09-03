browser.browserAction.onClicked.addListener((tab) => {
    const eventData = {
        type: "job_application",
        source: "firefox-extension",
        timestamp: new Date().toISOString(),
        title: tab.title,
        url: tab.url
    };
    
    fetch("http://localhost:8080/ingest", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(eventData)
    }).then(response => {
        console.log("Posted to the app:", response.status);
    }).catch(err => {
        console.error("Failed to POST to the app", err);
    });
});