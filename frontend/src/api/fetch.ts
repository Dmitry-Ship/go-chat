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
  url = "/api" + url;
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

export const makePaginatedQuery = async (
  url: string,
  page: number,
  pageSize = 50
): Promise<{ status: boolean; data: any; nextPage: number }> => {
  const paginationParams =
    (url.includes("?") ? "&" : "?") + "page=" + page + "&page_size=" + pageSize;

  const result = await makeQuery(url + paginationParams);

  const nextPage = result.data.status ? page + 1 : page;

  return { ...result, nextPage };
};

export const makeCommand = async (
  url: string,
  body?: Record<string, any>
): Promise<{ status: boolean; data: any }> => {
  url = "/api" + url;
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
