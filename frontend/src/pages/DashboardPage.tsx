// src/pages/DashboardPage.tsx
import { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { getBookmarks, createBookmark, deleteBookmark } from "../api/bookmarks";
import { useAuth } from "../hooks/useAuth";
import type { Bookmark } from "../types";

export default function DashboardPage() {
  const { isAuthed, signout } = useAuth();
  const navigate = useNavigate();
  const [bookmarks, setBookmarks] = useState<Bookmark[]>([]);
  const [url, setUrl] = useState("");
  const [tags, setTags] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

    const fetchBookmarks = useCallback(async () => {
    try {
        const { data } = await getBookmarks();
        setBookmarks(data.bookmarks);  // was: data
    } catch {
        setError("Failed to load bookmarks");
    }
    }, []);

  useEffect(() => {
    if (!isAuthed) { navigate("/login"); return; }
    fetchBookmarks();
  }, [isAuthed, navigate, fetchBookmarks]);

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url.trim()) return;
    setLoading(true);
    setError(null);
    try {
      const tagList = tags.split(",").map(t => t.trim()).filter(Boolean);
      await createBookmark({ url, tags: tagList });
      setUrl("");
      setTags("");
      await fetchBookmarks();
    } catch {
      setError("Failed to add bookmark");
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteBookmark(id);
      setBookmarks(prev => prev.filter(b => b.id !== id));
    } catch {
      setError("Failed to delete");
    }
  };

  const handleSignout = () => { signout(); navigate("/login"); };

  return (
    <div className="dashboard">
      <header>
        <h1>LinkVault</h1>
        <button onClick={handleSignout}>Sign out</button>
      </header>

      <form onSubmit={handleAdd} className="add-form">
        <input
          type="url"
          placeholder="https://example.com"
          value={url}
          onChange={e => setUrl(e.target.value)}
          required
        />
        <input
          type="text"
          placeholder="tags (comma separated)"
          value={tags}
          onChange={e => setTags(e.target.value)}
        />
        <button type="submit" disabled={loading}>
          {loading ? "Adding..." : "Add"}
        </button>
      </form>

      {error && <p className="error">{error}</p>}

      <ul className="bookmark-list">
        {bookmarks.map(b => (
          <li key={b.id} className="bookmark-item">
            <div>
              <a href={b.url} target="_blank" rel="noopener noreferrer">
                {b.title || b.url}
              </a>
              {b.description && <p>{b.description}</p>}
              {b.tags?.length > 0 && (
                <div className="tags">
                  {b.tags.map(t => <span key={t} className="tag">{t}</span>)}
                </div>
              )}
            </div>
            <button onClick={() => handleDelete(b.id)}>✕</button>
          </li>
        ))}
        {bookmarks.length === 0 && <li className="empty">No bookmarks yet.</li>}
      </ul>
    </div>
  );
}