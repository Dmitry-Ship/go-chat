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
): Promise<{ status: boolean; data: any }> => {
  url = String(import.meta.env.VITE_DOMAIN) + url;
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
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      ...req,
    });

    if (result.status === 200) {
      const data = await result.json();
      return {
        status: true,
        data,
      };
    } else {
      return {
        status: false,
        data: null,
      };
    }
  } catch (error) {
    console.error(error);
    return {
      status: false,
      data: null,
    };
  }
};
