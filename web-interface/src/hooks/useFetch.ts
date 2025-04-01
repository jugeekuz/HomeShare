import { useState, useEffect } from 'react';
import api from '../api/api';
import { AxiosResponse } from 'axios';

function useFetch<T = any>(url: string): {
    data: T | null;
    isLoading: boolean;
    error: string | null;
} {
    const [data, setData] = useState<T | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async (): Promise<void> => {
        setIsLoading(true);
        setError(null);

        try {
            const response: AxiosResponse<{ body: T }> = await api({
                method: "GET",
                url,
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            setData(response.data.body);
        } catch (err: any) {
            setError(err.response ? err.response.statusText : err.message);
        } finally {
            setIsLoading(false);
        }
        };

        fetchData();
    }, [url]);

    return { data, isLoading, error };
}

export default useFetch;
