// sidebar toggle
function toggleSidebar() {
  document.getElementById("sidebar").classList.toggle("show");
}

// create new note
function createNote() {
  fetch("/api/notes", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title: "New Note", content: "" })
  })
  .then(r => r.json())
  .then(d => { if (d.id) window.location.href = "/note/" + d.id });
}

// import handler
document.getElementById("doImportBtn").addEventListener("click", () => {
  const code = document.getElementById("importCode").value.trim();
  if (!/^\d{8}$/.test(code)) {
    return alert("Please enter a valid 8-digit code.");
  }
  fetch("/api/notes/import", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ code })
  })
  .then(r => r.json())
  .then(d => {
    if (d.id) {
      const m = bootstrap.Modal.getInstance(document.getElementById("importModal"));
      m.hide();
      location.reload();
    } else {
      alert("Import failed");
    }
  })
  .catch(() => alert("Error importing"));
});

const grid = document.getElementById("noteGrid");

new Sortable(grid, {
  animation: 150,
  onEnd: function (evt) {
    const newOrder = Array.from(grid.children).map((el, index) => ({
      id: parseInt(el.dataset.id),
      order: index
    }));

    fetch("/api/notes/order", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(newOrder)
    });
  }
});