import { reducer } from "next/dist/client/components/reducer";
import { useEffect, useState } from "react";
import { useAPI } from "../contexts/apiContext";

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

  const { makeQuery } = useAPI();

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

  const { makeQuery } = useAPI();

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

  const { makePaginatedQuery } = useAPI();

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

      setCurrentResponse((prevState) => {
        const uniqItems = inverted
          ? new Set(
              [...response.data, ...prevState.items].map((val) =>
                JSON.stringify(val)
              )
            )
          : new Set(
              [...prevState.items, ...response.data].map((val) =>
                JSON.stringify(val)
              )
            );

        return {
          status: "done",
          items: Array.from(uniqItems).map((val) => JSON.parse(val)),
        };
      });
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
