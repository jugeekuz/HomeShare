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
    token:              string | null;
    claims:             TokenClaims | null;
    setToken:           (token: string | null) => void;
    logout :            () => void;
    isAuthenticated:    boolean;
    refreshLoading:     boolean;
}

interface TokenClaims {
    user_id:        string;
    folder_id:      string;
    folder_name?:   string;
    access:         "r" | "w" | "rw";
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
    const [claims, setClaims] = useState<TokenClaims | null>(null);
    const [refreshLoading, setRefreshLoading] = useState<boolean>(true);
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);

    useEffect(() => {
        setIsAuthenticated(!!token);
        const tokenClaims = extractClaims(token);
        setClaims(tokenClaims);
        if (tokenClaims === null) {
            return;
        }
    }, [token])

    const decodeIdToken = (token: string | null): TokenClaims | null => {
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

    const extractClaims = (token: string | null): TokenClaims | null => {
        const decodedToken = decodeIdToken(token);
        if (!decodedToken?.user_id || !decodedToken?.folder_id || !decodedToken?.access) return null;
        if (!["r", "w", "rw"].includes(decodedToken.access)) return null;
        const tokenClaims : TokenClaims = decodedToken;
        return tokenClaims;
    }

    const logout = () => {
        logoutCallback();
        setToken(null);
    }

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
                if (
                    error?.response?.status === 401 &&
                    !originalRequest._retry &&
                    (error?.response?.data?.trim() === 'Unauthorized' ||
                        error?.response?.data?.trim() === 'The incoming token has expired')
                ) {
                    originalRequest._retry = true;
                    try {
                        const response = await refresh();
                        setToken(response.access_token);
                        originalRequest.headers = {
                            ...originalRequest.headers,
                            Authorization: `Bearer ${response.access_token}`,
                        };
                        return api(originalRequest);
                    } catch (refreshError) {
                        setToken(null);
                        return Promise.reject(refreshError);
                    }
                }
    
                return Promise.reject(error);
            }
        );
    
        return () => {
            api.interceptors.response.eject(refreshInterceptor);
        };
    }, []);

    return (
        <AuthContext.Provider
            value={{ token, setToken, isAuthenticated, claims, logout, refreshLoading }}
        >
            {children}
        </AuthContext.Provider>
    );
};