import { useEffect, useState } from "react";
import { makeRequest, Request } from "./fetch";

export const useRequest = <T,>(
  url: string,
  request?: Request
): { data: T; loading: boolean } => {
  const [data, setData] = useState<any>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      const response = await makeRequest(url, request);
      setData(response.data);
      setLoading(false);
    };

    fetchData();
  }, []);

  return {
    data,
    loading,
  };
};
