import { useEffect, useState } from "react";
import { makeQuery } from "./fetch";

type Response<T> =
  | { status: "fetching" }
  | { status: "done"; data: T | null }
  | { status: "error" };

export const useQuery = <T,>(url: string): Response<T> => {
  const [response, setResponse] = useState<Response<T>>({
    status: "fetching",
  });

  useEffect(() => {
    const fetchData = async () => {
      const response = await makeQuery(url);
      setResponse({ status: "done", data: response.data });
    };

    fetchData();
  }, []);

  return response;
};
