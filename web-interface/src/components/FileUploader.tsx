import React, { useState, useEffect } from 'react';
import { Button } from '@heroui/button';
import { Spinner } from '@heroui/spinner';

import { useSearchParams } from 'react-router-dom';
import { useFileUploadContext } from '../contexts/FileUploadContext';

const FileUploader : React.FC = () => {
    const { files, filesReady, filesUploading, uploadFiles } = useFileUploadContext();
    const [searchParams] = useSearchParams();
    const [folderId, setFolderId] = useState<string | undefined>(undefined);
    
    useEffect(() => {
        const folderId = searchParams.get("folder-id")
        if (!folderId) {
            setFolderId(undefined);
            return;
        };
        setFolderId(folderId);
        
    }, [searchParams])

    return (
        <>
        
        <div className="w-full px-2">
            <Button
                color="primary"
                isDisabled={!files || Object.keys(files).length === 0 || !filesReady}
                className="text-md w-full bg-primary-gradient"
                size="md"
                radius="sm"
                onPress={() => uploadFiles(folderId)}
            >
                {((files && !filesReady) || filesUploading) ? <Spinner color="default"/> : "Send Files"}
            </Button>
        </div>
        </>
    )
}

export default FileUploader;