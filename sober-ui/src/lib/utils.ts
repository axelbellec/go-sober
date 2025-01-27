
import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"


export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export async function fetchWithAuth(url: string, options: RequestInit = {}) {
  const tokenKey = process.env.NEXT_PUBLIC_LOCALSTORAGE_TOKEN_KEY!;
  const token = localStorage.getItem(tokenKey);
  const headers = {
    ...options.headers,
    Authorization: `Bearer ${token}`,
  };

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (response.status === 401) {
    // Token expired or invalid
    localStorage.removeItem(tokenKey);
    window.location.href = "/login";
    throw new Error("Unauthorized");
  }

  return response;
}



