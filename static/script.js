document
  .getElementById("shorten-form")
  .addEventListener("submit", async function (event) {
    event.preventDefault();

    const longUrlInput = document.getElementById("long-url");
    const longUrl = longUrlInput.value;

    if (!longUrl) {
      alert("Please enter a valid URL.");
      return;
    }

    const response = await fetch("/shorten", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ url: longUrl }),
    });

    if (response.ok) {
      const data = await response.json();
      const shortUrl = `${window.location.origin}/redirect/${data.short_url}`;

      const resultDiv = document.getElementById("result");
      const shortUrlAnchor = document.getElementById("short-url");

      shortUrlAnchor.href = shortUrl;
      shortUrlAnchor.textContent = shortUrl;

      resultDiv.classList.remove("hidden");
      longUrlInput.value = ""; // Clear the input field
    } else {
      alert("Failed to shorten URL. Please try again.");
    }
  });
