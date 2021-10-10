import { useEffect, useState } from "react";
import { makeRequest, Request } from "./fetch";
import { connection } from "./ws";

export const useWS = () => {
  const [status, setStatus] = useState<"connecting" | "connected" | "offline">(
    "connecting"
  );

  connection.onopen = () => {
    console.log("Successfully Connected");
    setStatus("connected");
  };

  connection.onclose = (event) => {
    console.log("Socket Closed Connection: ", event);
    setStatus("offline");
  };

  connection.onerror = (error) => {
    console.log("Socket Error: ", error);
    setStatus("offline");
  };

  return {
    status,
  };
};

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
