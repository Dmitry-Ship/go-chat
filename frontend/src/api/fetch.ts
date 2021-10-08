export type Request =
  | {
      method: "GET";
    }
  | {
      method: "POST";
      body?: Record<string, any>;
    };

export const makeRequest = async (
  url: string,
  request: Request = { method: "GET" }
) => {
  if (import.meta.env.DEV) {
    url = "/api" + url;
  } else {
    url = String(import.meta.env.VITE_DOMAIN) + url;
  }

  try {
    let req: Record<string, any> = {};
    switch (request.method) {
      case "GET":
        req = {
          method: "GET",
        };
        break;

      case "POST":
        req = {
          method: "POST",
          body: request.body ? JSON.stringify(request.body) : undefined,
        };
        break;

      default:
        break;
    }

    const result = await fetch(url, {
      headers: {
        "Content-Type": "application/json",
      },
      ...req,
    });

    if (result.status === 200) {
      return result.json();
    } else {
      throw new Error(result.statusText);
    }
  } catch (error) {
    console.error(error);
    return null;
  }
};
