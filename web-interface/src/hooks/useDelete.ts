import { useState, useCallback } from "react";
import api from "../api/api";
const useDelete = (url: string): {
    success: boolean;
    loading: boolean;
    error: string | null;
    refetch: () => Promise<void>;
} => {
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<boolean>(false);

    const deleteItem = useCallback(async (): Promise<void> => {
        setLoading(true);
        setError(null);
        setSuccess(false);
        
        try {
            await api({
                method: "DELETE",
                url,
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            setSuccess(true); 
        } catch (e: any) {
            setError(e.response ? e.response.statusText : e.message);
        } finally {
            setLoading(false);
        }
    }, [url]);

    return { success, loading, error, refetch: deleteItem };
};

export default useDelete