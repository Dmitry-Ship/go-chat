export const makeRequest = async (url: string, method = "GET", body?: any) => {
  if (import.meta.env.DEV) {
    url = "/api" + url;
  } else {
    url = String(import.meta.env.VITE_DOMAIN) + url;
  }

  const result = await fetch(url, {
    method,
    headers: {
      "Content-Type": "application/json",
    },
    body: body ? JSON.stringify(body) : undefined,
  });

  if (result.status === 200) {
    return result.json();
  }
};
