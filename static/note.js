// sidebar
function toggleSidebar() {
  document.getElementById("sidebar").classList.toggle("show");
}

// delete note
function deleteNote(button) {
  const noteId = button.getAttribute("data-note-id");
  if (!confirm("Are you sure you want to delete this note?")) return;
  fetch(`/api/notes/${noteId}`, { method: "DELETE" })
    .then(res => res.json())
    .then(data => {
      if (data.message === "Note deleted") {
        window.location.href = "/";
      } else {
        alert("Failed to delete the note.");
      }
    })
    .catch(err => console.error("Error:", err));
}

// share note
function shareNote(button) {
  const noteId = button.getAttribute("data-note-id");
  fetch(`/api/notes/share/${noteId}`, { method: "POST" })
    .then(res => res.json())
    .then(data => {
      if (!data.code) {
        return alert("Failed to retrieve share code.");
      }
      // copy to clipboard
      navigator.clipboard.writeText(data.code)
        .then(() => {
          document.getElementById("shareCodeText").textContent = data.code;
          const shareModalEl = document.getElementById("shareModal");
          const shareModal = new bootstrap.Modal(shareModalEl);
          shareModal.show();
          shareModalEl.focus();
        })
        .catch(() => {
          alert("Failed to copy share code to clipboard.");
        });
    })
    .catch(() => {
      alert("Error fetching share code.");
    });
}