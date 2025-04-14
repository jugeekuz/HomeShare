import axios from 'axios';
import config from '../configs/config';
import { SharingGatewayDetails, AuthSharingResponse, TokenResponse } from '../types';

export const authenticate = async (username: string, password: string) : Promise<TokenResponse> =>{
    const loginUrl = `${config.BASE_URL}/login`
    try {
        const payload = {
            username: username,
            password: password,
        };

        const response = await axios.post(loginUrl, payload, {
            withCredentials: true,
            headers: {
                'Content-Type': 'application/json',
            },
        });
        return response.data;
    } catch (error) {
        throw error;
    }      
};

export const authenticateSharing = async (sharingPayload: SharingGatewayDetails) : Promise<AuthSharingResponse> => {
    const shareUrl = `${config.AUTH_SHARE_URL}`
    try {
        const response = await axios.post(shareUrl, sharingPayload, {
            withCredentials: true,
            headers: {
                'Content-Type': 'application/json',
            },
        });
        return response.data;
    } catch (error) {
        throw error;
    }      
}

export const logout = async () : Promise<TokenResponse> => {
    const logoutUrl = `${config.BASE_URL}/logout`
    try {
        const response = await axios.post(logoutUrl, null, {
            withCredentials: true,
            headers: {
              'Content-Type': 'application/json',
            },
        });
        return response.data;
    } catch (error) {
        throw error;
    }
    
};

export const refresh = async () : Promise<TokenResponse> => {
    const refreshUrl = `${config.BASE_URL}/refresh`;
    try {
        const response = await axios.post(refreshUrl, {}, {
            withCredentials: true,
            headers: {
                'Content-Type': 'application/json',
              }
        })
        return response.data;
    } catch (error) {
        throw error;
    }
}