import React, {useState, useEffect} from 'react'

import { refresh } from '../services/authenticate';
import { useSearchParams, useNavigate, Navigate } from 'react-router-dom';

import { useAuth } from '../contexts/AuthContext';
import { useNotificationContext } from '../contexts/NotificationContext';
import LoadingPage from './LoadingPage';

const setCookie = (name: string, value: string, minutes: number, domain?: string) => {
    const expires = new Date(Date.now() + minutes * 60 * 1000).toUTCString();
    document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/` + 
        (domain ? `; domain=${domain}` : '');
  };
  

const SharingGateway: React.FC = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();

    const {isAuthenticated, setToken} = useAuth();
    const {notifyError} = useNotificationContext();

    const [folderId, setFolderId] = useState<string | null>(null);
    const [folderName, setFolderName] = useState<string | null>(null);
    const [refreshToken, setRefreshToken] = useState<string | null>(null);
    const [_, setRefreshLoading] = useState<boolean>(true);

    useEffect(() => {
        const folderId = searchParams.get("folder-id");
        const folderName = searchParams.get("folder-name");
        const refreshToken = searchParams.get("refresh");

        if (!folderId || !folderName || !refreshToken) {
            notifyError("Invalid Link","Missing values on query parameters");
            navigate("/");
        }

        setFolderId(folderId);
        setFolderName(folderName);
        setRefreshToken(refreshToken);
    },[])

    useEffect(() => {
        if (!folderId || !folderName || !refreshToken) return;
        if (isAuthenticated) return;

        setCookie("refresh_token", refreshToken, 30)
        setRefreshLoading(true);
        refresh()
            .then((data) => {
                setToken(data.access_token);
                setRefreshLoading(false);
            })
            .catch(() => {
                notifyError("Invalid Credentials", "Encountered an issue while obtaining tokens")
                setToken(null);
                setRefreshLoading(false);
            })
    },[folderId, folderName, refreshToken])

    useEffect(() => {

    }, [isAuthenticated])

    return <>{
        isAuthenticated ?
            <Navigate to={`/sharing?folder-id=${folderId}&folder-name=${folderName}`}/> 
        :
            <LoadingPage/>
        }</>;

};

export default SharingGateway;