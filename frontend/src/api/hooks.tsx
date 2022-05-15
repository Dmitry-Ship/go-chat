import { useEffect, useState } from "react";
import { makePaginatedQuery, makeQuery } from "./fetch";

type Response<T> =
  | { status: "fetching" }
  | { status: "done"; data: T }
  | { status: "error" };

export const useQuery = <T,>(
  url: string
): [Response<T>, (changes: any) => void] => {
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

  return [
    response,
    (changes) =>
      setResponse((prevState) => {
        if (prevState.status === "done") {
          return { ...prevState, data: { ...prevState.data, ...changes } };
        }
        return prevState;
      }),
  ];
};

export const useQueryOnDemand = <T,>(
  url: string
): [Response<T>, () => void] => {
  const [response, setResponse] = useState<Response<T>>({
    status: "fetching",
  });

  return [
    response,
    async () => {
      const response = await makeQuery(url);

      if (response.status) {
        setResponse({ status: "done", data: response.data });
      } else {
        setResponse({ status: "error" });
      }
    },
  ];
};

type PaginatedResponse<T> = {
  status: "fetching" | "done" | "error";
  items: T[];
};

type appendFunc<T> = (items: T[]) => void;
type loadNextFunc = () => void;

export const usePaginatedQuery = <T,>(
  url: string,
  inverted: boolean = false
): [PaginatedResponse<T>, appendFunc<T>, loadNextFunc] => {
  const [page, setPage] = useState(1);

  const [currentResponse, setCurrentResponse] = useState<PaginatedResponse<T>>({
    status: "fetching",
    items: [],
  });

  useEffect(() => {
    const fetchData = async (url: string, page: number) => {
      const response = await makePaginatedQuery(url, page);

      if (!response.status) {
        setCurrentResponse((prevState) => ({
          ...prevState,
          status: "error",
        }));

        return;
      }

      setCurrentResponse((prevState) => ({
        status: "done",
        items: inverted
          ? [...response.data, ...prevState.items]
          : [...prevState.items, ...response.data],
      }));
    };

    fetchData(url, page);
  }, [page]);

  return [
    currentResponse,
    (items: T[]) => {
      setCurrentResponse((prevState) => ({
        ...prevState,
        items: [...prevState.items, ...items],
      }));
    },
    () => {
      setPage((prevState) => prevState + 1);
    },
  ];
};
