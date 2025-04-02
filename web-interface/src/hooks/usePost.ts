import { useState } from 'react';
import api from '../api/api';

interface UsePostReturn<T> {
    postItem:   (item: any) => Promise<T | undefined>;
    loading:    boolean;
    success:    boolean;
    error:      string | null;
    data:       T | null;
}

const usePost = <T>(url: string): UsePostReturn<T> => {
    const [success, setSuccess] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [data, setData] = useState<T | null>(null);

    const postItem = async (item: any): Promise<T | undefined> => {
        try {
            setLoading(true)
            const response = await api({
                method: 'POST',
                url: url,
                headers: {
                'Content-Type': 'application/json',
                },
                data: item,
            });

            setSuccess(true);
            setData(response.data);
            setLoading(false);
            return response.data;
        } catch (e: any) {
            setLoading(false);
            setError(e.response ? e.response.statusText : e.message);
            console.error('Failed to post data:', e);
        } finally {
            setTimeout(() => {
                setSuccess(false);
            }, 200);
        }
    };

    return { postItem, loading, success, error, data };
};

export default usePost;
