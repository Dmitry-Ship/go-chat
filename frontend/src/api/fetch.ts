const returnResult = async (response: Response) => {
  if (response.ok) {
    const data = await response.json();
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
};

export const makeQuery = async (
  url: string
): Promise<{ status: boolean; data: any }> => {
  url = process.env.NEXT_PUBLIC_API_URL + "/api" + url;
  try {
    const result = await fetch(url, {
      method: "GET",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    });

    return await returnResult(result);
  } catch (error) {
    console.error(error);
    return {
      status: false,
      data: null,
    };
  }
};

export const makeCommand = async (
  url: string,
  body?: Record<string, any>
): Promise<{ status: boolean; data: any }> => {
  url = process.env.NEXT_PUBLIC_API_URL + "/api" + url;
  try {
    const result = await fetch(url, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: body ? JSON.stringify(body) : null,
    });

    return await returnResult(result);
  } catch (error) {
    console.error(error);
    return {
      status: false,
      data: null,
    };
  }
};
