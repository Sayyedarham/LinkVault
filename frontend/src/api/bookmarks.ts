// src/api/bookmarks.ts
import client from "./client";
import type { Bookmark } from "../types";

export const getBookmarks = () =>
  client.get<{ bookmarks: Bookmark[]; count: number }>("/bookmarks");

export const createBookmark = (data: { url: string; tags?: string[] }) =>
  client.post<Bookmark>("/bookmarks", data);

export const deleteBookmark = (id: string) =>
  client.delete(`/bookmarks/${id}`);