import {
    createContext,
    useState,
    useContext,
    useLayoutEffect,
    ReactNode,
    useEffect,
} from 'react';
import { refresh } from '../services/authenticate.ts';
import api from '../api/api.ts';
import { logout as logoutCallback } from '../services/authenticate.ts';

export interface AuthContextType {
    username: string;
    token: string | null;
    setToken: (token: string | null) => void;
    logout : () => void;
    isAuthenticated: boolean;
    refreshLoading: boolean;
}

interface TokenPayload {
    nickname?: string;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
    const authContext = useContext(AuthContext);

    if (authContext === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }

    return authContext;
};

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [token, setToken] = useState<string | null>(null);
    const [refreshLoading, setRefreshLoading] = useState<boolean>(true);
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);

    useEffect(() => {
        setIsAuthenticated(!!token);
    }, [token])

    const decodeIdToken = (token: string | null): TokenPayload | null => {
        try {
            if (!token) return null;
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(
                atob(base64)
                .split('')
                .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
                .join('')
            );
            return JSON.parse(jsonPayload);
        } catch (error) {
            console.error('Invalid token', error);
            return null;
        }
    };

    const logout = () => {
        logoutCallback();
        setToken(null);
    }

    const tokenPayload = decodeIdToken(token);
    const username = tokenPayload?.nickname || 'User';

    useLayoutEffect(() => {
        refresh()
            .then((res) => {
                setToken(res.access_token);
            })
            .catch(() => setToken(null))
            .finally(() => setRefreshLoading(false));
    }, []);

    useLayoutEffect(() => {
        const authInterceptor = api.interceptors.request.use((config) => {
            config.headers.Authorization =
                !(config as any)._retry && token
                ? `Bearer ${token}`
                : config.headers.Authorization;
            return config;
        });
        return () => {
            api.interceptors.request.eject(authInterceptor);
        };
    }, [token]);

    useLayoutEffect(() => {
        const refreshInterceptor = api.interceptors.response.use(
        (response) => response,
        async (error) => {
            const originalRequest = error.config;
            if (error?.response?.status === 401 && 
                (error?.response?.data?.message === 'Unauthorized' ||
                error?.response?.data?.message === 'The incoming token has expired')
            ) {
                try {
                    const response = await refresh();
                    setToken(response.access_token);
                    originalRequest.headers.Authorization = `Bearer ${response.access_token}`;
                    (originalRequest as any)._retry = true;
                    return api(originalRequest);
                } catch {
                    setToken(null);
                }
            }
            return Promise.reject(error);
        }
        );
        return () => {
            api.interceptors.request.eject(refreshInterceptor);
        };
    }, []);

    return (
        <AuthContext.Provider
            value={{ username, token, setToken, isAuthenticated, logout, refreshLoading }}
        >
            {children}
        </AuthContext.Provider>
    );
};