import React, {useState, useEffect} from 'react'

import { useSearchParams, useNavigate, Navigate } from 'react-router-dom';

import { useAuth } from '../contexts/AuthContext';
import { useNotificationContext } from '../contexts/NotificationContext';
import { InputOtp } from '@heroui/react'
import { authenticateSharing } from '../services/authenticate';
import { SharingGatewayDetails } from '../types';
import config from '../configs/config';
import LoadingPage from './LoadingPage';


const SharingGateway: React.FC = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();

    const {isAuthenticated, setToken} = useAuth();
    const {notifyError} = useNotificationContext();

    const [linkUrl, setLinkUrl] = useState<string | null>(null);
    const [folderId, setFolderId] = useState<string | null>(null);
    const [otpValue, setOtpValue] = useState<string>("");
    const [authLoading, setAuthLoading] = useState<boolean>(false);

    useEffect(() => {
        const linkUrl = searchParams.get("l");

        if (!linkUrl) {
            notifyError("Invalid Link","Missing values on query parameters");
            navigate("/");
        }

        setLinkUrl(linkUrl);
    },[])


    useEffect(() => {
        if (otpValue.length < config.OTP_LENGTH) return;
        if (!linkUrl) return;
        setAuthLoading(true);

        const sharingPayload : SharingGatewayDetails = {
            link_url:   linkUrl,
            otp:       otpValue,
        }
        authenticateSharing(sharingPayload)
            .then((res) => {
                if (!res?.access_token || !res?.folder_id) return;
                setToken(res.access_token);
                setFolderId(res.folder_id)
            })
            .catch(() => {
                notifyError("Authentication Failure", "Encountered issue during authentication process")
            })

    }, [linkUrl, otpValue])

    useEffect(() => {

    }, [isAuthenticated])

    return <>{
        isAuthenticated && folderId ?
            <Navigate to={`/sharing?fid=${folderId}`}/> 
        :   
            authLoading ? 
                <LoadingPage/>
            :
                <div className="flex justify-center items-center w-full h-full">
                    <div className="flex flex-col justify-center items-center gap-2">
                        <span className="font-linik text-wsecondary text-2xl">Folder Password: </span>
                        <InputOtp 
                            length={config.OTP_LENGTH}
                            size='lg'
                            radius="md"
                            variant='faded'
                            classNames={{
                                segment: "bg-wsecondary text-gray-700 text-[14px] border border-gray-300",
                                segmentWrapper: "gap-1 "
                            }}
                            value={otpValue}
                            onValueChange={setOtpValue}
                        />
                    </div>
                </div>
        }</>;

};

export default SharingGateway;